package config

import (
	"gopkg.in/yaml.v2"
)

//自定义配置，不能改结构体名
type Custom struct {
	Url  			string  `yaml:"url"`
	GetRoute  		string  `yaml:"get_route"`
	PostRoute  		string  `yaml:"post_route"`
	ContentType  	string  `yaml:"content_type"`
	HttpProxy  		string  `yaml:"http_proxy"`
}

func (h Custom) String() string {
	s, _ := yaml.Marshal(h)
	return string(s)
}