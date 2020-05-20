package main

import "time"

const DEFAULT_INCIDENT_LAG_SECONDS = 150

func CheckEvaluation(is *incidentStore) []*component {

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
