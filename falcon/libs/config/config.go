package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
)


func LoadConfig(path string, conf interface{}) error {
	err := checkConfType(conf)
	if err != nil {
		return err
	}

	var data []byte
	data = []byte(os.Getenv("LOCUST_CONFIG"))
	if len(data) == 0 {
		_, err := os.Stat(path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		data, err = ioutil.ReadAll(file)
		if err != nil {
			return err
		}
	}

	err = LoadYamlData(data, conf)
	if err != nil {
		return err
	}
	return nil
}

func LoadYamlData(data []byte, conf interface{}) error {
	err := checkConfType(conf)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	return nil
}

func checkConfType(conf interface{}) error {
	if reflect.TypeOf(conf).Kind().String() != "ptr" {
		return errors.New("conf type is not ptr")
	}
	if reflect.TypeOf(conf).Elem().Kind().String() != "struct" {
		return errors.New("*conf type is not struct")
	}
	return nil
}
