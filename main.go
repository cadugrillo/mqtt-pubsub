package main

import (
	"fmt"
	"mqtt-pubsub/handlers"
	mqttpubsub "mqtt-pubsub/modules/mqtt-pubsub"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func main() {
	debug.SetGCPercent(10)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/mqtt-pubsub/config", handlers.GetConfigHandler)
	r.POST("/mqtt-pubsub/config", handlers.SetConfigHandler)

	go mqttpubsub.Run()
	defer CloseConnections()

	err := r.Run(":9091")
	if err != nil {
		panic(err)
	}

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, accept")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func CloseConnections() {
	fmt.Println("signal caught - exiting")
	mqttpubsub.ClientSub.Disconnect(1000)
	mqttpubsub.ClientPub.Disconnect(1000)
	fmt.Println("shutdown complete")
}
