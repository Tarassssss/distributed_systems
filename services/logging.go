package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"basement"
)

var messagesByKeys = make(map[string]string)

func main() {
	panic(http.ListenAndServe(basement.LoggingServiceAddress, &LoggingResponse{}))
}

type LoggingResponse struct{}

func (m *LoggingResponse) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		for _, value := range messagesByKeys {
    		//fmt.Println("Key:", key, "Value:", value)
    		logs, _ := json.Marshal(value)
    		_, _ = writer.Write(logs)
		}		
	} else if request.Method == "POST"{
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		var requestInfo basement.RequestInfo
		if err = json.Unmarshal(body, &requestInfo); err != nil {
			panic(err)
		} else {
			fmt.Println("requestInfo.Msg => " + requestInfo.Msg)
			fmt.Println("requestInfo.Id => " + requestInfo.Id)
		}
		messagesByKeys[requestInfo.Id] = requestInfo.Msg
	}
}