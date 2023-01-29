package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
)

type IOutputSink interface {
	Write(data string) (error)
}

type FileOutputSink struct {
	FilePath string
}

func getOutputSinks(config map[string]interface{}) []IOutputSink {

	var sinkMap = config["Sinks"].(map[string]interface{})
	var sinks = make([]IOutputSink, len(sinkMap))
	var counter = 0
	for k, v := range sinkMap {
		value := v.(interface{}).(string)
		if (strings.EqualFold(k, "OutputFile")) {
			sinks[counter] = &FileOutputSink {
				FilePath: value,
			}
		} else if (strings.EqualFold(k, "S3Bucket")) {
			// TODO
		} else {
			fmt.Println("Unknown sink type specified in config file: %s", k)
		}
	}
	return sinks
}

func (f FileOutputSink) Write(data string) (error) {
	return ioutil.WriteFile(f.FilePath, []byte(data), 0644)
}