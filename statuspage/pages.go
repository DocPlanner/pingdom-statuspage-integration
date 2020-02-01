package statuspage

type Page struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (client *Client) ListPages() (pages []Page, err error) {
	rsp, err := client.doGET("/pages", nil)
	if err != nil {
		return nil, err
	}

	err = client.unmarshal(rsp, &pages)
	if err != nil {
		return nil, err
	}

	return pages, nil
}
