package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"movieexample.com/pkg/discovery"
)

type Registry struct {
	client *consul.Client
}

func NewRegistry(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {

		return nil, err

	}

	return &Registry{client: client}, nil
}

func (r *Registry) Register(ctx context.Context, instanceId string, svcName string, hostPort string) error {

	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in a form of <host>:<port>, example: localhost:8081")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      instanceId,
		Name:    svcName,
		Address: parts[0],
		Port:    port,
		Check:   &consul.AgentServiceCheck{CheckID: instanceId, TTL: "5s"},
	})
}

func (r *Registry) Deregister(ctx context.Context, instanceId string, _ string) error {
	return r.client.Agent().ServiceDeregister(instanceId)
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)

	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string

	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}

	return res, nil
}

func (r *Registry) ReportHealthyState(instanceID string, _ string) error {
	return r.client.Agent().PassTTL(instanceID, "")
}
