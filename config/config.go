package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/toxuin/alarmserver/servers/hikvision"
)

type Config struct {
	Debug     bool            `json:"debug"`
	Mqtt      MqttConfig      `json:"mqtt"`
	Webhooks  WebhooksConfig  `json:"webhooks"`
	Hikvision HikvisionConfig `json:"hikvision"`
}

type MqttConfig struct {
	Enabled   bool   `json:"enabled"`
	Server    string `json:"server"`
	Port      string `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	TopicRoot string `json:"topicRoot"`
}

type WebhooksConfig struct {
	Enabled bool            `json:"enabled"`
	Items   []WebhookConfig `json:"items"`
	Urls    []string        `json:"urls"`
}

type WebhookConfig struct {
	Url          string   `json:"url"`
	Method       string   `json:"method"`
	Headers      []string `json:"headers"`
	BodyTemplate string   `json:"bodyTemplate"`
}


type HikvisionConfig struct {
	Enabled bool                  `json:"enabled"`
	Cams    []hikvision.HikCamera `json:"cams"`
}


func (c *Config) SetDefaults() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")

	viper.SetDefault("debug", false)
	viper.SetDefault("mqtt.port", 1883)
	viper.SetDefault("mqtt.topicRoot", "camera-alerts")
	viper.SetDefault("mqtt.server", "mqtt.example.com")
	viper.SetDefault("mqtt.username", "anonymous")
	viper.SetDefault("mqtt.password", "")
	viper.SetDefault("hikvision.enabled", false)

	_ = viper.BindEnv("debug", "DEBUG")
	_ = viper.BindEnv("mqtt.port", "MQTT_PORT")
	_ = viper.BindEnv("mqtt.topicRoot", "MQTT_TOPIC_ROOT")
	_ = viper.BindEnv("mqtt.server", "MQTT_SERVER")
	_ = viper.BindEnv("mqtt.username", "MQTT_USERNAME")
	_ = viper.BindEnv("mqtt.password", "MQTT_PASSWORD")
	_ = viper.BindEnv("hikvision.enabled", "HIKVISION_ENABLED")
	_ = viper.BindEnv("hikvision.cams", "HIKVISION_CAMS")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, writing default config...")
			err := viper.SafeWriteConfig()
			if err != nil {
				panic(fmt.Errorf("error saving default config file: %s \n", err))
			}
		} else {
			panic(fmt.Errorf("error reading config file: %s \n", err))
		}
	}
}

func (c *Config) Load() *Config {
	myConfig := Config{
		Debug:     viper.GetBool("debug"),
		Mqtt:      MqttConfig{},
		Webhooks:  WebhooksConfig{},
		Hikvision: HikvisionConfig{
			Enabled: viper.GetBool("hikvision.enabled"),
		},
	}

	if viper.IsSet("mqtt") {
		err := viper.Sub("mqtt").Unmarshal(&myConfig.Mqtt)
		if err != nil {
			panic(fmt.Errorf("unable to decode mqtt config, %v", err))
		}
	}
	if viper.IsSet("webhooks") {
		err := viper.Sub("webhooks").Unmarshal(&myConfig.Webhooks)
		if err != nil {
			panic(fmt.Errorf("unable to decode webhooks config, %v", err))
		}
	}

	if !myConfig.Mqtt.Enabled && !myConfig.Webhooks.Enabled {
		panic("Both MQTT and Webhook buses are disabled. Nothing to do!")
	}


	if viper.IsSet("hikvision.cams") {
		hikvisionCamsConfig := viper.Sub("hikvision.cams")
		if hikvisionCamsConfig != nil {
			camConfigs := viper.GetStringMapString("hikvision.cams")

			for camName := range camConfigs {
				camConfig := viper.Sub("hikvision.cams." + camName)
				// CONSTRUCT CAMERA URL
				url := ""
				if camConfig.GetBool("https") {
					url += "https://"
				} else {
					url += "http://"
				}
				url += camConfig.GetString("address") + "/ISAPI/"

				camera := hikvision.HikCamera{
					Name:     camName,
					Url:      url,
					Username: camConfig.GetString("username"),
					Password: camConfig.GetString("password"),
				}
				if camConfig.GetBool("rawTcp") {
					camera.BrokenHttp = true
				}
				if myConfig.Debug {
					fmt.Printf("Added Hikvision camera:\n"+
						"  name: %s \n"+
						"  url: %s \n"+
						"  username: %s \n"+
						"  password set: %t\n"+
						"  rawRcp: %t\n",
						camera.Name,
						camera.Url,
						camera.Username,
						camera.Password != "",
						camera.BrokenHttp,
					)
				}

				myConfig.Hikvision.Cams = append(myConfig.Hikvision.Cams, camera)
			}
		}
	}
	return &myConfig
}

func (c *Config) Printout() {
	fmt.Printf("CONFIG:\n"+
		"  SERVER: Hikvision - enabled: %t\n"+
		"    camera count: %d\n"+
		"  BUS: MQTT - enabled: %t\n"+
		"    port: %s\n"+
		"    topicRoot: %s\n"+
		"    server: %s\n"+
		"    username: %s\n"+
		"    password set: %t\n"+
		"  BUS: Webhooks - enabled: %t\n"+
		"    count: %d\n",
		c.Hikvision.Enabled,
		len(c.Hikvision.Cams),
		c.Mqtt.Enabled,
		c.Mqtt.Port,
		c.Mqtt.TopicRoot,
		c.Mqtt.Server,
		c.Mqtt.Username,
		c.Mqtt.Password != "",
		c.Webhooks.Enabled,
		len(c.Webhooks.Items)+len(c.Webhooks.Urls),
	)
}