package consul

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/registry"

	"net"
	"strconv"
)

type ConsulRegistry struct {
	registrationConfig *consulapi.AgentServiceRegistration
	client             *consulapi.Client
}

func NewConsulRegistry(config *model.ConsulConfig) (*ConsulRegistry, model.AppError) {
	var (
		err error
	)
	entity := ConsulRegistry{}
	if config.Id == "" {
		return nil, model.NewBadRequestError("consul.registry.new_consul.check_args.service_id", "service id is empty! (set it by '-id' flag)")
	}
	ip, port, err := net.SplitHostPort(config.PublicAddress)
	if err != nil {
		return nil, model.NewBadRequestError("consul.registry.new_consul.parse_address.error", "unable to parse address")
	}
	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, model.NewBadRequestError("consul.registry.new_consul.parse_ip.error", "unable to parse ip")
	}

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = config.Address
	entity.client, err = consulapi.NewClient(consulConfig)
	if err != nil {
		return nil, model.NewBadRequestError("consul.registry.new_consul_registry.consulapi_creation.error", err.Error())
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

func (c *ConsulRegistry) Register() model.AppError {
	err := c.client.Agent().ServiceRegister(c.registrationConfig)
	if err != nil {
		return model.NewInternalError("consul.registry.consul.register.error", err.Error())
	}

	return nil
}

func (c *ConsulRegistry) Deregister() model.AppError {
	err := c.client.Agent().ServiceDeregister(c.registrationConfig.ID)
	if err != nil {
		return model.NewInternalError("consul.registry.consul.register.error", err.Error())
	}

	return nil
}
