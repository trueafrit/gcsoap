// msoap.go
package gcsoap

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type XMLQuery struct {
	Loc string `xml:",chardata"`
}

type MsoapBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	Body    interface{}
}

type MsoapRequest struct {
	XMLName      xml.Name `xml:"soapenv:Envelope"`
	XMLNsSoapEnv string   `xml:"xmlns:soapenv,attr"`
	XMLNsCBD2    string   `xml:"xmlns:cbd2,attr"`

	Header MsoapRequestHeader
	Body   MsoapBody
}

type MsoapRequestHeader struct {
	XMLName  xml.Name `xml:"soapenv:Header"`
	AuthData MsoapAuthData
}

type MsoapAuthData struct {
	XMLName  xml.Name `xml:"cbd2:AuthData"`
	Login    string   `xml:"cbd2:login"`
	Password string   `xml:"cbd2:password"`
}

func MsoapCall(url string, action string, reqb []byte, timeoutSec int) ([]byte, error) {
	if timeoutSec == 0 {
		timeoutSec = 30
	}

	timeout := time.Duration(timeoutSec) * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqb))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", action)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	// dump, err := httputil.DumpRequestOut(req, true)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("%q", dump)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Println(string(bodyBytes))
	defer response.Body.Close()
	return bodyBytes, nil
}

func MsoapCallToFile(fname string, url string, action string, reqb []byte, timeoutSec int) error {
	if timeoutSec == 0 {
		timeoutSec = 30
	}

	timeout := time.Duration(timeoutSec) * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqb))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", action)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	//fmt.Println(string(bodyBytes))
	defer response.Body.Close()

	out, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755) //os.Create(fname)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	return nil
}
