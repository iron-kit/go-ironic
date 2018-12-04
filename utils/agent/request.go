package agent

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(url string) (*http.Response, []byte) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("http.NewRequest:", err.Error())
		return nil, nil
	}

	fmt.Println(r.Proto)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Println("http.DefaultClient.Do:", err.Error())
		return nil, nil
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp.StatusCode not ok", resp.StatusCode)
		return nil, nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		fmt.Println("ioutil.ReadAll:", err.Error())
		return nil, nil
	}

	fmt.Println(string(data))
	// return data
	return resp, data
}

func Post(url string, headers map[string]string, body interface{}) (*http.Response, []byte) {

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		fmt.Println("http.NewRequest.ErrorBody:", err.Error())
		return nil, nil
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(bodyJSON)))

	if err != nil {
		fmt.Println("http.NewRequest:", err.Error())
		return nil, nil
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("http.request.error:", err.Error())
		return nil, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp.StatusCode not ok", resp.StatusCode)
		return nil, nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		fmt.Println("ioutil.ReadAll:", err.Error())
		return nil, nil
	}

	return resp, data
}

func Put(url string, headers map[string]string, body interface{}) (*http.Response, []byte) {

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		fmt.Println("http.NewRequest.ErrorBody:", err.Error())
		return nil, nil
	}

	req, err := http.NewRequest("PUT", url, strings.NewReader(string(bodyJSON)))

	if err != nil {
		fmt.Println("http.NewRequest:", err.Error())
		return nil, nil
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("http.request.error:", err.Error())
		return nil, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp.StatusCode not ok", resp.StatusCode)
		return nil, nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		fmt.Println("ioutil.ReadAll:", err.Error())
		return resp, nil
	}

	return resp, data
}

func Delete(url string, headers map[string]string, body interface{}) (*http.Response, []byte) {

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		fmt.Println("http.NewRequest.ErrorBody:", err.Error())
		return nil, nil
	}

	req, err := http.NewRequest("DELETE", url, strings.NewReader(string(bodyJSON)))

	if err != nil {
		fmt.Println("http.NewRequest:", err.Error())
		return nil, nil
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("http.request.error:", err.Error())
		return nil, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp.StatusCode not ok", resp.StatusCode)
		return nil, nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		fmt.Println("ioutil.ReadAll:", err.Error())
		return resp, nil
	}

	return resp, data
}
