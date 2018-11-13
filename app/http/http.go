package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"validator-monitor/utils"

	"github.com/astaxie/beego"
)

var (
	ClientRESTServerAddress string
	DefaultHTTPClient       *http.Client = nil
)

func init() {
	ClientRESTServerAddress = beego.AppConfig.String("client-rest-server::address")
	if "" == ClientRESTServerAddress {
		panic("client-rest-server::address invalid")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	DefaultHTTPClient = &http.Client{Transport: tr}
}

func Call(method, url string, reqData, result interface{}) error {
	var body io.Reader = nil
	beego.Debug(reqData)
	if nil != reqData {
		bin, err := utils.CDC.MarshalJSON(reqData)
		if nil != err {
			beego.Error(err)
			return err
		}

		body = bytes.NewReader(bin)
	}

	req, err := http.NewRequest(method, url, body)
	if nil != err {
		beego.Error(err)
		return err
	}

	resp, err := DefaultHTTPClient.Do(req)
	if nil != err {
		beego.Error(err)
		return err
	}
	defer resp.Body.Close()

	bin, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		beego.Error(err)
		return err
	}
	beego.Debug(method, url, "response:", string(bin))

	if http.StatusOK != resp.StatusCode {
		err = fmt.Errorf("%s %s,expected 200,got:%d,err:%s", method, url, resp.StatusCode, string(bin))
		beego.Error(err)

		return err
	}

	if 0 != len(bin) && nil != result {
		return json.Unmarshal(bin, result)
	}

	return nil
}

func POST(url string, reqData, result interface{}) error {
	return Call("POST", url, reqData, result)
}

func GET(url string, result interface{}) error {
	return Call("GET", url, nil, result)
}
