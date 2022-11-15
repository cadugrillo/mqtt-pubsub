package configurator

import yaml_parser "mqtt-pubsub/modules/yaml-parser"

func GetConfig() yaml_parser.Config {
	return yaml_parser.GetConfig()
}

func SetConfig(ConfigFile yaml_parser.Config) string {
	err := yaml_parser.SetConfig(ConfigFile)
	if err != nil {
		return err.Error()
	}
	return "Configuration updated successfully"
}
