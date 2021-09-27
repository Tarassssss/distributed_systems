package basement

import (
	"io/ioutil"
	"net/http"
)

const MsgServicesNotResponding = "all logging services are not responding"
const Localhost = "http://127.0.0.1"
const FacadeAddress = ":23000"
var LoggingServiceAddress = [...]string{":23001", ":23002", ":23003"}
var HazelcastAddress = [...]string{":5701", ":5702", ":5703"}
//const MessagesServiceAddress = ":23102"
var MessagesServiceAddress = [...]string{":23102", ":23103"}

var LoggerPortsSize = len(LoggingServiceAddress)

type RequestInfo struct {
	Id  string
	Msg string
}

type RequestData struct {
	Msg string
}

func GetData(serverAddress string) (string, error)  {
	loggingResp, err := http.Get(Localhost + serverAddress)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(loggingResp.Body)
	return string(body), nil
}