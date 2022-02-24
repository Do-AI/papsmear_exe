package internal

import (
	"bytes"
	"encoding/json"
	"gocv.io/x/gocv"
	"io/ioutil"
	"log"
	"net/http"
)

// 외부 통신을 담당하는 코드

// postServer 함수는 target url에 post method를 이용해서 정보를 전송하는 역할을 담당한다.
func postServer(info interface{}, endpoint string) string {
	infoBytes, _ := json.Marshal(info)
	buff := bytes.NewBuffer(infoBytes)

	client := &http.Client{}

	url := CONFIG.Svc.URL + "/" + endpoint
	req, _ := http.NewRequest("POST", url, buff)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", GetAuthHeaderValue())

	resp, err := client.Do(req)
	if err != nil {
		message := makeSlackMessage(info)
		sendSlackMessage("post error\n" + message)
		log.Fatalln(err)
	}
	readBody, _ := ioutil.ReadAll(resp.Body)

	return string(readBody)
}

// UploadGcp 함수는 signed url을 spring papsmear 서버에서 받아와서 Google Storage에 이미지를 저장하는 역할을 한다.
func UploadGcp(url string, imgBuff *gocv.NativeByteBuffer) bool {

	buff := bytes.NewBuffer(imgBuff.GetBytes())

	client := &http.Client{}

	uploadReq, _ := http.NewRequest("PUT", url, buff)
	uploadReq.Header.Set("Content-Type", "image/jpeg")

	uploadResp, err := client.Do(uploadReq)

	if uploadResp != nil && uploadResp.Status != "200 OK" {
		return false
	}
	if err != nil {
		return false
	}
	defer func() {
		imgBuff.Close()
		buff.Reset()
	}()

	return true
}
