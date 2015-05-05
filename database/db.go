package database

import (
	"github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/thrsafe"
	//sq "github.com/lann/squirrel"
)

const database = "homestead"
const host = "127.0.0.1:33060"
const user = "homestead"
const password = "secret"

var Con *autorc.Conn

func init() {
	Con = autorc.New("tcp", "", host, user, password, database)
}
