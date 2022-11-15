package mqttpubsub

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	yaml_parser "mqtt-pubsub/modules/yaml-parser"

	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	once       sync.Once
	ConfigFile yaml_parser.Config
	b          Mqttbuffer
	PubConnOk  bool
	SubConnOk  bool
	ClientSub  mqtt.Client
	ClientPub  mqtt.Client
)

type Mqttbuffer struct {
	Buffer       [300000]Message
	ReadPointer  int
	WritePointer int
}

// Message
type Message struct {
	Duplicate bool
	Qos       byte
	Retained  bool
	Topic     string
	MessageID uint16
	Payload   string
	Ack       bool
}

func init() {
	once.Do(initialise)
}

func initialise() {
	b = NewMqttbuffer()
	PubConnOk = false
	SubConnOk = false
}

type handler struct {
	f bool
}

func NewHandler() *handler {
	var f bool
	return &handler{f: f}
}

func (o *handler) handle(_ mqtt.Client, msg mqtt.Message) {

	var recmsg Message
	recmsg.Duplicate = msg.Duplicate()
	recmsg.Qos = msg.Qos()
	recmsg.Retained = msg.Retained()
	recmsg.MessageID = msg.MessageID()
	recmsg.Topic = msg.Topic()
	recmsg.Payload = string(msg.Payload())

	AddMessage(recmsg)
}

func NewTLSConfig(rootCAPath string, clientKeyPath string, privateKeyPath string, insecureSkipVerify bool) *tls.Config {

	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(rootCAPath)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	cert, err := tls.LoadX509KeyPair(clientKeyPath, privateKeyPath)
	if err != nil {
		panic(err)
	}

	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	//fmt.Println(cert.Leaf)

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: insecureSkipVerify,
		Certificates:       []tls.Certificate{cert},
	}
}

func NewMqttbuffer() Mqttbuffer {
	b := Mqttbuffer{}
	return b
}

func GetReadPointer() int {
	return b.ReadPointer
}

func GetWritePointer() int {
	return b.WritePointer
}

func AddMessage(message Message) {
	if b.WritePointer == len(b.Buffer)-1 {
		b.Buffer[b.WritePointer] = message
		b.WritePointer = 0
		return
	}
	b.Buffer[b.WritePointer] = message
	b.WritePointer++
	return
}

func ReadMessage(index int) (Message, error) {
	if index < len(b.Buffer) {
		return b.Buffer[index], nil
	}
	msg := Message{}
	return msg, errors.New(fmt.Sprintf("Index %d greater then buffer size [%d]", index, len(b.Buffer)))
}

func NextMessage() {
	if b.ReadPointer == len(b.Buffer)-1 {
		b.ReadPointer = 0
		//return b.ReadPointer
		return
	}
	if b.ReadPointer != b.WritePointer {
		b.ReadPointer++
		//return b.ReadPointer
		return
	}
	fmt.Println("No new messages on the buffer")
	//return b.ReadPointer
}

func (b Mqttbuffer) NewMessage() bool {
	return b.WritePointer != b.ReadPointer
}

func Run() {

	ConfigFile = yaml_parser.LoadConfig()

	//logs
	if ConfigFile.Logs.Error {
		mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	}
	if ConfigFile.Logs.Critical {
		mqtt.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
	}
	if ConfigFile.Logs.Warning {
		mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	}
	if ConfigFile.Logs.Debug {
		mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	}

	h := NewHandler()

	optsSub := mqtt.NewClientOptions()
	optsSub.AddBroker(ConfigFile.ClientSub.ServerAddress)

	switch ConfigFile.ClientSub.TlsConn {
	case true:
		tlsSub := NewTLSConfig("./certs/sub/RootCA.pem", "./certs/sub/Certificate.crt", "./certs/sub/PrivateKey.key", ConfigFile.ClientSub.InsecureSkipVerify)
		optsSub.SetClientID(ConfigFile.ClientSub.ClientId).SetTLSConfig(tlsSub)
	case false:
		optsSub.SetClientID(ConfigFile.ClientSub.ClientId)
		optsSub.SetUsername(ConfigFile.ClientSub.UserName)
		optsSub.SetPassword(ConfigFile.ClientSub.Password)
	}

	optsSub.SetOrderMatters(ConfigFile.ClientSub.OrderMaters)                                      // Allow out of order messages (use this option unless in order delivery is essential)
	optsSub.ConnectTimeout = (time.Duration(ConfigFile.ClientSub.ConnectionTimeout) * time.Second) // Minimal delays on connect
	optsSub.WriteTimeout = (time.Duration(ConfigFile.ClientSub.WriteTimeout) * time.Second)        // Minimal delays on writes
	optsSub.KeepAlive = int64(ConfigFile.ClientSub.KeepAlive)                                      // Keepalive every 10 seconds so we quickly detect network outages
	optsSub.PingTimeout = (time.Duration(ConfigFile.ClientSub.PingTimeout) * time.Second)          // local broker so response should be quick
	optsSub.ConnectRetry = ConfigFile.ClientSub.ConnectRetry                                       // Automate connection management (will keep trying to connect and will reconnect if network drops)
	optsSub.AutoReconnect = ConfigFile.ClientSub.AutoConnect
	optsSub.DefaultPublishHandler = func(_ mqtt.Client, msg mqtt.Message) { fmt.Printf("SUB BROKER - UNEXPECTED : %s\n", msg) }
	optsSub.OnConnectionLost = func(cl mqtt.Client, err error) { fmt.Println("SUB BROKER - CONNECTION LOST") } // Log events

	optsSub.OnConnect = func(c mqtt.Client) {
		fmt.Println("SUB BROKER - CONNECTION STABLISHED")

		// Establish the subscription - doing this here means that it will happen every time a connection is established
		// (useful if opts.CleanSession is TRUE or the broker does not reliably store session data)
		for i := 0; i < len(ConfigFile.TopicsSub.Topic); i++ {
			t := c.Subscribe(ConfigFile.TopicsSub.Topic[i], byte(ConfigFile.ClientSub.Qos), h.handle)
			id := i

			// the connection handler is called in a goroutine so blocking here would not cause an issue. However as blocking
			// in other handlers does cause problems its best to just assume we should not block
			go func() {
				_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
				if t.Error() != nil {
					fmt.Printf("SUB BROKER - ERROR SUBSCRIBING TO : %s\n", t.Error())
				} else {
					fmt.Println("SUB BROKER - SUBSCRIBED TO : ", ConfigFile.TopicsSub.Topic[id])
				}
			}()
		}
	}

	optsSub.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) { fmt.Println("SUB BROKER - ATTEMPTING TO RECONNECT") }

	/////opts for Pub Broker
	optsPub := mqtt.NewClientOptions()
	optsPub.AddBroker(ConfigFile.ClientPub.ServerAddress)

	switch ConfigFile.ClientPub.TlsConn {
	case true:
		tlsPub := NewTLSConfig("./certs/pub/RootCA.pem", "./certs/pub/Certificate.crt", "./certs/pub/PrivateKey.key", ConfigFile.ClientPub.InsecureSkipVerify)
		tlsPub.MaxVersion = tls.VersionTLS12
		optsPub.SetClientID(ConfigFile.ClientPub.ClientId).SetTLSConfig(tlsPub)
	case false:
		optsPub.SetClientID(ConfigFile.ClientPub.ClientId)
		optsPub.SetUsername(ConfigFile.ClientPub.UserName)
		optsPub.SetPassword(ConfigFile.ClientPub.Password)
	}

	optsPub.SetOrderMatters(ConfigFile.ClientPub.OrderMaters)                                      // Allow out of order messages (use this option unless in order delivery is essential)
	optsPub.ConnectTimeout = (time.Duration(ConfigFile.ClientPub.ConnectionTimeout) * time.Second) // Minimal delays on connect
	optsPub.WriteTimeout = (time.Duration(ConfigFile.ClientPub.WriteTimeout) * time.Second)        // Minimal delays on writes
	optsPub.KeepAlive = int64(ConfigFile.ClientPub.KeepAlive)                                      // Keepalive every 10 seconds so we quickly detect network outages
	optsPub.PingTimeout = (time.Duration(ConfigFile.ClientPub.PingTimeout) * time.Second)          // local broker so response should be quick
	optsPub.ConnectRetry = ConfigFile.ClientPub.ConnectRetry                                       // Automate connection management (will keep trying to connect and will reconnect if network drops)
	optsPub.AutoReconnect = ConfigFile.ClientPub.AutoConnect
	optsPub.DefaultPublishHandler = func(_ mqtt.Client, msg mqtt.Message) { fmt.Printf("PUB BROKER - UNEXPECTED : %s\n", msg) }

	optsPub.OnConnectionLost = func(cl mqtt.Client, err error) {
		fmt.Println("PUB BROKER - CONNECTION LOST")
		PubConnOk = false
	}

	optsPub.OnConnect = func(c mqtt.Client) {
		fmt.Println("PUB BROKER - CONNECTION STABLISHED")
		PubConnOk = true
	}

	optsPub.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) { fmt.Println("PUB BROKER - ATTEMPTING TO RECONNECT") }

	//
	// Connect to the SUB broker
	//
	ClientSub = mqtt.NewClient(optsSub)

	// If using QOS2 and CleanSession = FALSE then messages may be transmitted to us before the subscribe completes.
	// Adding routes prior to connecting is a way of ensuring that these messages are processed
	for i := 0; i < len(ConfigFile.TopicsSub.Topic); i++ {
		ClientSub.AddRoute(ConfigFile.TopicsSub.Topic[i], h.handle)
	}

	if tokenSub := ClientSub.Connect(); tokenSub.Wait() && tokenSub.Error() != nil {
		panic(tokenSub.Error())
	}
	fmt.Println("SUB BROKER  - CONNECTION IS UP")

	//
	//connect to PUB broker
	//
	ClientPub = mqtt.NewClient(optsPub)

	if tokenPub := ClientPub.Connect(); tokenPub.Wait() && tokenPub.Error() != nil {
		panic(tokenPub.Error())
	}
	fmt.Println("PUB BROKER  - CONNECTION IS UP")

	go func() {
		for {
			if b.NewMessage() && PubConnOk {
				msg, err := ReadMessage(GetReadPointer())
				if err != nil {
					panic(err.Error())
				}
				if ConfigFile.Logs.SubPayload {
					fmt.Println(msg.Payload)
					fmt.Println(GetWritePointer())
					fmt.Println(GetReadPointer())
				}
				switch ConfigFile.ClientPub.TranslateTopic {
				case false:
					ClientPub.Publish(msg.Topic, msg.Qos, msg.Retained, msg.Payload)
				case true:
					for i := 0; i < len(ConfigFile.TopicsSub.Topic); i++ {
						if ConfigFile.TopicsSub.Topic[i] == msg.Topic {
							TranslatedTopic := ConfigFile.TopicsPub.Topic[i]
							ClientPub.Publish(TranslatedTopic, msg.Qos, msg.Retained, msg.Payload)
							break
						}
					}
				}
				NextMessage()
				//b.ReadPointer = b.NextMessage()
			}
			time.Sleep(time.Duration(ConfigFile.ClientPub.PublishInterval) * time.Millisecond)
		}
	}()

	// Messages will be delivered asynchronously so we just need to wait for a signal to shutdown
	// sig := make(chan os.Signal, 2)
	// signal.Notify(sig, os.Interrupt)
	// signal.Notify(sig, syscall.SIGTERM)

	// <-sig
	// fmt.Println("signal caught - exiting")
	// ClientSub.Disconnect(1000)
	// ClientPub.Disconnect(1000)
	// fmt.Println("shutdown complete")

}
