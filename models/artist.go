package models

import (
	"errors"
	"fmt"
	db "github.com/AlienStream/Shared-Go/database"
	mysql "github.com/ziutek/mymysql/mysql"
	"time"
)

type Artist struct {
	Id             int
	Name           string
	Thumbnail      string
	Favorite_count string
	Play_count     int
	Updated_at     time.Time
	Created_at     time.Time
}

func AllArtists() []Post {
	rows, _, err := db.Con.Query("select * from artists")
	if err != nil {
		panic(err)
	}

	return RowsToArtists(rows)
}

func (a Artist) FromId(Id int) (Artist, error) {
	rows, _, err := db.Con.Query("select * from artists where `id`=%s", Id)
	if err != nil {
		return nil, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return nil, errors.New("Artist Not Found For ID %d", Id)
	}

	return RowsToArtists(rows)[0], nil
}

func (a Artist) IsNew() (bool, int) {
	rows, _, err := db.Con.Query("select * from artists where `name` = '%s'", a.Name)
	if err != nil {
		panic("Error When Querying the database")
	}

	if len(rows) == 0 {
		return true, 0
	}

	return false, RowsToArtists(rows)[0].Id

}

func (a Artist) Insert() error {
	fmt.Printf("Inserting New Artist %s \n", a.Name)
	stmt, err := db.Con.Prepare("insert into artists (`name`, `thumbnail`, `favorite_count`, `play_count`, `created_at`, `updated_at`) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(a.Name, a.Thumbnail, a.Favorite_count, a.Play_count, time.Now(), time.Now())

	return nil
}

func (a Artist) Save() error {
	fmt.Printf("Updating Artist %s \n", a.Name)
	if a.Id < 1 {
		return errors.New("Invalid ID for Artist")
	}
	stmt, err := db.Con.Prepare("update artists set `name`=?, `thumbnail`=?, `favorite_count`=?, `play_count`=?,`updated_at`=? where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(a.Name, a.Thumbnail, a.Favorite_count, a.Play_count, time.Now(), a.Id)

	return nil

}

func (a Artist) Delete() error {
	fmt.Printf("Deleting Artist %s \n", a.Name)
	if p.Id < 1 {
		return errors.New("Invalid ID for Artist")
	}
	stmt, err := db.Con.Prepare("delete from artists where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(a.Id)

	return nil
}

func RowsToArtists(rows []mysql.Row) []Post {
	var artists = []Artist{}
	for _, row := range rows {
		artists = append(artists, Artist{}.FromRow(row))
	}
	return artists
}

func (a Artist) FromRow(row mysql.Row) Artist {
	a.Id = row.Int(0)
	a.Name = row.Str(1)
	a.Thumbnail = row.Str(2)
	a.Favorite_count = row.Int(3)
	a.Play_count = row.Int(4)
	return a
}
