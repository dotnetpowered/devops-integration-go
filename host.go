package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Host struct {
	bun.BaseModel `bun:"table:host"`

	Name        string
	Description string
	Template    string
	Ip          string
	OS          string
	Tags        []string `bun:",array"`
	Status      string
	NumCpu      int
	MemSize     int
	Notes       string
	Uptime      int
	Source      string
}

func printHosts(hosts []Host) {
	for i := 0; i < len(hosts); i++ {
		h := hosts[i]
		h1, _ := json.Marshal(h)
		fmt.Println(string(h1))
		//fmt.Println(h.name + " (" + h.ip + ") " + h.os)
	}
}

func saveHosts(hosts []Host, dsn string) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	ctx := context.Background()
	for i := 0; i < len(hosts); i++ {
		h := hosts[i]
		_, err := db.NewInsert().Model(&h).Exec(ctx)

		if err != nil {
			panic(err)
		}
	}
}
