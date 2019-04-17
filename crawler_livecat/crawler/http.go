package crawler

import (
	"bytes"
	"crawler_livecat/mongogo"
	"encoding/json"
	"net/http"
)

func sendToElasticByHTTPPost(data mongogo.MongoDB) {
	url := "http://120.126.16.88:17777/add"

	dataByte, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataByte))
	handleError(err, "[http.go] NewRequest err")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	handleError(err, "[http.go] send data err")
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
}

func updateToElasticByHTTPPost(data mongogo.MongoDB) {

}
