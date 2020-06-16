package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/gorilla/websocket"
)

// type SupercontestRaw struct {
// 	Date string
// 	RawData string
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const okMessage = "\"s\":\"ok\""
const testFilename = "output_test.json"

var globalConfiguration map[string]interface{}

func main() {
	var configFile = "config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	fmt.Println("Using config file " + configFile)
	configuration, err := parseConfigFile(configFile)
	globalConfiguration = configuration
	if err != nil {
		fmt.Println("Error parsing configuration file: " + err.Error())
	}

	fmt.Println()

	//success, data := getDataFromFile()
	success, data := getSupercontestData()
	if !success {
		fmt.Println(data)
		return
	}

	parseJsonData(data)

	// err := WriteToFile(testFilename, data)
	// if err != nil {
	// 	fmt.Println(err)
	// }
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

func getDataFromFile() (bool, string) {
	data, err := ioutil.ReadFile(testFilename)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	return true, string(data)
}

func parseJsonData(json string) {
	data := []byte(json)
	//var counter int = 0
	jsonparser.ObjectEach(data, func(key []byte, value []byte, datatype jsonparser.ValueType, offset int) error {
		fmt.Println(string(value))
		return nil
	}, "d", "b", "d", "data")
}

func getSupercontestData() (bool, string) {

	wsHeader := http.Header{}

	headers := globalConfiguration["HttpRequestHeader"].(map[string]interface{}) //["data"].(interface{})
	for k, v := range headers {
		value := v.(interface{}).(string)
		//fmt.Println(k + " " + value)
		wsHeader.Add(k, value)
	}

	serverAddress := globalConfiguration["WebSocketServer"].(string)
	c, _, err := websocket.DefaultDialer.Dial(serverAddress, wsHeader)
	if err != nil {
		fmt.Println(err.Error())
		return false, err.Error()
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		// handle error
		return false, err.Error()
	}
	fmt.Println(string(message))

	for _, msg := range [2]string{globalConfiguration["InitialMsg"].(string), globalConfiguration["GetSelectionsMsg"].(string)} {
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			return false, err.Error()
		}
	}

	_, message, err = c.ReadMessage()
	if err != nil {
		// handle error
		return false, err.Error()
	} else if !strings.Contains(string(message), okMessage) {
		return false, "Unexpected response" + string(message)
	}

	var sb strings.Builder
	// receive message
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			// handle error
			fmt.Println(err.Error())
			break
		}

		strMessage := string(message)
		if strMessage == "16" {
			continue
		} else if strings.Contains(strMessage, okMessage) {
			fmt.Println("done.")
			break
		}

		sb.WriteString(strMessage)
	}
	return true, sb.String()
}
