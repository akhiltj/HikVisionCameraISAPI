package main

import (
	"fmt"
	"github.com/toxuin/alarmserver/buses/mqtt"
	"github.com/toxuin/alarmserver/buses/webhooks"
	conf "github.com/toxuin/alarmserver/config"
	"github.com/toxuin/alarmserver/servers/hikvision"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var config *conf.Config

func init() {
	config.SetDefaults()
}

func main() {
	config = config.Load()
	fmt.Println("STARTING...")
	if config.Debug {
		config.Printout()
	}

	processesWaitGroup := sync.WaitGroup{}

	// INIT BUSES
	mqttBus := mqtt.Bus{Debug: config.Debug}
	if config.Mqtt.Enabled {
		mqttBus.Initialize(config.Mqtt)
		if config.Debug {
			fmt.Println("MQTT BUS INITIALIZED")
		}
	}

	webhookBus := webhooks.Bus{Debug: config.Debug}
	if config.Webhooks.Enabled {
		webhookBus.Initialize(config.Webhooks)
		if config.Debug {
			fmt.Println("WEBHOOK BUS INITIALIZED")
		}
	}

	messageHandler := func(cameraName string, eventType string, extra string) {
		if config.Mqtt.Enabled {
			mqttBus.SendMessage(config.Mqtt.TopicRoot+"/"+cameraName+"/"+eventType, extra)
			mqttBus.SendMessage("test", "extra")
		}
		if config.Webhooks.Enabled {
			webhookBus.SendMessage(cameraName, eventType, extra)
		}
	}


	if config.Hikvision.Enabled {
		// START HIKVISION ALARM SERVER
		hikvisionServer := hikvision.Server{
			Debug:          config.Debug,
			WaitGroup:      &processesWaitGroup,
			Cameras:        &config.Hikvision.Cams,
			MessageHandler: messageHandler,
		}
		hikvisionServer.Start()
		if config.Debug {
			fmt.Println("STARTED HIKVISION SERVER")
		}
	}

	processesWaitGroup.Wait()

	// START INFINITE LOOP WAITING FOR SERVERS
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
