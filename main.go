package main

import (
	"DocPlanner/pingdom-statuspage-integration/statuspage"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Response struct {
	Status string `json:"status"`
}

func main() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	componentsStoreChan := make(chan *componentsStore)
	go AsyncRefresh(ticker, componentsStoreChan)

	secret := getSecret()
	statusPageClient := setupStatusPageClient()

	router := SetupRouter(statusPageClient, secret, componentsStoreChan)

	_ = router.Run(":80")
}

func getSecret() string {
	secret := os.Getenv("SECRET")
	if len(secret) == 0 {
		fmt.Println("Environment variable \"SECRET\" not set!")
		os.Exit(1)
	}

	return secret
}

func setupStatusPageClient() *statuspage.Client {
	statuspageToken := os.Getenv("STATUSPAGE_TOKEN")
	if len(statuspageToken) == 0 {
		fmt.Println("Environment variable \"STATUSPAGE_TOKEN\" not set!")
		os.Exit(1)
	}

	statusPageClient := statuspage.NewClient(statuspageToken)

	envMaxRetries := os.Getenv("MAX_RETRIES")
	if len(envMaxRetries) > 0 {
		retries, _ := strconv.Atoi(envMaxRetries)
		statusPageClient.MaxRetries = retries
	}

	envRetryInterval := os.Getenv("RETRY_INTERVAL")
	if len(envRetryInterval) > 0 {
		interval, _ := strconv.Atoi(envRetryInterval)
		statusPageClient.RetryInterval = time.Duration(interval) * time.Second
	}

	return statusPageClient
}
