package database

import (
	"github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/thrsafe"
	"os"
	//sq "github.com/lann/squirrel"
)

var Con *autorc.Conn

func init() {
	var database string = os.Getenv("DB_DATABASE") 
	if database == "" {
		database = "homestead"
	}

	var host string = os.Getenv("DB_HOST") 
	if host == "" {
		host = "127.0.0.1:33060"
	}

	var user string = os.Getenv("DB_USERNAME")
	if user == "" {
		 user = "homestead"
	}
	
	var password string = os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "secret"
	}

	Con = autorc.New("tcp", "", host, user, password, database)
}
