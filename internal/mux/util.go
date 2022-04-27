package mux

import (
	"strings"
)

func (mux *Mux) GetServiceByName(name string) Service {
	for _, service := range mux.services {
		if service.Name() == name {
			return service
		}
	}

	return nil
}

func (mux *Mux) GetInteractorByName(name string) Interactor {
	for _, interactor := range mux.interactors {
		if interactor.Name() == name {
			return interactor
		}
	}

	return nil
}

func (mux *Mux) GetServicesNames() []string {
	services := make([]string, len(mux.services))
	for _, s := range mux.services {
		services = append(services, strings.Title(s.Name()))
	}
	return services
}
