package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type DataSource interface {
	GetData() (string, error)
}

type WebSocketDataSource struct {
	Config map[string]interface{}
}

type FileDataSource struct {
	File string
}

func getDataSource(config map[string]interface{}) DataSource {
	isFileInput := false
	if val, ok := config["IsFileInput"]; ok {
		isFileInput = val.(bool)
	}

	if isFileInput {
		return &FileDataSource{
			File: config["DebugFilePath"].(string),
		}
	} else {
		return &WebSocketDataSource{
			Config: config,
		}
	}
}

// Specific to WebSocket data source
const okMessage = "\"s\":\"ok\""

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ws WebSocketDataSource) GetData() (string, error) {
	wsHeader := http.Header{}

	headers := ws.Config["HttpRequestHeader"].(map[string]interface{})
	for k, v := range headers {
		value := v.(interface{}).(string)
		wsHeader.Add(k, value)
	}

	serverAddress := ws.Config["WebSocketServer"].(string)
	c, _, err := websocket.DefaultDialer.Dial(serverAddress, wsHeader)
	if err != nil {
		return "", err
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		return "", err
	}
	fmt.Println(string(message))

	for _, msg := range [2]string{ws.Config["InitialMsg"].(string), ws.Config["GetSelectionsMsg"].(string)} {
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			return "", err
		}
	}

	_, message, err = c.ReadMessage()
	if err != nil {
		// handle error
		return "", err
	} else if !strings.Contains(string(message), okMessage) {
		return "", errors.New("Unexpected response" + string(message))
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
	return sb.String(), nil
}

func (f FileDataSource) GetData() (string, error) {
	data, err := ioutil.ReadFile(f.File)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(data), nil
}