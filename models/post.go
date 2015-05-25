package models

import (
	"errors"
	"fmt"
	db "github.com/AlienStream/Shared-Go/database"
	mysql "github.com/ziutek/mymysql/mysql"
	"time"
)

type Post struct {
	Id                 int
	Title              string
	Number_of_comments int
	Permalink          string
	Thumbnail          string
	Likes              int
	Dislikes           int
	Submitter          string
	Source_id          int
	Embed_url          string
	Is_new             bool
	Posted_at          time.Time
	Updated_at         time.Time
	Created_at         time.Time
}

func AllPosts() []Post {
	rows, _, err := db.Con.Query("select * from posts")
	if err != nil {
		panic(err)
	}

	return RowsToPosts(rows)
}

func (p Post) FromId(Id int) (Post, error) {
	rows, _, err := db.Con.Query("select * from posts where `id`=%s", Id)
	if err != nil {
		return p, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return p, errors.New("Post Not Found")
	}

	return RowsToPosts(rows)[0], nil
}

func (p Post) IsNew() (bool, int) {
	rows, _, err := db.Con.Query("select * from posts where `source_id`=%d and `embed_url` = '%s'", p.Source_id, p.Embed_url)
	if err != nil {
		panic("Error When Querying the database")
	}

	if len(rows) == 0 {
		return true, 0
	}

	return false, RowsToPosts(rows)[0].Id

}

func (p Post) Insert() error {
	fmt.Printf("Inserting New Post %s \n", p.Title)
	stmt, err := db.Con.Prepare("insert into posts (`title`, `number_of_comments`, `permalink`, `thumbnail`, `likes`, `dislikes`, `submitter`, `source_id`, `is_new`, `embed_url`, `posted_at`, `created_at`, `updated_at`) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(p.Title, p.Number_of_comments, p.Permalink, p.Thumbnail, p.Likes, p.Dislikes, p.Submitter, p.Source_id, true, p.Embed_url, p.Posted_at, time.Now(), time.Now())

	return nil
}

func (p Post) Save() error {
	fmt.Printf("Updating Post %s \n", p.Title)
	if p.Id < 1 {
		return errors.New("Invalid ID for Post")
	}
	stmt, err := db.Con.Prepare("update posts set `title`=?, `number_of_comments`=?, `permalink`=?, `thumbnail`=?, `likes`=?, `dislikes`=?, `submitter`=?, `source_id`=?, `embed_url`=?, `posted_at`=?, `updated_at`=? where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(p.Title, p.Number_of_comments, p.Permalink, p.Thumbnail, p.Likes, p.Dislikes, p.Submitter, p.Source_id, p.Embed_url, p.Posted_at, time.Now(), p.Id)

	return nil

}

func (p Post) Delete() error {
	fmt.Printf("Deleting Post %s \n", p.Title)
	if p.Id < 1 {
		return errors.New("Invalid ID for Post")
	}
	stmt, err := db.Con.Prepare("delete from posts where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(p.Id)

	return nil
}

func RowsToPosts(rows []mysql.Row) []Post {
	var posts = []Post{}
	for _, row := range rows {
		posts = append(posts, Post{}.FromRow(row))
	}
	return posts
}

func (p Post) FromRow(row mysql.Row) Post {
	p.Id = row.Int(0)
	p.Title = row.Str(1)
	p.Number_of_comments = row.Int(2)
	p.Permalink = row.Str(3)
	p.Thumbnail = row.Str(4)
	p.Likes = row.Int(5)
	p.Dislikes = row.Int(6)
	p.Submitter = row.Str(7)
	p.Source_id = row.Int(8)
	p.Is_new = row.Bool(9)
	p.Embed_url = row.Str(10)
	p.Posted_at = row.Localtime(11)
	p.Updated_at = row.Localtime(12)
	p.Created_at = row.Localtime(13)
	return p
}
