package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang191119/Buoi2/hw4/lodash"

	"github.com/patrickmn/go-cache"
)

var c = cache.New(time.Duration(5)*time.Minute, time.Duration(10)*time.Minute)

type Request struct {
	List    []int `json:"list"`
	Element int   `json:"element"`
}

func ParseData(data interface{}) (string, error) {
	byteData, err := json.Marshal(data)
	if err == nil {
		return string(byteData), nil
	}
	return "", err
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is a website server by a Go HTTP server.")
	})

	http.HandleFunc("/contain", func(w http.ResponseWriter, r *http.Request) {
		//parse params
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("\nCan not parse the api params")
			panic(err)
		}
		var req Request
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Println("\nCan not decode")
			panic(err)
		}

		// Find element in list
		begin := time.Now()
		result := lodash.Contain(c, &req.List, req.Element, cache.DefaultExpiration)
		end := time.Now()

		fmt.Printf("\nSpend: %d miliseconds\n", end.Sub(begin).Nanoseconds()/1000000)
		res, err := ParseData(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
		}
		fmt.Fprintf(w, res)
	})

	http.HandleFunc("/reverse", func(w http.ResponseWriter, r *http.Request) {
		//parse params
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("\nCan not parse the api params")
			panic(err)
		}
		var req Request
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Println("\nCan not decode")
			panic(err)
		}

		//reverse slice
		result, err := ParseData(lodash.Reverse(req.List))
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
		} else {
			log.Println(result)
			fmt.Fprintf(w, result)
		}
	})

	http.ListenAndServe(":5001", nil)
}
