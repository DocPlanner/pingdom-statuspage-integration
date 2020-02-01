package statuspage

type Component struct {
	ID          string `json:"id,omitempty"`
	PageID      string `json:"page_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func (client *Client) ListComponents(page Page) (components []Component, err error) {
	rsp, err := client.doGET("/pages/"+page.ID+"/components", nil)
	if err != nil {
		return nil, err
	}

	err = client.unmarshal(rsp, &components)
	if err != nil {
		return nil, err
	}

	return components, nil
}

func (client *Client) UpdateComponent(component Component) error {
	return client.doPATCH("/pages/"+component.PageID+"/components/"+component.ID, component)
}
