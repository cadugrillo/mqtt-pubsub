package configurator

import config_parser "mqtt-pubsub/modules/config-parser"

func GetConfig() config_parser.Config {
	return config_parser.GetConfig()
}

func SetConfig(ConfigFile config_parser.Config) string {
	err := config_parser.SetConfig(ConfigFile)
	if err != nil {
		return err.Error()
	}
	return "Configuration updated successfully"
}
