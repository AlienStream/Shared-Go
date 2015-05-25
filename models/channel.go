package models

import (
	"errors"
	"fmt"
	db "github.com/AlienStream/Shared-Go/database"
	mysql "github.com/ziutek/mymysql/mysql"
	"time"
)

type Channel struct {
	Id         int
	Url        string
	Artist_id  int
	Updated_at time.Time
	Created_at time.Time
}

func AllChannels() []Post {
	rows, _, err := db.Con.Query("select * from channels")
	if err != nil {
		panic(err)
	}

	return RowsToChannels(rows)
}

func (c Channel) FromId(Id int) (Channel, error) {
	rows, _, err := db.Con.Query("select * from channels where `id`=%s", Id)
	if err != nil {
		return nil, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return nil, errors.New("Channel Not Found For ID %d", Id)
	}

	return RowsToChannels(rows)[0], nil
}

func (c Channel) IsNew() (bool, int) {
	rows, _, err := db.Con.Query("select * from channels where `url` = '%s'", c.Url)
	if err != nil {
		panic("Error When Querying the database")
	}

	if len(rows) == 0 {
		return true, 0
	}

	return false, RowsToChannels(rows)[0].Id

}

func (c Channel) Insert() error {
	fmt.Printf("Inserting New Channel %s \n", c.Url)
	stmt, err := db.Con.Prepare("insert into channels (`url`, `artist_id`, `created_at`, `updated_at`) values (?, ?, ?, ?)")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(c.Url, c.Artist_id, time.Now(), time.Now())

	return nil
}

func (c Channel) Save() error {
	fmt.Printf("Updating Channel %s \n", c.Url)
	if c.Id < 1 {
		return errors.New("Invalid ID for Channel")
	}
	stmt, err := db.Con.Prepare("update channel set `url`=?, `artist_id`=?, `updated_at`=? where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(c.Url, c.Artist_id, time.Now(), c.Id)

	return nil

}

func (c Channel) Delete() error {
	fmt.Printf("Deleting Post %s \n", c.Url)
	if c.Id < 1 {
		return errors.New("Invalid ID for Channel")
	}
	stmt, err := db.Con.Prepare("delete from channels where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(c.Id)

	return nil
}

func RowsToChannels(rows []mysql.Row) []Channel {
	var channels = []Channel{}
	for _, row := range rows {
		channels = append(channels, Channel{}.FromRow(row))
	}
	return channels
}

func (c Channel) FromRow(row mysql.Row) Channel {
	c.Id = row.Int(0)
	c.Url = row.Str(1)
	c.Artist_id = row.Int(2)
	c.Updated_at = row.Localtime(3)
	c.Created_at = row.Localtime(4)
	return c
}
