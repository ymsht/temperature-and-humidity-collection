package main

import (
	"fmt"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"

	sdk "github.com/ymsht/nature-remo-sdk"
)

type device struct {
	serial_number string
	target_date time.Time
	temperature float32
	humidity float32
	illumination float32
	movement float32
}

func main()  {
	sdk := sdk.NatureRemoSdk{Token: "Bearer "}
	devices, err := sdk.GetDevice()
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf("%f\n", devices[0].Newest_events.Te.Val)

	db, err := sql.Open("mysql", ":@tcp(localhost:3306)/db?parseTime=true")
	if err != nil {
		fmt.Printf(err.Error())
	}

	dialect := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}
	dbmap := &gorp.DbMap{Db: db, Dialect: dialect}
	dbmap.AddTableWithName(device{}, "deveices")
	defer dbmap.Db.Close()

	tx, err := dbmap.Begin()
	if err != nil {
		fmt.Printf(err.Error())
	}
	
	device := device {
		devices[0].Serial_number,
		time.Now(),
		devices[0].Newest_events.Te.Val,
		devices[0].Newest_events.Hu.Val,
		devices[0].Newest_events.Il.Val,
		devices[0].Newest_events.Mo.Val,
	}

	err = tx.Insert(&device)
	if err != nil {
		fmt.Printf(err.Error())
	}

	tx.Commit()
}
