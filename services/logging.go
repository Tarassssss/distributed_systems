package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"basement"
	"net"
	"context"
	"github.com/hazelcast/hazelcast-go-client"
)

var messagesByKeys *hazelcast.Map
var ctx context.Context
func main() {
	port, hazelcastPort := findAvailablePort()
	messagesByKeys, ctx = getHazelcastMap(hazelcastPort)
	_ = http.ListenAndServe(port, &LoggingResponse{})
}

type LoggingResponse struct{}

func (m *LoggingResponse) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
    		valuePairs, _ := messagesByKeys.GetEntrySet(ctx)
			msg := ""
			if len(valuePairs) > 0 {
				for i := 0; i < len(valuePairs); i++ {
					msg += valuePairs[i].Key.(string) + " : " + valuePairs[i].Value.(string) + "\n"
				}
			} else {
				msg = "There are no log messages"
			}
			_, _ = writer.Write([]byte(msg))	
	} else if request.Method == "POST" {
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
		_, _ = messagesByKeys.Put(ctx, requestInfo.Id, requestInfo.Msg)
	}
}

func getHazelcastMap(port string) (*hazelcast.Map, context.Context) {
	ctx := context.Background()
	config := hazelcast.NewConfig()
	config.Cluster.Network.SetAddresses("localhost" + port)
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	if err != nil {
		fmt.Println(err)
	}
	logMap, err := client.GetMap(ctx, "log")
	if err != nil{
		fmt.Println("Error is", err)
	}
	return logMap, ctx
}

func findAvailablePort() (string, string) {
	for i := 0; i < len(basement.LoggingServiceAddress); i++ {
		listener, err := net.Listen("tcp", basement.LoggingServiceAddress[i])
		if err != nil {
			fmt.Println(basement.LoggingServiceAddress[i], " is already taken, try next one")
		} else {
			_ = listener.Close()
			return basement.LoggingServiceAddress[i], basement.HazelcastAddress[i]
		}
	}
	panic("All ports are unavailable")
}