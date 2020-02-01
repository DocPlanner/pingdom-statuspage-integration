package main

import (
	"DocPlanner/pingdom-statuspage-integration/statuspage"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)
import "github.com/gin-gonic/gin"

type Response struct {
	Status string `json:"status"`
}

func main() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	componentsStoreChan := make(chan *componentsStore)
	go func() {
		var cs *componentsStore
		for {
			select {
			case a := <-componentsStoreChan:
				cs = a
			case <-ticker.C:
				err := cs.Refresh()
				fmt.Println(fmt.Sprintf("[%s] Refreshing StatusPage components state! %s", time.Now().Format(time.RFC1123Z), err.Error()))
			}
		}
	}()

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

func InitializeComponentsStore(statusPageClient *statuspage.Client, componentsStoreChan chan *componentsStore) gin.HandlerFunc {
	componentStore := NewComponentsStore(statusPageClient)
	err := componentStore.Refresh()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if componentsStoreChan != nil {
		componentsStoreChan <- componentStore
	}

	return func(context *gin.Context) {
		context.Set("component_store", componentStore)
		context.Next()
	}
}

func BananaAuthMiddleware(secret string) gin.HandlerFunc {
	return func(context *gin.Context) {
		isHealthCheck := context.Request.URL.Path == "/healthcheck"
		secretIsCorrect := context.Query("secret") == secret

		if !isHealthCheck && !secretIsCorrect {
			_ = context.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Authorization Required"))
		}

		context.Next()
	}
}

func SetupRouter(statusPageClient *statuspage.Client, secret string, componentsStoreChan chan *componentsStore) *gin.Engine {
	router := gin.Default()

	router.Use(BananaAuthMiddleware(secret))
	router.Use(InitializeComponentsStore(statusPageClient, componentsStoreChan))

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
		return
	})

	router.POST("/", pingdomHandler)

	return router
}
