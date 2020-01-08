package models

import (
	"database/sql"
	"log"

	"github.com/Toringol/forum/services"
)

type Forum struct {
	Posts   int    `json:"posts"`
	Slug    string `json:"slug"`
	Threads int    `json:"threads"`
	Title   string `json:"title"`
	Author  string `json:"user"`
}

func (f *Forum) scanForum(rows *sql.Rows) error {
	if rows.Next() == true {
		err := rows.Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.Author)
		if err != nil {
			log.Println("Error in scanForum:", err)
			return err
		}
		for rows.Next() {
			err := rows.Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.Author)
			if err != nil {
				log.Println("Error in scanForum:", err)
				return err
			}
		}
	} else {
		return sql.ErrNoRows
	}
	return nil
}

func CreateForum(db *sql.DB, forum *Forum) error {
	_, err := db.Exec("INSERT INTO forums (slug,title,author) VALUES ($1, $2, $3)", forum.Slug, forum.Title, forum.Author)
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Forum %v, Error: %v", funcname, forum, err)
		return err
	}
	return nil
}

func GetForumBySlug(db *sql.DB, slug string) (*Forum, error) {
	rows, err := db.Query("select * from forums where slug = $1", slug)
	defer rows.Close()
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return nil, err
	}
	forum := &Forum{}
	err = forum.scanForum(rows)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return forum, nil
	default:
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return nil, err
	}
	//return forum, nil
}
