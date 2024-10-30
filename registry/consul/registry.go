package consul

import (
	"net"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	conf "github.com/webitel/cases/config"
	rerr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/registry"
)

type ConsulRegistry struct {
	registrationConfig *consulapi.AgentServiceRegistration
	client             *consulapi.Client
}

// NewConsulRegistry creates a new Consul registry instance.
func NewConsulRegistry(config *conf.ConsulConfig) (*ConsulRegistry, error) {
	var err error
	entity := ConsulRegistry{}
	if config.Id == "" {
		return nil, rerr.NewRegistryError("consul.registry.new_consul.check_args.service_id", "service id is empty! (set it by '-id' flag)")
	}
	ip, port, err := net.SplitHostPort(config.PublicAddress)
	if err != nil {
		return nil, rerr.NewRegistryError("consul.registry.new_consul.parse_address.error", "unable to parse address")
	}
	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, rerr.NewRegistryError("consul.registry.new_consul.parse_ip.error", "unable to parse ip")
	}

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = config.Address
	entity.client, err = consulapi.NewClient(consulConfig)
	if err != nil {
		return nil, rerr.NewRegistryError("consul.registry.new_consul_registry.consulapi_creation.error", err.Error())
	}

	entity.registrationConfig = &consulapi.AgentServiceRegistration{
		ID:      config.Id,
		Name:    registry.ServiceName,
		Port:    parsedPort,
		Address: ip,
		Check: &consulapi.AgentServiceCheck{
			DeregisterCriticalServiceAfter: registry.DeregisterCriticalServiceAfter.String(),
			CheckID:                        config.Id,
			TCP:                            config.PublicAddress,
			Interval:                       registry.CheckInterval.String(),
		},
	}

	return &entity, nil
}

// Register registers the service with Consul.
func (c *ConsulRegistry) Register() error {
	err := c.client.Agent().ServiceRegister(c.registrationConfig)
	if err != nil {
		return rerr.NewRegistryError("consul.registry.consul.register.error", err.Error())
	}

	return nil
}

// Deregister deregisters the service from Consul.
func (c *ConsulRegistry) Deregister() error {
	err := c.client.Agent().ServiceDeregister(c.registrationConfig.ID)
	if err != nil {
		return rerr.NewRegistryError("consul.registry.consul.deregister.error", err.Error())
	}

	return nil
}
