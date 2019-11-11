package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"bytes"
	"encoding/json"
	"strconv"

	//"net/url"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

//const message = "foobar"
const message = "Srinivas"

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {
	//go func() { log.Fatal(echoServer()) }()
	for {
		err :=	echoServer()
		if err != nil {
			fmt.Println("Error ... \n")
		}
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
	Send2ssc(string(b))
	return w.Writer.Write(b)
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

func Send2ssc(msg  string) (int ) {

	//logger := util.NewLogger()

	url := "http://i-0447166d3d40f72ed.ec2.splunkit.io:8088/services/collector"
	url1 := "http://i-0c30550d1f4cfb90a.ec2.splunkit.io:8088/services/collector"
	url2 := "http://10.202.22.20:8088/services/collector"
	//url2 := "http://10.202.22.20:8088/services/collector"


	fmt.Printf("Send2ssc - sending event to %s ", url)

	token := "Splunk 8d51f577-b1b9-414e-8e26-579a08ee6d90";
	token1 := "Splunk bfb2fb1f-3e9e-457e-a5a2-01142f72cc29";

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
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Headers:", resp.Header)
	result := fmt.Sprintln("%s%s", url+strconv.Itoa(1))
	fmt.Println(result)


	req1, err1 := http.NewRequest("POST", url1, bytes.NewBuffer(jsonBufferVal))
	req1.Header.Set("Authorization", token1)
	req1.Header.Set("Content-Type", "application/json")
	client1 := &http.Client{}
	resp1, err1 := client1.Do(req1)
	if err1 != nil {
		panic(err1)
	}
	defer resp1.Body.Close()
	fmt.Println("Response Status:", resp1.Status)
	fmt.Println("Response Headers:", resp1.Header)

	req2, err2 := http.NewRequest("POST", url2, bytes.NewBuffer(jsonBufferVal))
	req2.Header.Set("Authorization", token1)
	req2.Header.Set("Content-Type", "application/json")
	client2 := &http.Client{}
	resp2, err2 := client2.Do(req2)
	if err2 != nil {
		panic(err2)
	}
	defer resp2.Body.Close()
	fmt.Println("Response Status:", resp1.Status)
	fmt.Println("Response Headers:", resp1.Header)



  urls:=[8]string{"http://i-0447166d3d40f72ed.ec2.splunkit.io:8088/services/collector",
								 "http://i-0c30550d1f4cfb90a.ec2.splunkit.io:8088/services/collector",
								 "http://10.202.21.77:8088/services/collector",
									"http://10.202.17.2:8088/services/collector",
									"http://10.202.20.35:8088/services/collector",
									"http://10.202.20.96:8088/services/collector",
									"http://10.202.16.226:8088/services/collector",
									"http://10.202.22.20:8088/services/collector"}

	for i := 0;  i<=7; i++ {
		fmt.Printf("Welcome %d %s times\n",i,urls[i])

		req2, err2 := http.NewRequest("POST", urls[i], bytes.NewBuffer(jsonBufferVal))
		req2.Header.Set("Authorization", token1)
		req2.Header.Set("Content-Type", "application/json")
		client2 := &http.Client{}
		resp2, err2 := client2.Do(req2)
		if err2 != nil {
			panic(err2)
		}
		defer resp2.Body.Close()
		fmt.Println("Response Status:", resp1.Status)
		fmt.Println("Response Headers:", resp1.Header)

	}

	//logMsg = fmt.Sprintln("Response Status:", resp.Status)
	//logger.Debug(logMsg)
	//logMsg = fmt.Sprintln("Response Headers:", resp.Header)
	//logger.Debug(logMsg)
	//body, _ := ioutil.ReadAll(resp.Body)
	//logMsg = fmt.Sprintln("Response Body:", string(body))
	//logger.Debug(logMsg)
	return resp.StatusCode
}

