package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var host = "http://176.9.59.110:8000"
var address = "<SiaAddress>"
var getworkpath = "/miner/getwork/"
var submitblockpath = "/miner/submitwork/"

type job struct {
	WorkerAddress string
	JobID         int64
	BlockTarget   []byte
	ShareTarget   []byte
	Header        []byte
	Nonce         uint64
	Legacy        []byte
}

func getHeaderForWork() (share job, target, header []byte, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", host+getworkpath+address, nil)
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", "Sia-Agent")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
	case 400:
		err = fmt.Errorf("Invalid siad response, status code %d, is your wallet initialized and unlocked?", resp.StatusCode)
		return
	default:
		err = fmt.Errorf("Invalid siad, status code %d", resp.StatusCode)
		return
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if len(buf) < 112 {
		err = fmt.Errorf("Invalid siad response, only received %d bytes, is your wallet initialized and unlocked?", len(buf))
		return
	}

	//	target = buf[:32]
	//	header = buf[32:112]

	err = json.Unmarshal(buf, &share)

	if err != nil {
		log.Fatal(err)
	}

	target = share.ShareTarget
	header = share.Header

	return
}

func submitHeader(header []byte, share job) (err error) {
	share.Legacy = header

	enc, err := json.Marshal(share)
	if err != nil {
		log.Fatal("I can't json :(")
	}

	req, err := http.NewRequest("POST", host+submitblockpath+address, bytes.NewReader(enc))
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", "Sia-Agent")

	client := &http.Client{}
	_, err = client.Do(req)

	return
}
