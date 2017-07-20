package internal

import (
	"my-game/conf"
	"github.com/name5566/leaf/db/mongodb"
	"github.com/name5566/leaf/log"
)

var mongoDB *mongodb.DialContext



func init() {
	// mongodb
	if conf.Server.DBMaxConnNum <= 0 {
		conf.Server.DBMaxConnNum = 100
	}
	db, err := mongodb.Dial(conf.Server.DBUrl, conf.Server.DBMaxConnNum)
	if err != nil {
		log.Fatal("dial mongodb error: %v", err)
	}
	mongoDB = db

	/*
		// users
		err = db.EnsureUniqueIndex((DBName), "users", []string{"accid"})
		if err != nil {
			log.Fatal("ensure index error: %v", err)
		}
	*/
	err = db.EnsureCounter(DBName, "counters", "myusers")
	if err != nil {
		log.Fatal("ensure counter error: %v", err)
	}

	err = db.EnsureCounter(DBName, "counters", "myrooms")
	if err != nil {
		log.Fatal("ensure counter error: %v", err)
	}
}

func mongoDBDestroy() {
	mongoDB.Close()
	mongoDB = nil
}

func mongoDBNextSeq(id string) (int, error) {
	return mongoDB.NextSeq(DBName, "counters", id)
}