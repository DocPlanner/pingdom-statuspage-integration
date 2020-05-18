package main

import (
	"DocPlanner/pingdom-statuspage-integration/statuspage"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func InitializeComponentsStore(statusPageClient *statuspage.Client, componentsStoreChan chan *componentsStore, is *incidentStore) gin.HandlerFunc {
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
		context.Set("incident_store", is)
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

func SetupRouter(statusPageClient *statuspage.Client, secret string, componentsStoreChan chan *componentsStore, is *incidentStore) *gin.Engine {
	router := gin.Default()

	router.Use(BananaAuthMiddleware(secret))
	router.Use(InitializeComponentsStore(statusPageClient, componentsStoreChan, is))

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
		return
	})

	router.POST("/", pingdomHandler)

	return router
}
