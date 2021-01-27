package main

import (
	"time"
	"os"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"

	sdk "github.com/ymsht/nature-remo-sdk"

	log "github.com/sirupsen/logrus"
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

type GorpTracer struct{}

func (t *GorpTracer) Printf(format string, v ...interface{}) {
	log.Info(v...)
}

func main()  {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	sdk := sdk.NatureRemoSdk{Token: "Bearer "}
	devices, err := sdk.GetDevice()
	if err != nil {
		log.Error(err.Error())
	}

	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/db?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		log.Error(err.Error())
	}

	dialect := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}
	dbmap := &gorp.DbMap{Db: db, Dialect: dialect}
	dbmap.AddTableWithName(Device{}, "deveices")
	defer dbmap.Db.Close()

	tracer := &GorpTracer{}
	dbmap.TraceOn("[SQL]", tracer)

	tx, err := dbmap.Begin()
	if err != nil {
		log.Error(err.Error())
	}

	utc := time.Now().UTC()
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
    	log.Error(err.Error())
	}
	
	device := Device {
		SerialNumber: devices[0].SerialNumber,
		TargetDate: utc.In(jst),
		Temperature: devices[0].NewestEvents.Te.Val,
		Humidity: devices[0].NewestEvents.Hu.Val,
		Illumination: devices[0].NewestEvents.Il.Val,
		Movement: devices[0].NewestEvents.Mo.Val,
	}

	err = tx.Insert(&device)
	if err != nil {
		log.Error(err.Error())
	}

	tx.Commit()
}
