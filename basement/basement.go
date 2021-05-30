package basement

import (
	"io/ioutil"
	"net/http"
)

const Localhost = "http://127.0.0.1"
const FacadeAddress = ":23000"
const LoggingServiceAddress = ":23001"
const MessagesServiceAddress = ":23002"

type RequestInfo struct {
	Id  string
	Msg string
}

type RequestData struct {
	Msg string
}

func GetData(serverAddress string) string {
	loggingResp, err := http.Get(Localhost + serverAddress)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(loggingResp.Body)
	return string(body)
}