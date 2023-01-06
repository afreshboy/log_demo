package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func main() {
	AsyncPrintLog()
	http.HandleFunc("/v1/ping", Ping1)
	http.HandleFunc("/api/v1/ping", Ping2)
	http.ListenAndServe(":8000", nil)

}

func Ping1(w http.ResponseWriter, req *http.Request) {
	fmt.Println("helloping1")
}

func Ping2(w http.ResponseWriter, req *http.Request) {
	fmt.Println("helloping2")
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
