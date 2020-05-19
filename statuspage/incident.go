package statuspage

type Incident struct {
	Name         string   `json:"name,omitempty"`
	ID           string   `json:"id,omitempty"`
	Status       string   `json:"status,omitempty"`
	ComponentIDs []string `json:"component_ids,omitempty"`
	ShortLink    string   `json:"shortlink,omitempty"`
}

type IncidentPostPayload struct {
	Incident *Incident `json:"incident"`
}

const incidentName = "Performance degraded"
const incidentStatus = "investigating"

func (client *Client) CreateIncident(components []string, pageID string) (incident *Incident, err error) {
	rsp, err := client.doPOST("/pages/"+pageID+"/incidents", &IncidentPostPayload{
		Incident: &Incident{
			Name:         incidentName,
			Status:       incidentStatus,
			ComponentIDs: components,
		}})
	if err != nil {
		return nil, err
	}

	err = client.unmarshal(rsp, &incident)
	if err != nil {
		return nil, err
	}

	return incident, nil
}
