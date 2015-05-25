package models

import (
	"errors"
	"fmt"
	db "github.com/AlienStream/Shared-Go/database"
	mysql "github.com/ziutek/mymysql/mysql"
	"time"
)

type Track struct {
	Id             int
	Rank           int
	Title          string
	Thumbnail      string
	Favorite_count int
	Play_count     int
	Artist_id      int
	Updated_at     time.Time
	Created_at     time.Time
}

func AllTracks() []Track {
	rows, _, err := db.Con.Query("select * from tracks")
	if err != nil {
		panic(err)
	}

	return RowsToTracks(rows)
}

func (t Track) FromId(Id int) (Track, error) {
	rows, _, err := db.Con.Query("select * from tracks where `id` = '%d'", Id)
	if err != nil {
		return t, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return t, errors.New("Track not found")
	}

	return RowsToTracks(rows)[0], nil
}

func (t Track) Insert() error {
	fmt.Printf("Inserting New Track %s \n", t.Title)

	stmt, err := db.Con.Prepare("insert into tracks (`title`, `rank`, `thumbnail`, `favorite_count`, `play_count`, `artist_id`, `created_at`, `updated_at`) values (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("Error Querying the Database")
	}
	stmt.Exec(t.Title, t.Rank, t.Thumbnail, 0, 0, t.Artist_id, time.Now(), time.Now())

	return nil
}

func (t Track) Save() error {
	fmt.Printf("Updating Track %s \n", t.Title)

	if t.Id < 1 {
		return errors.New("Invalid ID for Track")
	}
	stmt, err := db.Con.Prepare("update tracks set `title`=?, `rank`=?, `thumbnail`=?, `favorite_count`=?, `play_count`=?, `artist_id`=?, `updated_at`=? where `id`=?")
	if err != nil {
		panic(err)
	}
	stmt.Exec(t.Title, t.Rank, t.Thumbnail, t.Favorite_count, t.Play_count, t.Artist_id, time.Now(), t.Id)

	return nil
}

func (t Track) Delete() error {
	fmt.Printf("Deleting Track %s \n", t.Title)
	stmt, err := db.Con.Prepare("delete from tracks where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(t.Id)

	return nil
}

func RowsToTracks(rows []mysql.Row) []Track {
	var tracks = []Track{}
	for _, row := range rows {
		tracks = append(tracks, Track{}.FromRow(row))
	}
	return tracks
}

func (t Track) FromRow(row mysql.Row) Track {
	t.Id = row.Int(0)
	t.Title = row.Str(1)
	t.Rank = row.Int(2)
	t.Thumbnail = row.Str(3)
	t.Favorite_count = row.Int(5)
	t.Play_count = row.Int(6)
	t.Artist_id = row.Int(7)
	t.Updated_at = row.Localtime(8)
	t.Created_at = row.Localtime(9)
	return t
}
