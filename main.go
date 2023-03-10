package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"log"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Item struct {
	ID    string  `json:"id"`
	Price float32 `json:"price"`
}

type CountReq struct {
	Num int `json:"num"`
}

func main() {
	http.HandleFunc("/api/v1/get_count", GetCount)
	http.HandleFunc("/api/v1/post_count", PostCount)
	http.HandleFunc("/api/v1/internal_call", InternalCall)
	http.HandleFunc("/v1/ping", Ping)
	http.HandleFunc("/api/v1/ping", Ping2)
	http.ListenAndServe(":8000", nil)
}

func GetCount(w http.ResponseWriter, req *http.Request) {
	fmt.Println("GetCount start")
	req.ParseForm()
	numStr := req.Form.Get("num")
	num, _ := strconv.Atoi(numStr)
	num++
	fmt.Printf("GetCount return num: %d\n", num)
	w.Write([]byte(strconv.Itoa(num)))
}

func PostCount(w http.ResponseWriter, req *http.Request) {
	fmt.Println("PostCount start")
	body, _ := ioutil.ReadAll(req.Body)
	countReq := new(CountReq)
	err := json.Unmarshal(body, countReq)
	if err != nil {
		fmt.Printf("json.Unmarshal failed, err: %+v", err)
		w.Write([]byte(fmt.Sprintf("json.Unmarshal failed, err: %+v", err)))
		return
	}
	num := countReq.Num
	num++
	fmt.Printf("PostCount return num: %d\n", num)
	w.Write([]byte(strconv.Itoa(num)))
}

func InternalCall(w http.ResponseWriter, req *http.Request) {
	toServiceID := req.Header.Get("X-SERVICE-ID")
	method := req.Header.Get("X-SERVICE-METHOD")
	v := req.Header.Get("X-SERVICE-VALUE")
	uri := req.Header.Get("X-SERVICE-URI")
	var resp *http.Response
	var err error
	if method == "GET" {
		resp, err = InternalCallGet(uri, toServiceID, map[string]string{"num": v}, map[string]string{"X-Test-Header1": "testHeader1"})
	} else if method == "POST" {
		num, _ := strconv.Atoi(v)
		body, _ := json.Marshal(CountReq{Num: num})
		resp, err = InternalCallPost(uri, toServiceID, bytes.NewBuffer(body), map[string]string{"X-Test-Header1": "testHeader1"})
	} else {
		w.Write([]byte(fmt.Sprintf("error: %+v", fmt.Errorf("invalid method: %s", method))))
	}

	if err != nil {
		w.Write([]byte(fmt.Sprintf("error: %+v", err)))
		fmt.Printf("error: %+v\n", err)
		return
	}
	defer resp.Body.Close()
	s, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error: %+v", err)))
	}
	fmt.Printf("resp:%+v, write: %s\n", resp, s)
	w.Write([]byte(s))
	fmt.Println("CallAnotherService")
}

func InternalCallGet(uri, toServiceID string, paramMap map[string]string, headers map[string]string) (*http.Response, error) {
	fromServiceID := os.Getenv("SERVICE_ID")
	urlValue := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s-%s.dycloud.service", fromServiceID, toServiceID),
		Path:   uri,
	}
	params := url.Values{}
	for k, v := range paramMap {
		params.Set(k, v)
	}
	urlValue.RawQuery = params.Encode()
	urlPath := urlValue.String()
	fmt.Println(urlPath)
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		fmt.Printf("http.NewRequest failed, err: %+v\n", err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	fmt.Printf("req: %+v\n", req)
	return http.DefaultClient.Do(req)
}

func InternalCallPost(uri, toServiceID string, body io.Reader, headers map[string]string) (*http.Response, error) {
	fromServiceID := os.Getenv("SERVICE_ID")
	urlPath := fmt.Sprintf("http://%s-%s.dycloud.service%s", fromServiceID, toServiceID, uri)
	req, err := http.NewRequest("POST", urlPath, body)
	if err != nil {
		fmt.Printf("http.NewRequest failed, err: %+v\n", err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	fmt.Printf("req: %+v\n", req)
	return http.DefaultClient.Do(req)
}

func Ping(w http.ResponseWriter, req *http.Request) {
	fmt.Println("hello /v1/ping")
}

func Ping2(w http.ResponseWriter, req *http.Request) {
	fmt.Println("hello /api/v1/ping")
}

func AsyncPrintLog() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Panic(err)
			}
		}()
		var (
			name1, name2, id1, id2 string  = "Tom", "Jack", "Item1", "Item2"
			age1, age2             int     = 18, 28
			price1, price2         float32 = 100, 200
		)
		p1 := Person{Name: name1, Age: age1}
		p2 := Person{Name: name2, Age: age2}
		i1 := Item{ID: id1, Price: price1}
		i2 := Item{ID: id2, Price: price2}
		p1Str, _ := json.Marshal(p1)
		p2Str, _ := json.Marshal(p2)
		i1Str, _ := json.Marshal(i1)
		i2Str, _ := json.Marshal(i2)
		for {
			fmt.Printf("%s\n", string(p1Str))
			fmt.Printf("%s\n", string(p2Str))
			fmt.Printf("%s\n", string(i1Str))
			fmt.Printf("%s\n", string(i2Str))
			time.Sleep(5 * time.Second)

		}
	}()
}
