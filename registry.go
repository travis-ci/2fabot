package main

import (
	"sync"

	"github.com/pkg/errors"
)

var (
	serviceRegistry      = map[string]*serviceRegistryEntry{}
	serviceRegistryMutex sync.Mutex
)

type serviceRegistryEntry struct {
	Name        string
	DocURL      string
	ServiceFunc func(config string) (Service, error)
}

func NewService(name, config string) (Service, error) {
	serviceRegistryMutex.Lock()
	defer serviceRegistryMutex.Unlock()

	registryEntry, ok := serviceRegistry[name]
	if !ok {
		return nil, errors.Errorf("unknown service: %s", name)
	}

	return registryEntry.ServiceFunc(config)
}

func docURLForService(name string) string {
	serviceRegistryMutex.Lock()
	defer serviceRegistryMutex.Unlock()

	return serviceRegistry[name].DocURL
}

func registerService(name, docURL string, serviceFunc func(config string) (Service, error)) {
	serviceRegistryMutex.Lock()
	defer serviceRegistryMutex.Unlock()

	serviceRegistry[name] = &serviceRegistryEntry{
		Name:        name,
		DocURL:      docURL,
		ServiceFunc: serviceFunc,
	}
}
