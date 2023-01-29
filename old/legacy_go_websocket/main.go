package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/buger/jsonparser"
)

func main() {
	var configFile = "config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	fmt.Println("Using config file " + configFile)
	configuration, err := parseConfigFile(configFile)
	if err != nil {
		fmt.Println("Error parsing configuration file: " + err.Error())
		return
	}
	sinks := getOutputSinks(configuration)

	dataSource := getDataSource(configuration)
	data, err := dataSource.GetData()
	if err != nil {
		fmt.Println("Error getting data: " + err.Error())
		return
	}

	if (len(sinks) > 0) {
		fmt.Println(data)
		for _, sink := range sinks {
			if (sink != nil) {
				 sink.Write(data)
			}
		}
	}
}

func parseConfigFile(file string) (map[string]interface{}, error) {
	configData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(configData), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func parseJSONData(json string) {
	data := []byte(json)
	//var counter int = 0
	jsonparser.ObjectEach(data, func(key []byte, value []byte, datatype jsonparser.ValueType, offset int) error {
		fmt.Println(string(value))
		return nil
	}, "d", "b", "d", "data")
}
