package boomer

import (
	caseConfig "falcon/slave/config"
	siteConfig "falcon/libs/config"
	"gopkg.in/yaml.v2"
	"time"
)

type BoomerConfig struct {
	Cases  CasesConf  	`yaml:"cases"`
	caseConfig.Custom 	`yaml:"custom"`
}

type SlaveConf struct {
	MinWait             int64  `yaml:"min_wait"`
	MaxWait             int64  `yaml:"max_wait"`
	StopTimeout         int64  `yaml:"stop_timeout"`
	ReportInterval 		int64  `yaml:"report_interval"`
	MonitorInterval 	int64  `yaml:"monitor_interval"`
	HeartbeatInterval   int64  `yaml:"heartbeat_interval"`
	RendezvousInterval  int64  `yaml:"rendezvous_interval"`
	RunTimes            int64  `yaml:"run_times"`
	RandomMode          bool   `yaml:"random_mode"`
}

type WeightConf struct {
	Type   				string `yaml:"type"`
	Fn     				string `yaml:"fn"`
	Name   				string `yaml:"name"`
	Enable 				bool   `yaml:"enable"`
	Weight 				int64  `yaml:"weight"`
}

type SlaveConfig struct {
	BaseConfig 				*BoomerConfig
	masterConfig 			*BoomerConfig
	mutexInterval			time.Duration
	checkInterval			time.Duration
	breakMillisecond		int64
	restartWait				int
}

type CasesConf []WeightConf

var Conf *SlaveConfig


//读取配置
func LoadBaseConfig() *BoomerConfig {
	var conf BoomerConfig
	if err := siteConfig.LoadConfig(Param.ConfigPath, &conf); err != nil {
		control.log.error("Load base config fail,error:%s", err)
	}

	if Conf != nil {
		if Conf.masterConfig != nil {
			conf = *Conf.masterConfig
		}
	}
	return &conf
}

//读取Data配置
func LoadDataConfig(data []byte) *BoomerConfig {
	var conf BoomerConfig
	if err := siteConfig.LoadYamlData(data, &conf); err != nil {
		control.log.error("Load config data fail,error:%s", err)
	}
	return &conf
}

func (h BoomerConfig) String() string {
	s, _ := yaml.Marshal(h)
	return string(s)
}

func (h SlaveConf) String() string {
	s, _ := yaml.Marshal(h)
	return string(s)
}

func (h WeightConf) String() string {
	s, _ := yaml.Marshal(h)
	return string(s)
}