package main

import (
	"log"
	//"splunk.com/agent/plugin"
	//"splunk.com/agent/context"
	//"github.com/rcrowley/go-metrics"
	//"fmt"
	//"io/ioutil"
	"splunk.com/agent/util"
	"net/http"
	"bytes"
	"encoding/json"
)

func main() {
	//MakeRequest()
	Send2ssc("srini")
}
func MakeRequest() {
	message := map[string]interface{}{
		"hello": "Srinivas",
		"life":  42,
		"embedded": map[string]string{
			"yes": "of course!",
		},
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	//#resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(bytesRepresentation))
	resp, err := http.Post("http://i-0447166d3d40f72ed.ec2.splunkit.io:8088/services/collector", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["data"])
}
func Send2ssc(msg  string) (int ) {

	//logger := util.NewLogger()

	//curl -k http://localhost:8088/services/collector -H "Authorization: Splunk f4d49a2b-9b70-488e-9789-86c18557515a"
	// -H "Content-Type: application/json" -d '{"event":"testevent1101","sourcetype":"test-auto-extract"}'
	//url := "http://localhost:8088/services/collector"
	url := "http://i-0447166d3d40f72ed.ec2.splunkit.io:8088/services/collector"

	//url = "https://api.playground.splunkbeta.com:443/hyuan/ingest/v1/events/"

	//logMsg := fmt.Sprintf("Send2ssc - sending event to %s ", url)

	//logger.Debug(logMsg)

	//token := "Splunk f4d49a2b-9b70-488e-9789-86c18557515a";
	token := "Splunk 8d51f577-b1b9-414e-8e26-579a08ee6d90";
	//var jsonString = "{\"event\": \"test\"}"
	//jsonString := fmt.Sprintf("{\"event\": \"%s\"}", msg)

	var mapjson = map[string]string {
		"event":msg,
		"sourcetype":"test-auto-extract",
	}
	//mapjson ["event"] = msg;
	jsonBufferVal, err := json.Marshal(mapjson)
	if err != nil {
		log.Printf("error happend %s", err.Error())
	}


	//jsonString := msg;
	//logMsg = fmt.Sprintf(jsonString)
	//logger.Debug(logMsg)
	//var jsonBytes = []byte(jsonString)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBufferVal))

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()


	//logMsg = fmt.Sprintln("Response Status:", resp.Status)
	//logger.Debug(logMsg)
	//logMsg = fmt.Sprintln("Response Headers:", resp.Header)
	//logger.Debug(logMsg)
	//body, _ := ioutil.ReadAll(resp.Body)
	//logMsg = fmt.Sprintln("Response Body:", string(body))
	//logger.Debug(logMsg)
	return resp.StatusCode
}
