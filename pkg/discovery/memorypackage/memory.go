package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"movieexample.com/pkg/discovery"
)

type serviceName string
type instanceID string

type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func New() *Registry {
	return &Registry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instanceId string, svcName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(svcName)]; !ok {
		r.serviceAddrs[serviceName(svcName)] = map[instanceID]*serviceInstance{}
	}

	r.serviceAddrs[serviceName(svcName)][instanceID(instanceId)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}

	return nil
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(ctx context.Context, instanceId string , svcName string) error {
	r.Lock()
	defer r.Unlock()

	if _,ok := r.serviceAddrs[serviceName(svcName)]; !ok {
		return nil
	}

	delete(r.serviceAddrs[serviceName(svcName)],instanceID(instanceId))
	return nil
}

// ReportHealthyState is a push mechanism for
// reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceId string, srvName string) error {
	if _,ok := r.serviceAddrs[serviceName(srvName)];!ok {
		return errors.New("service is not registered yet")
	}
	if _,ok := r.serviceAddrs[serviceName(srvName)][instanceID(instanceId)]; !ok {
		return errors.New("service instance is not registered yet")
	}

	r.serviceAddrs[serviceName(srvName)][instanceID(instanceId)].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of
// active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, srvName string) ([]string, error) {
	r.Lock()
	defer r.Unlock()
	if len(r.serviceAddrs[serviceName(srvName)]) == 0 {
        return nil, discovery.ErrNotFound
    }

    var res []string
	for _,i := range r.serviceAddrs[serviceName(srvName)] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}

		res = append(res, i.hostPort)
	}

	return res,nil
}