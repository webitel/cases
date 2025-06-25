package consul

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	conf "github.com/webitel/cases/config"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/registry"
)

type ConsulRegistry struct {
	registrationConfig *consulapi.AgentServiceRegistration
	client             *consulapi.Client
	stop               chan any
	checkId            string
}

// NewConsulRegistry creates a new Consul registry instance.
func NewConsulRegistry(config *conf.ConsulConfig) (*ConsulRegistry, error) {
	var err error
	entity := ConsulRegistry{}
	if config.Id == "" {
		return nil, errors.Internal(
			"service id is empty! (set it by '-id' flag)",
			errors.WithID("consul.registry.new_consul.check_args.service_id"),
		)
	}
	ip, port, err := net.SplitHostPort(config.PublicAddress)
	if err != nil {
		return nil, errors.Internal(
			"unable to parse address",
			errors.WithID("consul.registry.new_consul.parse_address.error"),
		)
	}
	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.Internal(
			"unable to parse ip",
			errors.WithID("consul.registry.new_consul.parse_ip.error"),
		)
	}

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = config.Address
	entity.client, err = consulapi.NewClient(consulConfig)
	if err != nil {
		return nil, errors.Internal(
			err.Error(),
			errors.WithID("consul.registry.new_consul_registry.consulapi_creation.error"),
		)
	}

	entity.registrationConfig = &consulapi.AgentServiceRegistration{
		ID:      config.Id,
		Name:    registry.ServiceName,
		Port:    parsedPort,
		Address: ip,
		Check: &consulapi.AgentServiceCheck{
			DeregisterCriticalServiceAfter: registry.DeregisterCriticalServiceAfter.String(),
			TTL:                            registry.CheckInterval.String(),
		},
	}
	entity.stop = make(chan any)

	return &entity, nil
}

// Register registers the service with Consul.
func (c *ConsulRegistry) Register() error {
	err := c.client.Agent().ServiceRegister(c.registrationConfig)
	if err != nil {
		return errors.Internal(
			err.Error(),
			errors.WithID("consul.registry.consul.register.error"),
		)
	}
	var checks map[string]*consulapi.AgentCheck
	if checks, err = c.client.Agent().Checks(); err != nil {
		return errors.Internal(
			err.Error(),
			errors.WithID("consul.registry.consul.register.get_checks.error"),
		)
	}

	var serviceCheck *consulapi.AgentCheck
	for _, check := range checks {
		if check.ServiceID == c.registrationConfig.ID {
			serviceCheck = check
		}
	}

	if serviceCheck == nil {
		return errors.Internal(
			"service check not found",
			errors.WithID("consul.registry.consul.register.error"),
		)
	}
	c.checkId = serviceCheck.CheckID
	go c.RunServiceCheck()
	return nil
}

func (c *ConsulRegistry) Deregister() error {
	err := c.client.Agent().ServiceDeregister(c.registrationConfig.ID)
	if err != nil {
		return errors.Internal(
			err.Error(),
			errors.WithID("consul.registry.consul.deregister.error"),
		)
	}
	c.stop <- true
	slog.Info(fmtConsulLog("service was deregistered"))
	return nil
}

func (c *ConsulRegistry) doUpdateTTL() error {
	err := c.client.Agent().UpdateTTL(c.checkId, "success", "pass")
	if err != nil {
		slog.Error("consul: failed to complete regular check-in", "error", fmtConsulLog(err.Error()))
		return err
	}
	return nil // [OK]
}

func (c *ConsulRegistry) RunServiceCheck() error {
	// register: now !
	if err := c.doUpdateTTL(); err == nil {
		slog.Info(fmtConsulLog("service was registered"))
	}
	defer slog.Info(fmtConsulLog("stopped service checker"))
	slog.Info(fmtConsulLog("started service checker"))
	ticker := time.NewTicker(registry.CheckInterval / 2)
	defer ticker.Stop()
	for {
		select {
		case <-c.stop:
			// gracefull stop
			return nil
		case <-ticker.C:
			_ = c.doUpdateTTL() // regular: check-in
			// TODO: seems that connection is lost, reconnect?
		}
	}
}

func fmtConsulLog(s string) string {
	return fmt.Sprintf("consul: %s", s)
}
