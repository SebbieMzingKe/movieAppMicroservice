package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"movieapp.com/pkg/discovery"
)

type serviceName string
type instanceId string

type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceId]*serviceInstance
}

type serviceInstance struct {
	hostPort  string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[serviceName]map[instanceId]*serviceInstance{}}
}

func (r *Registry) Register(ctx context.Context, instId string, sName string, hostPort string) error {
	r.Lock()

	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(sName)]; !ok {
		r.serviceAddrs[serviceName(sName)] = map[instanceId]*serviceInstance{}
	}
	r.serviceAddrs[serviceName(sName)][instanceId(instId)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (r *Registry) Deregister(ctx context.Context, instId string, sName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(sName)]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceName(sName)], instanceId(instId))
	return nil
}


func (r *Registry) ReportHealthyStatus(instID string, sName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(sName)]; !ok {
		return errors.New("service is not registered yet")
	}

	if _, ok := r.serviceAddrs[serviceName(sName)][instanceId(instID)]; !ok {
		return errors.New("service instance not registered yet")
	}

	r.serviceAddrs[serviceName(sName)][instanceId(instID)].lastActive = time.Now()

	return nil
}

func (r *Registry) ServiceAddresses(ctx context.Context, sName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.serviceAddrs[serviceName(sName)]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string

	for _, i := range r.serviceAddrs[serviceName(sName)] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}

		res =  append(res, i.hostPort)
	}

	return res, nil
}