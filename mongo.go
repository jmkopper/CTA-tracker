// mongo.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoSearchHandler struct {
	config *configSettings
}

type MongoResponse struct {
	Contents []BusStop `json:"mongo-response"`
}

type BusStop struct {
	ID                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	LocationType       string             `bson:"location_type,omitempty"`
	StopCode           string             `bson:"stop_code,omitempty"`
	StopDesc           string             `bson:"stop_desc,omitempty"`
	StopID             string             `json:"stpid" bson:"stop_id,omitempty"`
	StopName           string             `json:"stpnm" bson:"stop_name,omitempty"`
	StopLat            string             `bson:"stop_lat,omitempty"`
	StopLon            string             `bson:"stop_lon,omitempty"`
	WheelchairBoarding string             `bson:"wheelchair_boarding,omitempty"`
}

type mongoSearchRequest struct {
	QueryString string `json:"queryString"`
}

// called by /search endpoint
// pass the query data to ctaStopSearch() which handles the mongo search, then transmit the results back to the frontend
func (mh mongoSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("content-type", "application/json")

	var postBody mongoSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&postBody); err != nil {
		log.Println(err)
	}

	searchResponse := make(chan MongoResponse)
	go ctaStopSearch(searchResponse, postBody.QueryString, mh.config)
	// error handler
	select {
	case resp := <-searchResponse:
		json.NewEncoder(w).Encode(resp)
	case <-time.After(time.Second * 5):
		fmt.Fprintf(w, "timeout")
	}
}

// call the mongo search and pass the response back to the HTTP handler
func ctaStopSearch(mongoResponse chan<- MongoResponse, queryString string, config *configSettings) {
	connectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", config.mongoUsername, config.mongoPassword, config.mongoClusterAddress)
	/*
	   Connect to the cluster
	*/
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	stopsCollection := client.Database(config.db).Collection(config.collection)

	// formulate the search filter
	filter := bson.M{
		"$text": bson.M{
			"$search": queryString,
		},
	}
	// sort by "textScore" (it's a Mongo thing)
	opts := options.Find().SetSort(bson.M{
		"score": bson.M{
			"$meta": "textScore",
		},
	})

	cursor, err := stopsCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		panic(err)
	}

	var results []BusStop
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	mongoResponse <- MongoResponse{Contents: results}
}
