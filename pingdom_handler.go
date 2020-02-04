package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type PingdomPayload struct {
	CheckId             int      `json:"check_id"`
	CheckName           string   `json:"check_name"`
	CheckType           string   `json:"check_type"`
	Tags                []string `json:"tags"`
	PreviousState       string   `json:"previous_state"`
	CurrentState        string   `json:"current_state"`
	StateChangedUtcTime string   `json:"state_changed_utc_time"`
}

func pingdomHandler(c *gin.Context) {
	cs := c.MustGet("component_store").(*componentsStore)

	var pingdomPayload PingdomPayload
	if err := c.ShouldBindJSON(&pingdomPayload); err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//operational - UP
	//major_outage, partial_outage, degraded_performance - DOWN
	status := "operational"
	if pingdomPayload.CurrentState != "UP" {
		status = "major_outage"
	}

	components := cs.FindComponentsByName(pingdomPayload.CheckName)
	if len(components) == 0 {
		// TODO: implement Slack notification
		err := fmt.Errorf("Components matching name \"%s\" not found on StatusPage. Pingdom check status changed to: %s", pingdomPayload.CheckName, pingdomPayload.CurrentState)
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var errs []string
	for _, cmp := range components {
		err := cs.UpdateComponentStatus(cmp, status)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		_ = c.Error(fmt.Errorf(strings.Join(errs, ",")))

		c.JSON(http.StatusInternalServerError, gin.H{"errors": errs})
		return
	}

	c.JSON(http.StatusOK, Response{
		Status: "OK",
	})
}
