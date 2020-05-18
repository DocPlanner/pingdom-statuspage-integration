package main

import (
	"DocPlanner/pingdom-statuspage-integration/statuspage"
	"fmt"
	"sync"
	"time"
)

const DEFAULT_INCIDENT_LAG_SECONDS = 60

type incidentClient struct {
	StatusPageClient *statuspage.Client
}

type incidentStore struct {
	components []*component
	mux        sync.Mutex
}

type component struct {
	timestamp int64
	id        string
	pageId    string
}

func NewIncidentClient(statusPageClient *statuspage.Client) *incidentClient {
	return &incidentClient{
		StatusPageClient: statusPageClient,
	}
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
	now := time.Now().Unix()

	for _, candidate := range incidentStore {
		relative := now - candidate.timestamp
		if relative > DEFAULT_INCIDENT_LAG_SECONDS {
			components = append(components, candidate)
		}
	}

	return components
}

func AsyncIncidentCheck(ticker *time.Ticker, is *incidentStore, ic *incidentClient) {

	for {
		select {
		case <-ticker.C:
			components := is.CheckEvaluation()
			if len(components) > 0 {

				// currently support single page
				var componentsList []string
				var pageID string
				for _, component := range components {
					componentsList = append(componentsList, component.id)
					pageID = component.pageId
				}

				// create incidents
				incident, err := ic.StatusPageClient.CreateIncident(componentsList, pageID)
				if err != nil {
					fmt.Errorf("StatusPage: " + err.Error())
					return
				}
				fmt.Printf("Incident %s created for components: %s", incident.ID, componentsList)

				// TODO: add here slack notifications

				// if incident created, remove component from store
				for _, component := range components {
					is.Remove(component.id)
				}
			}
		}
	}

}

func shouldAddAsCandidate(status string) bool {
	return status == "major_outage"
}

func (is *incidentStore) updateIncidentStore(cmp statuspage.Component, status string) {

	timestamp := time.Now().Unix()
	switch shouldAddAsCandidate(status) {
	case true:
		is.Add(cmp.ID, cmp.PageID, timestamp)
	case false:
		is.Remove(cmp.ID)
	}

}
