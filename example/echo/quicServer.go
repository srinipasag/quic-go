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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

const message = "hyuan:this is a log message from Oracle Weblogic Server 1030.001";

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {

	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		//return err
		fmt.Printf("111")
	}

	for i := 1;  i<=500000000000000; i++ {
		fmt.Println("---------------------------------------------")
		fmt.Printf("quicServer starting, iteration = %d \n", i);


		sess, err := listener.Accept(context.Background())
		if err != nil {
			//panic(err)
			fmt.Printf("112")
		}

		//go func() {

		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			//panic(err)
			fmt.Printf("113")
		}

		//fmt.Printf("Server: Got '%s'\n", string(stream))
		//Send2Splunk(string(stream));
		// Echo through the loggingWriter
		_, err = io.Copy(loggingWriter{stream}, stream)
		//time.Sleep(1 * time.Second)
  }
/*
	go func() { log.Fatal(echoServer()) }()

	err := clientMain()
	if err != nil {
		panic(err)
	}
*/
}

// Start a server that echos all data on the first stream opened by the client
func echoServer() error {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		//return err
		fmt.Printf("111")
	}

	sess, err := listener.Accept(context.Background())
	if err != nil {
		//panic(err)
		fmt.Printf("112")
	}

	//go func() {

		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			//panic(err)
			fmt.Printf("113")
		}

		//fmt.Printf("Server: Got '%s'\n", string(stream))
		//Send2Splunk(string(stream));
		// Echo through the loggingWriter
		_, err = io.Copy(loggingWriter{stream}, stream)
	time.Sleep(1 * time.Second)
		//panic(err)
	//}()

	return nil
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

	buf := make([]byte, len(message))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		return err
	}
	fmt.Printf("Client: Got '%s'\n", buf)

	return nil
}

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))

	// HEC sink goes here
	Send2Splunk(string(b));

	return 0, errors.New("no error")

	//return w.Writer.Write(b)

}

func Send2Splunk(msg  string) (int ) {

	//logger := util.NewLogger()

	//curl -k http://localhost:8088/services/collector -H "Authorization: Splunk f4d49a2b-9b70-488e-9789-86c18557515a"
	// -H "Content-Type: application/json" -d '{"event":"testevent1101","sourcetype":"test-auto-extract"}'
	url := "http://i-0447166d3d40f72ed.ec2.splunkit.io:8088/services/collector"

	//url = "https://api.playground.splunkbeta.com:443/hyuan/ingest/v1/events/"

	logMsg := fmt.Sprintf("Send2Splunk - sending event to %s ", url)

	fmt.Printf(logMsg)

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
	
	// Sending Request to Multiple instances

   	token1 := "Splunk bfb2fb1f-3e9e-457e-a5a2-01142f72cc29";
		 urls:=[8]string{"http://i-0447166d3d40f72ed.ec2.splunkit.io:8088/services/collector",
			 "http://i-0c30550d1f4cfb90a.ec2.splunkit.io:8088/services/collector",
			 "http://10.202.21.77:8088/services/collector",
			"http://10.202.17.2:8088/services/collector",
			"http://10.202.20.35:8088/services/collector",
			"http://10.202.20.96:8088/services/collector",
			"http://10.202.16.226:8088/services/collector",
			"http://10.202.22.20:8088/services/collector"}

		 //http://i-0c30550d1f4cfb90a.ec2.splunkit.io:8000
		 // http://10.202.21.77:8000
		 //http://10.202.17.2:8000
		 // http://10.202.20.35:8000
		 // http://10.202.20.96:8000
		 // http://10.202.16.226:8000
		 // http://10.202.22.20:8000


	for i := 0;  i<=7; i++ {
		//fmt.Printf("Welcome %d %s times\n",i,urls[i])
		logMsg := fmt.Sprintf("Send2Splunk - sending event to %d %s ", i, urls[i])

		fmt.Printf(logMsg)

		req2, err2 := http.NewRequest("POST", urls[i], bytes.NewBuffer(jsonBufferVal))
		req2.Header.Set("Authorization", token1)
		req2.Header.Set("Content-Type", "application/json")
		client2 := &http.Client{}
		resp2, err2 := client2.Do(req2)
		if err2 != nil {
			panic(err2)
		}
		defer resp2.Body.Close()
		logMsg = fmt.Sprintln("Response Status:", resp2.Status)
		fmt.Printf(logMsg)
		logMsg = fmt.Sprintln("Response Headers:", resp2.Header)
		fmt.Printf(logMsg)
		body, _ := ioutil.ReadAll(resp2.Body)
		logMsg = fmt.Sprintln("Response Body:", string(body))
		fmt.Printf(logMsg)
		//fmt.Println("Response Status:", resp1.Status)
		//fmt.Println("Response Headers:", resp1.Header)
	}

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
