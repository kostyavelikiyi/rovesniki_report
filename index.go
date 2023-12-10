package main

import (
	"context"
	"log"
	"os"
	"strconv"
)

func Handler() (string, error) {

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	ctx := context.Background()
	start := "2023-12-01"
	end := "2023-12-02"
	restaurant := "historySea"

	client := createClient(ctx)
	logger.Println("FB client created")
	defer client.Close()

	// db, err := sql.Open("ydb", dsn)
	// if err != nil {
	// 	log.Fatalf("connect error: %v", err)
	// }
	// defer func() { _ = db.Close() }()
	// db.SetMaxOpenConns(50)
	// db.SetMaxIdleConns(50)
	// db.SetConnMaxIdleTime(time.Second)

	res := parceData(getData(start, end, restaurant, client, ctx))
	logger.Println("Data from FB parced")

	res2, error := readData(ctx)
	//	fillTablesWithData(ctx, db, prefix)
	logger.Println(*res2, error)

	text := "Данные синхронизированы " + strconv.Itoa(len(*res))
	message := Message{
		ChatID: -860192892,
		Text:   text,
	}
	SendMessage(&message)

	return text, nil
}
