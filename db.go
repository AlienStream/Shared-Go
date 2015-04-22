package Shared

import (
	"github.com/ziutek/mymysql/mysql"
    	_ "github.com/ziutek/mymysql/thrsafe"
	//sq "github.com/lann/squirrel"
)

const database = "alien";
const host = "127.0.0.1:33060";
const user = "homestead";
const password = "secret";

func getDBConnection() mysql.Conn {
	db := mysql.New("tcp", "", host, user, password, database)
	err := db.Connect()
    	if err != nil {
        	panic(err)
    	}
	return db;
}
