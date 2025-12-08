package main

import (
   "context"
   "encoding/json"
   "fmt"
   "log"

   "go.mongodb.org/mongo-driver/bson"
   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/mongo/options"
   "go.mongodb.org/mongo-driver/mongo/readpref"
)

type Movie struct {
   Title string `bson:"title"`
}

func main() {
   // Replace with your connection string
   uri := "<CONNECTION-STRING>"

   // Connect to MongoDB
   client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
   if err != nil {
      log.Fatalf("Failed to connect to MongoDB: %v", err)
   }
   defer func() {
      if err = client.Disconnect(context.TODO()); err != nil {
         log.Fatalf("Failed to disconnect MongoDB client: %v", err)
      }
   }()

   // Ping to confirm connection
   if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
      log.Fatalf("Failed to ping MongoDB: %v", err)
   }
   fmt.Println("Successfully connected to MongoDB")

   // Insert a document to a collection
   coll := client.Database("test").Collection("movies")
   newMovie := Movie{Title: "Back to the Future"}
   result, err := coll.InsertOne(context.TODO(), newMovie)
   if err != nil {
      log.Fatalf("Failed to insert document: %v", err)
   }
   fmt.Printf("Inserted document ID: %v\n", result.InsertedID)

   // List search indexes
   listOpts := options.ListSearchIndexesOptions{}
   ctx := context.TODO()
   cursor, err := coll.SearchIndexes().List(ctx, nil, &listOpts)
   if err != nil {
      log.Fatalf("Failed to list search indexes: %v", err)
   }

   var results []bson.D
   if err = cursor.All(ctx, &results); err != nil {
      log.Fatalf("Failed to iterate over cursor: %v", err)
   }

   res, err := json.Marshal(results)
   if err != nil {
      log.Fatalf("Failed to marshal results to JSON: %v", err)
   }
   fmt.Println("Search indexes found:", string(res))
}