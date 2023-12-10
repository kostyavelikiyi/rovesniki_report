package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/retry"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	yc "github.com/ydb-platform/ydb-go-yc"
)

type table364 struct {
	id   uint64
	name string
}

func writeData() error {

	return fmt.Errorf("Error")
}

func readData(ctx context.Context) (*table364, error) {

	res := table364{}

	db, err := ydb.Open(ctx,
		os.Getenv("YDB_CONNECTION_STRING"),
		yc.WithMetadataCredentials(),
		yc.WithInternalCA(), // append Yandex Cloud certificates
	)
	if err != nil {
		panic(err)
	}
	defer db.Close(ctx)

	log.Println("DB connetcted")

	return &res, nil
}

func render(t *template.Template, data interface{}) string {
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func sliceToInterfaces[T any](v []T) []interface{} {
	ii := make([]interface{}, len(v))
	for i, vv := range v {
		ii[i] = vv
	}
	return ii
}

func fillTablesWithData(ctx context.Context, db *sql.DB, prefix string) (err error) {
	series, seasonsData, episodesData := getData2()
	args := []sql.NamedArg{
		sql.Named("seriesData", types.ListValue(series...)),
		sql.Named("seasonsData", types.ListValue(seasonsData...)),
		sql.Named("episodesData", types.ListValue(episodesData...)),
	}
	declares, err := sugar.GenerateDeclareSection(args)
	if err != nil {
		return err
	}
	err = retry.DoTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		if _, err = tx.ExecContext(ctx,
			render(
				template.Must(template.New("").Parse(`
					PRAGMA TablePathPrefix("{{ .TablePathPrefix }}");

					{{ .Declares }}
				
					REPLACE INTO series
					SELECT
						series_id,
						title,
						series_info,
						release_date,
						comment
					FROM AS_TABLE($seriesData);
						
					REPLACE INTO seasons
					SELECT
						series_id,
						season_id,
						title,
						first_aired,
						last_aired
					FROM AS_TABLE($seasonsData);
						
					REPLACE INTO episodes
					SELECT
						series_id,
						season_id,
						episode_id,
						title,
						air_date
					FROM AS_TABLE($episodesData);
				`)), struct {
					TablePathPrefix string
					Declares        string
				}{
					TablePathPrefix: prefix,
					Declares:        declares,
				},
			),
			sliceToInterfaces(args)...,
		); err != nil {
			return err
		}
		return nil
	}, retry.WithDoTxRetryOptions(retry.WithIdempotent(true)))
	if err != nil {
		return fmt.Errorf("upsert query failed: %w", err)
	}
	return nil
}
