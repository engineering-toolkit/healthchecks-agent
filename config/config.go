package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represent configuration attribute
type Config struct {
	HC struct {
		ServerURL string `json:"ServerURL" yaml:"ServerURL"`
	} `json:"HC" yaml:"HC"`

	Log struct {
		Level string `json:"Level" yaml:"Level"`
	} `json:"Log" yaml:"Log"`

	HttpClient struct {
		RetryCount       int `json:"RetryCount" yaml:"RetryCount"`
		RetryWaitTime    int `json:"RetryWaitTime" yaml:"RetryWaitTime"`
		RetryMaxWaitTime int `json:"RetryMaxWaitTime" yaml:"RetryMaxWaitTime"`
	} `json:"HttpClient" yaml:"HttpClient"`

	ServiceGroups []ServiceGroup `json:"ServiceGroups" yaml:"ServiceGroups"`

	ClientChecks []ClientCheck `json:"ClientChecks" yaml:"ClientChecks"`

	SGClientChecks map[string][]ClientCheck
}

// ServiceGroup represent ServiceGroup config attribute
type ServiceGroup struct {
	Enable      bool     `json:"Enable" yaml:"Enable"`
	Name        string   `json:"Name" yaml:"Name"`
	ConfigPaths []string `json:"ConfigPaths" yaml:"ConfigPaths"`
}

// SGClientCheck represent SGClientCheck config attribute
type SGClientCheck struct {
	ClientChecks []ClientCheck `json:"ClientChecks" yaml:"ClientChecks"`
}

// ClientCheck represent ClientCheck config attribute
type ClientCheck struct {
	Enable      bool   `json:"Enable" yaml:"Enable"`
	ServiceName string `json:"ServiceName" yaml:"ServiceName"`
	ClientUUID  string `json:"ClientUUID" yaml:"ClientUUID"`
	Customs     struct {
		HC struct {
			ServerURL      string `json:"ServerURL" yaml:"ServerURL"`
			SendSuccessLog bool   `json:"SendSuccessLog" yaml:"SendSuccessLog"`
			SendFailLog    bool   `json:"SendFailLog" yaml:"SendFailLog"`
		} `json:"HC" yaml:"HC"`
	} `json:"Customs" yaml:"Customs"`

	CheckSchedule struct {
		// ScheduleType: FIXED, CRON
		ScheduleType  string `json:"ScheduleType" yaml:"ScheduleType"`
		CheckInterval int    `json:"CheckInterval" yaml:"CheckInterval"`
		// CheckDuration: SECOND, MINUTE, HOUR
		CheckDuration string `json:"CheckDuration" yaml:"CheckDuration"`
		// CheckCron: * * * * *
		CheckCron string `json:"CheckCron" yaml:"CheckCron"`
	} `json:"CheckSchedule" yaml:"CheckSchedule"`

	// ServiceType: HTTP, Port, Database, OS
	ServiceType string `json:"ServiceType" yaml:"ServiceType"`
	// ServiceProviderType: HTTP.Default, Port.TCP, Database.MySQL, OS.service
	ServiceProviderType string `json:"ServiceProviderType" yaml:"ServiceProviderType"`
	// ServiceProviderProperties custom ServiceProvider Properties
	ServiceProviderProperties interface{}
}

// ConfigHandler represent Config Handler
type ConfigHandler struct {
	handler      *viper.Viper
	ClientChecks map[string][]ClientCheck
}

// NewConfigHandler create New ConfigHandler
func NewConfigHandler(testingMode bool) (ConfigHandler, error) {
	ch := ConfigHandler{}

	ch.handler = viper.New()
	ch.handler.SetConfigType("yaml")
	ch.handler.SetConfigName("config")
	if testingMode {
		ch.handler.AddConfigPath("../conf/test-data")
	} else {
		ch.handler.AddConfigPath("./conf")
		ch.handler.AddConfigPath("./")
	}
	ch.handler.ReadInConfig()

	return ch, nil
}

// Read configuration
func (ch *ConfigHandler) Read() (Config, error) {
	var c Config
	err := ch.handler.Unmarshal(&c)

	ch.ClientChecks = make(map[string][]ClientCheck, len(c.ServiceGroups)+1)
	ch.ClientChecks["Root"] = c.ClientChecks

	// Service Group
	if len(c.ServiceGroups) > 0 {
		for _, sg := range c.ServiceGroups {
			if sg.Enable {
				for fdx, cfgPath := range sg.ConfigPaths {
					cc, err := ch._readServiceGroup(cfgPath)
					if err != nil {
						fmt.Printf("Failed load service group: %s\n", cfgPath)
					}
					if len(cc.ClientChecks) > 0 {
						ch.ClientChecks[fmt.Sprintf("'%s-%d: %s'", sg.Name, fdx, cfgPath)] = cc.ClientChecks
					}
				}
			}
		}
	}
	c.SGClientChecks = ch.ClientChecks

	return c, err
}

// _readServiceGroup read clientcheck in service group
func (ch *ConfigHandler) _readServiceGroup(filePath string) (SGClientCheck, error) {
	var aCC SGClientCheck

	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigFile(filePath)
	vp.ReadInConfig()

	err := vp.Unmarshal(&aCC)

	return aCC, err
}
