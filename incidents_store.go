package main

import (
	"DocPlanner/pingdom-statuspage-integration/statuspage"
	"fmt"
	"sync"
	"time"
)

const DEFAULT_INCIDENT_LAG_SECONDS = 60

type incidentStore struct {
	components []*component
	mux        sync.Mutex
}

type component struct {
	timestamp int64
	id        string
	pageId    string
}

func (is *incidentStore) Add(componentID string, pageID string, timestamp int64) {
	is.mux.Lock()
	candidate := &component{
		id:        componentID,
		pageId:    pageID,
		timestamp: timestamp,
	}
	is.components = append(is.components, candidate)
	is.mux.Unlock()
}

func (is *incidentStore) GetAll() []*component {
	is.mux.Lock()
	defer is.mux.Unlock()
	return is.components
}

func (is *incidentStore) Remove(componentID string) {
	is.mux.Lock()
	componentIndex := 0
	for index, component := range is.components {
		if component.id == componentID {
			componentIndex = index
			break
		}
	}
	if len(is.components) > 0 {
		is.components = append(is.components[:componentIndex], is.components[componentIndex+1:]...)
	}
	is.mux.Unlock()
}

func (is *incidentStore) CheckEvaluation() []*component {

	var components []*component
	incidentStore := is.GetAll()
	now := time.Now()

	for _, candidate := range incidentStore {
		diff := now.Sub(time.Unix(candidate.timestamp, 0)).Seconds()
		if diff > DEFAULT_INCIDENT_LAG_SECONDS {
			components = append(components, candidate)
		}
	}

	return components
}

func shouldAddAsCandidate(status string) bool {
	return status == "major_outage"
}

func (is *incidentStore) updateIncidentStore(cmp statuspage.Component, status string) {

	timestamp := time.Now().Unix()
	if shouldAddAsCandidate(status) {
		is.Add(cmp.ID, cmp.PageID, timestamp)
		return
	}

	is.Remove(cmp.ID)

}
