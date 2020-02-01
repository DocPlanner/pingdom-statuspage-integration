package main

import (
	"DocPlanner/pingdom-statuspage-integration/statuspage"
	"errors"
)

type componentsStore struct {
	StatusPageClient *statuspage.Client
	Components       []statuspage.Component
}

func NewComponentsStore(statusPageClient *statuspage.Client) *componentsStore {
	return &componentsStore{
		StatusPageClient: statusPageClient,
	}
}

func (cs *componentsStore) Refresh() error {
	pages, err := cs.StatusPageClient.ListPages()
	if err != nil {
		return errors.New("StatusPage: " + err.Error())
	}

	var components []statuspage.Component
	for _, page := range pages {
		cmp, err := cs.StatusPageClient.ListComponents(page)
		if err != nil {
			return errors.New("StatusPage: " + err.Error())
		}

		components = append(components, cmp...)
	}

	cs.Components = components

	return nil
}

func (cs *componentsStore) FindComponentsByName(name string) []statuspage.Component {
	var components []statuspage.Component
	for _, cmp := range cs.Components {
		if cmp.Name == name {
			components = append(components, cmp)
		}
	}

	return components
}

func (cs *componentsStore) UpdateComponent(component statuspage.Component) error {
	return cs.StatusPageClient.UpdateComponent(component)
}
