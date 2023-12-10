package main

import (
	"context"
	"log"
	"strconv"
)

func Handler() (string, error) {

	ctx := context.Background()
	start := "2023-12-01"
	end := "2023-12-02"
	restaurant := "historySea"

	client := createClient(ctx)
	log.Println("FB client created")
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
	log.Println("Data from FB parced")

	res2, error := readData(ctx)
	//	fillTablesWithData(ctx, db, prefix)
	log.Println(*res2, error)

	text := "Данные синхронизированы " + strconv.Itoa(len(*res))
	message := Message{
		ChatID: -860192892,
		Text:   text,
	}
	SendMessage(&message)

	return text, nil
}
