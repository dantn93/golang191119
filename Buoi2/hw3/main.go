package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type Student struct {
	Id        string `json:"_id" bson:"_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Age       int    `json:"age" bson:"age"`
	Class     string `json:"class_name" bson:"class_name"`
}

func GetClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:wBui63uwex1Imla9@cluster0-nmhgi.mongodb.net/test?retryWrites=true&w=majority")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Println(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Println(err)
	}
	return client
}

func Insert(client *mongo.Client, data []interface{}) error {
	collection := client.Database("db").Collection("students")
	_, err := collection.InsertMany(context.TODO(), data)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	return err
}

func RemoveAll(client *mongo.Client) (err error) {
	collection := client.Database("db").Collection("students")
	err = collection.Drop(context.TODO())
	return err
}

func GetAll(client *mongo.Client) ([]Student, error) {
	collection := client.Database("db").Collection("students")
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}

	students := []Student{}
	for cur.Next(context.TODO()) {
		var s Student
		err = cur.Decode(&s)
		if err != nil {
			log.Println("Error on Decoding the document", err)
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func main() {

	// ========== GET STUDENTS ========== //
	res, err := http.Get("http://localhost:3001/getStudent")
	if err != nil {
		log.Println("ERROR1: ", err.Error())
		return
	}

	defer res.Body.Close()

	// var students interface{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ERROR2: ", err.Error())
		return
	}

	var students []Student
	err = json.Unmarshal(body, &students)
	if err != nil {
		log.Println("ERROR3: ", err.Error())
		return
	}

	// ========== WRITE INTO MONGO ========== //

	c := GetClient()
	err = c.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Println("Couldn't connect to the database", err)
		return
	} else {
		log.Println("Connected!")
	}

	// Remove all records
	err = RemoveAll(c)
	if err != nil {
		log.Println("Can not remove all records: ", err.Error())
		return
	}

	// Add _id to records
	var s []interface{}
	for _, v := range students {
		item := map[string]interface{}{
			"_id":        bson.NewObjectId().Hex(),
			"first_name": v.FirstName,
			"last_name":  v.LastName,
			"age":        v.Age,
			"class_name": v.Class,
		}
		s = append(s, item)
	}

	// Insert into db
	err = Insert(c, s)
	if err != nil {
		log.Println("ERROR5: ", err.Error())
	}

	// Get all record after add into db
	allStudent, err := GetAll(c)
	if err != nil {
		log.Println("Can not get all records: ", err.Error())
		return
	}
	fmt.Printf("\n\nAll student in db: %+v\n", allStudent)
}
