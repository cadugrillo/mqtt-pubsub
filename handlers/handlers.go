package handlers

import (
	"encoding/json"
	"io"
	config_parser "mqtt-pubsub/modules/config-parser"
	"mqtt-pubsub/modules/configurator"
	mqttpubsub "mqtt-pubsub/modules/mqtt-pubsub"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, configurator.GetConfig())
	return
}

func SetConfigHandler(c *gin.Context) {

	configFile, statusCode, err := convertHTTPBodyConfig(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, configurator.SetConfig(configFile))
	return
}

func StartServiceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, mqttpubsub.StartService())
	return
}

func StopServiceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, mqttpubsub.StopService())
	return
}

func GetServiceStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, mqttpubsub.GetServiceStatus())
	return
}

///////////////CONVERSIONs OF HTTP BODY TO SPECIFIC STRUCTURES////////////////////////////

func convertHTTPBodyConfig(httpBody io.ReadCloser) (config_parser.Config, int, error) {
	body, err := io.ReadAll(httpBody)
	if err != nil {
		return config_parser.Config{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	var Config config_parser.Config
	err = json.Unmarshal(body, &Config)
	if err != nil {
		return config_parser.Config{}, http.StatusBadRequest, err
	}
	return Config, http.StatusOK, nil
}
