package main

import (
	"mqtt-pubsub/handlers"
	mqttpubsub "mqtt-pubsub/modules/mqtt-pubsub"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	debug.SetGCPercent(10)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./webapp/dist/index.html")
		} else {
			c.File("./webapp/dist/" + path.Join(dir, file))
		}
	})

	r.GET("/mqtt-pubsub/config", handlers.GetConfigHandler)
	r.POST("/mqtt-pubsub/config", handlers.SetConfigHandler)

	go mqttpubsub.Run()

	go func() {
		err := r.Run(":9091")
		if err != nil {
			panic(err)
		}
	}()

	<-sig
	mqttpubsub.CloseConnections()
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
