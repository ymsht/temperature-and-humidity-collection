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
	device_id int `db:"device_id primarykey, autoincrement"`
	serial_number string `db:"serial_number"`
	target_date time.Time `db:"target_date"`
	temperature float32 `db:"temperature"`
	humidity float32 `db:"humidity"`
	illumination float32 `db:"illumination"`
	movement float32 `db:"movement"`
}

func main()  {
	sdk := sdk.NatureRemoSdk{Token: "Bearer "}
	devices, err := sdk.GetDevice()
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf("%f\n", devices[0].Newest_events.Te.Val)

	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/db?parseTime=true")
	if err != nil {
		fmt.Printf(err.Error())
	}

	dialect := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}
	dbmap := &gorp.DbMap{Db: db, Dialect: dialect}
	dbmap.AddTableWithName(device{}, "deveices").SetKeys(false, "device_id")
	defer dbmap.Db.Close()

	tx, err := dbmap.Begin()
	if err != nil {
		fmt.Printf(err.Error())
	}
	
	device := device {
		serial_number: devices[0].Serial_number,
		target_date: time.Now(),
		temperature: devices[0].Newest_events.Te.Val,
		humidity: devices[0].Newest_events.Hu.Val,
		illumination: devices[0].Newest_events.Il.Val,
		movement: devices[0].Newest_events.Mo.Val,
	}

	err = tx.Insert(&device)
	if err != nil {
		fmt.Printf(err.Error())
	}

	tx.Commit()
}
