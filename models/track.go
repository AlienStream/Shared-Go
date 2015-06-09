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
	Rank           float64
	Title          string
	Thumbnail      string
	Favorite_count int
	Play_count     int
	Channel_id     int
        Content_flags  int
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
	rows, _, err := db.Con.Query("select * from tracks where `id` = %d", Id)
	if err != nil {
		return t, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return t, errors.New("Track not found")
	}

	return RowsToTracks(rows)[0], nil
}

func (t Track) FromTitle(Title string) (Track, error) {
	rows, _, err := db.Con.Query("select * from tracks where `title` = '%s'", db.Con.Escape(Title))
	if err != nil {
		return t, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return t, errors.New("Track not found")
	}

	return RowsToTracks(rows)[0], nil
}

func (t Track) IsNew() bool {
	rows, _, err := db.Con.Query("select * from tracks where `title` = '%s'", db.Con.Escape(t.Title))
	if err != nil {
		panic(err)
	}

	if len(rows) == 0 {
		return true
	}

	return false
}

func (t Track) Insert() error {
	fmt.Printf("Inserting New Track %s \n", t.Title)

	if t.Title == "" {
		return errors.New("Invalid track Title")
	}

	stmt, err := db.Con.Prepare("insert into tracks (`title`, `rank`, `thumbnail`, `favorite_count`, `play_count`, `channel_id`, `content_flags`, `created_at`, `updated_at`) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
                panic(err);
	}
	_, _, err = stmt.Exec(t.Title, t.Rank, t.Thumbnail, 0, 0, t.Channel_id, 0, t.Created_at, time.Now())
	if err != nil {
                panic(err);
        }
	return nil
}

func (t Track) Save() error {
	fmt.Printf("Updating Track %s \n", t.Title)

	if t.Id < 1 {
		return errors.New("Invalid ID for Track")
	}
	stmt, err := db.Con.Prepare("update tracks set `title`=?, `rank`=?, `thumbnail`=?, `favorite_count`=?, `play_count`=?, `channel_id`=?, `content_flags`=?, `updated_at`=? where `id`=?")
	if err != nil {
		panic(err)
	}
	stmt.Exec(t.Title, t.Rank, t.Thumbnail, t.Favorite_count, t.Play_count, t.Channel_id, t.Content_flags, time.Now(), t.Id)

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

func CreateTrackSourcePivot(s Source, t Track) {
	fmt.Printf("Joining Track %s To Source %s \n", t.Title, s.Title)

	rows, _, err := db.Con.Query("select * from source_track where `source_id` = %d and `track_id` = %d", s.Id, t.Id)
	if err != nil {
		panic("Error When Querying the database")
	}

	if len(rows) == 0 {
		stmt, err := db.Con.Prepare("insert into source_track (`source_id`, `track_id`) values (?, ?)")
		if err != nil {
			panic("Error Querying the Database")
		}
		stmt.Exec(s.Id, t.Id)
	}
}

func DeleteTrackSourcePivot(s Source, t Track) {
	fmt.Printf("Joining Track %s To Source %s \n", t.Title, s.Title)

	stmt, err := db.Con.Prepare("delete from source_track where `source_id`=? and `track_id`=?")
	if err != nil {
		panic("Error Querying the Database")
	}
	stmt.Exec(s.Id, t.Id)
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
	t.Rank = row.Float(2)
	t.Thumbnail = row.Str(3)
	t.Favorite_count = row.Int(4)
	t.Play_count = row.Int(5)
	t.Channel_id = row.Int(6)
	t.Content_flags = row.Int(7)
	t.Updated_at = row.Localtime(8)
	t.Created_at = row.Localtime(9)
	return t
}
