package main

import (
	"fmt"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"

	sdk "github.com/ymsht/nature-remo-sdk"
)

type Device struct {
	Device_id int `db:"device_id"`
	SerialNumber string `db:"serial_number"`
	TargetDate time.Time `db:"target_date"`
	Temperature float32 `db:"temperature"`
	Humidity float32 `db:"humidity"`
	Illumination float32 `db:"illumination"`
	Movement float32 `db:"movement"`
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
	dbmap.AddTableWithName(Device{}, "deveices").SetKeys(false, "Device_id")
	defer dbmap.Db.Close()

	tx, err := dbmap.Begin()
	if err != nil {
		fmt.Printf(err.Error())
	}

	jst, _ := time.LoadLocation("Asia/Tokyo")
	now, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-01-01 00:00:00", jst)
	
	device := Device {
		SerialNumber: devices[0].Serial_number,
		TargetDate: now,
		Temperature: devices[0].Newest_events.Te.Val,
		Humidity: devices[0].Newest_events.Hu.Val,
		Illumination: devices[0].Newest_events.Il.Val,
		Movement: devices[0].Newest_events.Mo.Val,
	}

	err = tx.Insert(&device)
	if err != nil {
		fmt.Printf(err.Error())
	}

	tx.Commit()
}
