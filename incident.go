package main

import (
	"DocPlanner/pingdom-statuspage-integration/statuspage"
	"fmt"
	"time"
)

func AsyncIncidentCheck(ticker *time.Ticker, statuspageClient *statuspage.Client, is *incidentStore) {

	for {
		select {
		case <-ticker.C:
			components := CheckEvaluation(is)
			if len(components) > 0 {
				createIncidents(components, statuspageClient, is)
			}
		}
	}

}

func createIncidents(components []*component, statuspageClient *statuspage.Client, is *incidentStore) {
	componentsGroup := make(map[string][]string)

	for _, component := range components {
		componentsGroup[component.pageId] = append(componentsGroup[component.pageId], component.id)
	}

	// create incidents
	for pageID, componentsList := range componentsGroup {
		incident, err := statuspageClient.CreateIncident(componentsList, pageID)
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
