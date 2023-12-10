package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Order struct {
	Id                  string
	Name                string
	DateTime            time.Time
	PackingMode         string
	Positions           string
	Total               float32
	Address             string
	Phone               string
	Duration            int
	Discount            string
	DiscountType        string
	Restaurant          string
	PaymentMethod       string
	ExternalOrderNumber string
	IsReady             bool
}

func createClient(ctx context.Context) *firestore.Client {

	projectID := "rovesnikiprod"

	json, exists := os.LookupEnv("GOOGLE_SERVICE_ACCOUNT")
	if !exists {
		panic("GOOGLE_SERVICE_ACCOUNT environment variable not defined")
	}

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(json)))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func getData(start string, end string, restaurant string, client *firestore.Client, ctx context.Context) *firestore.DocumentIterator {

	const layout = "2006-01-02 03:04:05 PM"
	starttm, _ := time.Parse(layout, start+" 00:00:00.935 AM")
	endtm, _ := time.Parse(layout, end+" 11:59:59.935 PM")

	iter := client.Collection(restaurant).
		Where("dateTime", ">", starttm).
		Where("dateTime", "<", endtm).
		Limit(3).
		Documents(ctx)

	return iter
}

func parceData(iter *firestore.DocumentIterator) *[]Order {

	res := []Order{}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		var data Order

		config := &mapstructure.DecoderConfig{
			ErrorUnused: false,
			Result:      &data,
		}
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			log.Fatal(err)
		}
		if err := decoder.Decode(doc.Data()); err != nil {
			log.Fatal(err)
		}

		res = append(res, data)
	}

	return &res
}
