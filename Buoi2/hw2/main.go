package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Student struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Age       int    `json:"age" bson:"age"`
	Class     string `json:"class_name" bson:"class_name"`
}

func ParseData(data interface{}) (string, error) {
	byteData, err := json.Marshal(data)
	if err == nil {
		return string(byteData), nil
	}
	return "", err
}

var students = []Student{
	Student{
		FirstName: "Nguyen",
		LastName:  "Hoang",
		Age:       15,
		Class:     "10A1",
	},
	Student{
		FirstName: "Le",
		LastName:  "Van",
		Age:       15,
		Class:     "10A1",
	},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is a website server by a Go HTTP server.")
	})

	http.HandleFunc("/getStudent", func(w http.ResponseWriter, r *http.Request) {
		resp, err := ParseData(students)
		if err != nil {
			//Return with status Internal Server Error
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
		}
		//Return with status OK
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, resp)

	})

	http.ListenAndServe(":3001", nil)
}
