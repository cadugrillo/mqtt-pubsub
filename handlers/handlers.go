package handlers

import (
	"encoding/json"
	"io"
	"mqtt-pubsub/modules/configurator"
	yaml_parser "mqtt-pubsub/modules/yaml-parser"
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

///////////////CONVERSIONs OF HTTP BODY TO SPECIFIC STRUCTURES////////////////////////////

func convertHTTPBodyConfig(httpBody io.ReadCloser) (yaml_parser.Config, int, error) {
	body, err := io.ReadAll(httpBody)
	if err != nil {
		return yaml_parser.Config{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	var Config yaml_parser.Config
	err = json.Unmarshal(body, &Config)
	if err != nil {
		return yaml_parser.Config{}, http.StatusBadRequest, err
	}
	return Config, http.StatusOK, nil
}
