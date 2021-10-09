package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// flagPort is the open port the application listens on
	flagPort = flag.String("port", "9000", "Port to listen on")
)

var results []map[string]interface{}

//GetHandler handles the index route
func GetHandler(w http.ResponseWriter, r *http.Request) {
	// err := json.NewDecoder(r.Body).Decode(&results)
	// if err != nil {
	// 	panic(err)
	// }
	jsonBody, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error converting results to json",
			http.StatusInternalServerError)
	}
	fmt.Fprint(w, "Get done")
	w.Write(jsonBody)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		//results = append(results, string(body))
		w.Write(body)
		fmt.Fprint(w, "POST done")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
}

func connectDB() {
	jsonBody1, err := json.Marshal(results)
	if err != nil {

	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	//fmt.Println("ClientOptopm TYPE:", reflect.TypeOf(clientOptions), "\n")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Mongo.connect() ERROR: ", err)
		os.Exit(1)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	col := client.Database("Users").Collection("Details")
	fmt.Println("Collection Type: ", reflect.TypeOf(col), "\n")
	result, insertErr := col.InsertOne(ctx, jsonBody1)
	if insertErr != nil {
		fmt.Println("InsertONE Error:", insertErr)
		os.Exit(1)
	} else {
		fmt.Println("InsertOne() result type: ", reflect.TypeOf(result))
		fmt.Println("InsertOne() api result type: ", result)

		newID := result.InsertedID
		fmt.Println("InsertedOne(), newID", newID)
		fmt.Println("InsertedOne(), newID type:", reflect.TypeOf(newID))
	}
}

func returnSearchResult() {

}

func main() {

	results := []map[string]interface{}{
		{"Id": 2, "Name": "keyTestTest", "Email": "com.app", "Mobile": 1234},
	}

	// loop over elements of slice
	for _, m := range results {

		// loop over keys in map
		for k, v := range m {
			fmt.Println(k, "value is", v)
		}
	}

	connectDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/", GetHandler)
	mux.HandleFunc("/users", PostHandler)
	//http.HandleFunc("/user/search", returnSearchResult)
	//http.HandleFunc("/posts", returnAllArticles)

	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe(":"+*flagPort, mux))

}
