package main

import (
	"DocPlanner/pingdom-statuspage-integration/mock"
	"DocPlanner/pingdom-statuspage-integration/statuspage"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func performRequest(secret string, r http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))

	query := req.URL.Query()
	query.Set("secret", secret)
	req.URL.RawQuery = query.Encode()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestBananaAuth(t *testing.T) {
	http.DefaultTransport = buildTransport(nil)

	router := SetupRouter(statuspage.NewClient("FAKE_TAXI_I_MEAN_TOKEN"), "SUPER_SECRET", nil, nil)

	rsp := performRequest("INCORRECT_SECRET", router, http.MethodPost, "/", PingdomPayloadHTTPUpToDown)

	assert.Emptyf(t, rsp.Body.String(), "Response not empty!")
	assert.Equal(t, http.StatusUnauthorized, rsp.Code)
}

func TestIntegrationHappyPath(t *testing.T) {
	transport := buildTransport(&map[mock.Request]mock.Response{
		{
			Method:  "PATCH",
			Host:    "api.statuspage.io",
			URLPath: "/v1/pages/2y38pys158vc/components/123123",
		}: {
			Body: `{}`, // there is a body but if response code is OK then we don't care
		},
	})
	http.DefaultTransport = transport

	router := SetupRouter(statuspage.NewClient("FAKE_TAXI_I_MEAN_TOKEN"), "SUPER_SECRET", nil, nil)

	var response Response

	rsp := performRequest("SUPER_SECRET", router, http.MethodPost, "/", PingdomPayloadHTTPUpToDown)
	unmarshallErr := json.Unmarshal(rsp.Body.Bytes(), &response)

	assert.Nil(t, unmarshallErr)

	assert.Equal(t, http.StatusOK, rsp.Code)
	assert.Equal(t, "OK", response.Status)

	assert.Emptyf(t, transport.Responses, "Not all responses used")
}

func buildTransport(additionalResponses *map[mock.Request]mock.Response) *mock.Transport {
	responseMap := map[mock.Request]mock.Response{
		{
			Method:  "GET",
			Host:    "api.statuspage.io",
			URLPath: "/v1/pages",
		}: {
			Body: `[
					  {
						"id": "2y38pys158vc",
						"created_at": "2020-01-31T21:23:23Z",
						"updated_at": "2020-01-31T21:23:23Z",
						"name": "My Company Status",
						"page_description": "string",
						"headline": "string",
						"branding": "string",
						"subdomain": "your-subdomain.statuspage.io",
						"domain": "status.mycompany.com",
						"url": "https://www.mycompany.com",
						"support_url": "string",
						"hidden_from_search": true,
						"allow_page_subscribers": true,
						"allow_incident_subscribers": true,
						"allow_email_subscribers": true,
						"allow_sms_subscribers": true,
						"allow_rss_atom_feeds": true,
						"allow_webhook_subscribers": true,
						"notifications_from_email": "no-reply@status.mycompany.com",
						"notifications_email_footer": "string",
						"activity_score": 0,
						"twitter_username": "string",
						"viewers_must_be_team_members": true,
						"ip_restrictions": "string",
						"city": "string",
						"state": "string",
						"country": "string",
						"time_zone": "UTC",
						"css_body_background_color": "string",
						"css_font_color": "string",
						"css_light_font_color": "string",
						"css_greens": "string",
						"css_yellows": "string",
						"css_oranges": "string",
						"css_blues": "string",
						"css_reds": "string",
						"css_border_color": "string",
						"css_graph_color": "string",
						"css_link_color": "string",
						"favicon_logo": "string",
						"transactional_logo": "string",
						"hero_cover": "string",
						"email_logo": "string",
						"twitter_logo": "string"
					  }
					]`,
		},
		{
			Method:  "GET",
			Host:    "api.statuspage.io",
			URLPath: "/v1/pages/2y38pys158vc/components",
		}: {
			Body: `[
					  {
						"id": "123123",
						"page_id": "2y38pys158vc",
						"group_id": "string",
						"created_at": "2020-01-31T21:23:24Z",
						"updated_at": "2020-01-31T21:23:24Z",
						"group": true,
						"name": "Name of HTTP check",
						"description": "string",
						"position": 0,
						"status": "operational",
						"showcase": true,
						"only_show_if_degraded": true,
						"automation_email": "string"
					  }
					]`,
		},
	}

	if additionalResponses != nil {
		for req, rsp := range *additionalResponses {
			responseMap[req] = rsp
		}
	}

	return mock.NewTransport(&responseMap)
}

const PingdomPayloadHTTPUpToDown = `{
      "check_id": 12345,
      "check_name": "Name of HTTP check",
      "check_type": "HTTP",
      "check_params": {
        "basic_auth": false,
        "encryption": true,
        "full_url": "https://www.example.com/path",
        "header": "User-Agent:Pingdom.com_bot",
        "hostname": "www.example.com",
        "ipv6": false,
        "port": 443,
        "url": "/path"
      },
      "tags": [
        "example_tag"
      ],
      "previous_state": "UP",
      "current_state": "DOWN",
      "importance_level": "HIGH",
      "state_changed_timestamp": 1451610061,
      "state_changed_utc_time": "2016-01-01T01:01:01",
      "long_description": "Long error message",
      "description": "Short error message",
      "first_probe": {
        "ip": "123.4.5.6",
        "ipv6": "2001:4800:1020:209::5",
        "location": "Stockholm, Sweden"
      },
      "second_probe": {
        "ip": "123.4.5.6",
        "ipv6": "2001:4800:1020:209::5",
        "location": "Austin, US",
        "version": 1
      }
    }`
