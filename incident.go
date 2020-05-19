package main

import (
	"DocPlanner/pingdom-statuspage-integration/statuspage"
	"fmt"
	"time"
)

type incidentClient struct {
	StatusPageClient *statuspage.Client
}

func NewIncidentClient(statusPageClient *statuspage.Client) *incidentClient {
	return &incidentClient{
		StatusPageClient: statusPageClient,
	}
}

func AsyncIncidentCheck(ticker *time.Ticker, is *incidentStore, ic *incidentClient) {

	for {
		select {
		case <-ticker.C:
			components := is.CheckEvaluation()
			if len(components) > 0 {
				createIncidents(components, ic, is)
			}
		}
	}

}

func createIncidents(components []*component, ic *incidentClient, is *incidentStore) {
	componentsGroup := make(map[string][]string)

	for _, component := range components {
		componentsGroup[component.pageId] = append(componentsGroup[component.pageId], component.id)
	}

	// create incidents
	for pageID, componentsList := range componentsGroup {
		incident, err := ic.StatusPageClient.CreateIncident(componentsList, pageID)
		if err != nil {
			fmt.Errorf("StatusPage: " + err.Error())
			return
		}
		fmt.Printf("Incident %s (pageID %s) created for components: %s", incident.ID, pageID, componentsList)
	}

	// TODO: add here slack notifications

	// if incident created, remove component from store
	for _, component := range components {
		is.Remove(component.id)
	}
}
