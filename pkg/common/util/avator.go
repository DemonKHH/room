package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type ApiResult struct {
	Code   int    `json:"code"`
	Img    string `json:"img"`
	result int    `json:"result"`
}

// api
const API_URL string = "https://img.xjh.me/random_img.php"

func GetRandomAvator() string {
	var res ApiResult
	var err error
	var resp *http.Response
	params := url.Values{}
	Url, _ := url.Parse(API_URL)
	params.Set("type", "bg")
	params.Set("ctype", "acg")
	params.Set("return", "json")
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath) // https://httpbin.org/get?age=23&name=zhaofan
	resp, err = http.Get(urlPath)
	if err != nil {
		log.Println("request image api error: ", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("image api parse error: %v", err)
	}
	return "https:" + res.Img
}
