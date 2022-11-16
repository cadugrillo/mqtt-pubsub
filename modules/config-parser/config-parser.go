package config_parser

import (
	"encoding/json"
	"os"
)

type Config struct {
	ClientSub struct {
		ClientId           string `json:"ClientId"`
		ServerAddress      string `json:"ServerAddress"`
		Qos                int    `json:"Qos"`
		ConnectionTimeout  int    `json:"ConnectionTimeout"`
		WriteTimeout       int    `json:"WriteTimeout"`
		KeepAlive          int    `json:"KeepAlive"`
		PingTimeout        int    `json:"PingTimeout"`
		ConnectRetry       bool   `json:"ConnectRetry"`
		AutoConnect        bool   `json:"AutoConnect"`
		OrderMaters        bool   `json:"OrderMaters"`
		UserName           string `json:"UserName"`
		Password           string `json:"Password"`
		TlsConn            bool   `json:"TlsConn"`
		RootCA             string `json:"RootCA"`
		ClientKey          string `json:"ClientKey"`
		PrivateKey         string `json:"PrivateKey"`
		InsecureSkipVerify bool   `json:"InsecureSkipVerify"`
	} `json:"ClientSub"`
	ClientPub struct {
		ClientId           string `json:"ClientId"`
		ServerAddress      string `json:"ServerAddress"`
		Qos                int    `json:"Qos"`
		ConnectionTimeout  int    `json:"ConnectionTimeout"`
		WriteTimeout       int    `json:"WriteTimeout"`
		KeepAlive          int    `json:"KeepAlive"`
		PingTimeout        int    `json:"PingTimeout"`
		ConnectRetry       bool   `json:"ConnectRetry"`
		AutoConnect        bool   `json:"AutoConnect"`
		OrderMaters        bool   `json:"OrderMaters"`
		UserName           string `json:"UserName"`
		Password           string `json:"Password"`
		TlsConn            bool   `json:"TlsConn"`
		RootCA             string `json:"RootCA"`
		ClientKey          string `json:"ClientKey"`
		PrivateKey         string `json:"PrivateKey"`
		InsecureSkipVerify bool   `json:"InsecureSkipVerify"`
		TranslateTopic     bool   `json:"TranslateTopic"`
		PublishInterval    int    `json:"PublishInterval"`
	} `json:"ClientPub"`
	Logs struct {
		SubPayload bool `json:"SubPayload"`
		Debug      bool `json:"Debug"`
		Warning    bool `json:"Warning"`
		Error      bool `json:"Error"`
		Critical   bool `json:"Critical"`
	} `json:"Logs"`
	TopicsSub struct {
		Topic []string
	} `json:"TopicsSub"`
	TopicsPub struct {
		Topic []string
	} `json:"TopicsPub"`
}

func GetConfig() Config {
	f, err := os.Open("./config/config.json")
	if err != nil {
		return Config{}
	}
	defer f.Close()

	var cfg Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}
	}

	return cfg
}

func SetConfig(ConfigFile Config) error {
	f, err := os.Create("./config/config.json")
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	err = encoder.Encode(&ConfigFile)
	if err != nil {
		return err
	}
	return err

}

func LoadConfig() Config {
	cfg := GetConfig()

	err := os.WriteFile("./certs/pub/Certificate.crt", []byte(cfg.ClientPub.ClientKey), os.ModePerm)
	if err != nil {
		return Config{}
	}

	err = os.WriteFile("./certs/pub/PrivateKey.key", []byte(cfg.ClientPub.PrivateKey), os.ModePerm)
	if err != nil {
		return Config{}
	}

	err = os.WriteFile("./certs/pub/RootCA.pem", []byte(cfg.ClientPub.RootCA), os.ModePerm)
	if err != nil {
		return Config{}
	}

	err = os.WriteFile("./certs/sub/Certificate.crt", []byte(cfg.ClientSub.ClientKey), os.ModePerm)
	if err != nil {
		return Config{}
	}

	err = os.WriteFile("./certs/sub/PrivateKey.key", []byte(cfg.ClientSub.PrivateKey), os.ModePerm)
	if err != nil {
		return Config{}
	}

	err = os.WriteFile("./certs/sub/RootCA.pem", []byte(cfg.ClientSub.RootCA), os.ModePerm)
	if err != nil {
		return Config{}
	}

	return cfg
}
