package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

//const addr = "localhost:4242"
const addr = "10.140.127.166:4242"

const message = "HACK: This msg is sent by QUIC protocol";

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {
	for i := 1;  i<=5; i++ {
    fmt.Printf("quicClient: Iteration = %d \n", i)
		err := clientMain()
		if err != nil {
			fmt.Printf("error !!! \n")
		}
		time.Sleep(1 * time.Second)
	}
}

// Start a server that echos all data on the first stream opened by the client
func echoServer() error {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	sess, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}
	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}
	// Echo through the loggingWriter
	_, err = io.Copy(loggingWriter{stream}, stream)
	return err
}

func clientMain() error {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return err
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Client: Sending '%s'\n", message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		return err
	}
/*
	buf := make([]byte, len(message))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		return err
	}
	fmt.Printf("Client: Got '%s'\n", buf)
*/
	return nil
}

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))

	// HEC sink goes here
	Send2Splunk(string(b));

	return w.Writer.Write(b)
}

func Send2Splunk(msg  string) (int ) {

	//logger := util.NewLogger()

	//curl -k http://localhost:8088/services/collector -H "Authorization: Splunk f4d49a2b-9b70-488e-9789-86c18557515a"
	// -H "Content-Type: application/json" -d '{"event":"testevent1101","sourcetype":"test-auto-extract"}'
	url := "http://i-0447166d3d40f72ed.ec2.splunkit.io:8088/services/collector"

	//url = "https://api.playground.splunkbeta.com:443/hyuan/ingest/v1/events/"

	logMsg := fmt.Sprintf("Send2Splunk - sending event to %s ", url)

	fmt.Printf(logMsg)

	//token := "Bearer eyJraWQiOiI3cXV0SjFleUR6V2lOeGlTbktsakZHRWhmU0VzWFlMQWt0NUVNbzJaNFk4IiwiYWxnIjoiUlMyNTYifQ.eyJ2ZXIiOjEsImp0aSI6IkFULnZuTHpCVU01YmcxMDZFSmplamM3THNEX2thR09yT0FSQS1TSFdfc0Z2Z3ciLCJpc3MiOiJodHRwczovL3NwbHVuay1jaWFtLm9rdGEuY29tL29hdXRoMi9kZWZhdWx0IiwiYXVkIjoiYXBpOi8vZGVmYXVsdCIsImlhdCI6MTUyOTU0MjI4MSwiZXhwIjoxNTI5NTg1NDgxLCJjaWQiOiIwb2FwYmcyem1MYW1wV2daNDJwNiIsInVpZCI6IjAwdTFlMjMxMjJ1QnFjTGtYMnA3Iiwic2NwIjpbInByb2ZpbGUiLCJvcGVuaWQiLCJlbWFpbCJdLCJzdWIiOiJoeXVhbkBzcGx1bmsuY29tIn0.GLk7zPP8I15tE_FSPQA8yMVyJqjYUmITXVxfVh-t6HX8Qimk8hi8O1G1U-vhkOt566EWiPwpLFevo0m9M6UlOgBI-4RpILGWdFFcNvcWgYSh5eODgIdUqcb8OkxwM1R6KMjCI75qKI36oCYgOV0Ffw_pBUhT7XVYOgxy93Gdk8UK7pGIuNNXv_HCcxFOfbbtcZ-3tlM7d13vI4BivHHG2uYdDbWGsph1m5qLiuSm9Xscozd-hMjgX7FzruEwlTgmRGzu0Cci_wNefmO4HzCGEec6AML3KPeaFwL02DXmFKwioyCIAQJP9OGaRsmhnyiRUOQkxqLSPifWmE-Gc6bpMg"
	token := "Splunk 8d51f577-b1b9-414e-8e26-579a08ee6d90";
	//var jsonString = "{\"event\": \"test\"}"
	jsonString := fmt.Sprintf("{\"event\": \"%s\"}", msg)

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
	logMsg = fmt.Sprintf(jsonString)
	fmt.Printf(logMsg)
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

	logMsg = fmt.Sprintln("Response Status:", resp.Status)
	fmt.Printf(logMsg)
	logMsg = fmt.Sprintln("Response Headers:", resp.Header)
	fmt.Printf(logMsg)
	body, _ := ioutil.ReadAll(resp.Body)
	logMsg = fmt.Sprintln("Response Body:", string(body))
	fmt.Printf(logMsg)
	return resp.StatusCode
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
