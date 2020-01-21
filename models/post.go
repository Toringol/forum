package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Toringol/forum/database"
	"github.com/Toringol/forum/services"
	"github.com/lib/pq"
)

type Post struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	Id       int       `json:"id"`
	IsEdited bool      `json:"isEdited"`
	Message  string    `json:"message"`
	Parent   int       `json:"parent"` //идентификтор родительского сообщение
	Thread   int       `json:"thread"`
	Path     []int     `json:"path"`
}

type PostDetails struct {
	Post   *Post   `json:"post"`
	Thread *Thread `json:"thread"`
	Forum  *Forum  `json:"forum"`
	Author *User   `json:"author"`
}

func GetPostsIDByThreadID(db *sql.DB, threadID int) ([]int, error) {
	rows, err := db.Query("select id from posts where thread = $1", threadID)
	result := make([]int, 0)
	defer rows.Close()
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return []int{}, err
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			funcname := services.GetFunctionName()
			log.Printf("Function: %s, Error: %v, while scaning", funcname, err)
			return []int{}, err
		}
		result = append(result, id)
	}
	return result, nil
}

func GetPathById(id int) ([]int, error) {
	db := database.GetDataBase()
	var result []int = make([]int, 0)
	rows, err := db.Query(`SELECT path FROM posts WHERE id = $1`, id)
	defer rows.Close()
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v, while scaning", funcname, err)
		return []int{}, err
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			funcname := services.GetFunctionName()
			log.Printf("Function: %s, Error: %v, while scaning", funcname, err)
			return []int{}, err
		}
		result = append(result, id)
	}
	return result, nil
}

func CreatePosts(db *sql.DB, posts []*Post) ([]int, error) {
	valueStrings := make([]string, 0, len(posts))
	valueArgs := make([]interface{}, 0, len(posts)*7)
	for i, post := range posts {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7))
		//valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)",i*6+1,i*6+2,i*6+3,i*6+4,i*6+5,i*6+6))
		valueArgs = append(valueArgs, post.Author)
		valueArgs = append(valueArgs, post.Created)
		valueArgs = append(valueArgs, post.Forum)
		valueArgs = append(valueArgs, post.Message)
		valueArgs = append(valueArgs, post.Parent)
		valueArgs = append(valueArgs, post.Thread)
		valueArgs = append(valueArgs, pq.Array(post.Path))
		//fmt.Println("___________________________________")
		//fmt.Println("___________________________________")
		//fmt.Println("___________________________________")
		//fmt.Println("___________________________________")
		//fmt.Printf("POST %v %v %v CREATED at %v\n", post.Author,post.Forum, post.Thread, post.Created)
		//fmt.Println("___________________________________")
		//fmt.Println("___________________________________")
		//fmt.Println("___________________________________")
		//fmt.Println("___________________________________")
	}

	stmt := fmt.Sprintf("INSERT INTO posts (author,created,forum,message,parent,thread,path) VALUES %s returning id", strings.Join(valueStrings, ","))
	//fmt.Println("stmt:",stmt)
	//fmt.Println("valueArgs", valueArgs)
	rows, err := db.Query(stmt, valueArgs...)
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v, while scaning", funcname, err)
		return []int{}, err
		//return []int{},[]time.Time{}, err
	}
	defer rows.Close()
	result := make([]int, 0)
	//timeresult := make([]time.Time,0)
	//fmt.Println("check after Query")
	//fmt.Println(rows)
	for rows.Next() {
		id := 0
		//var t time.Time
		//err = rows.Scan(&id, &t)
		err = rows.Scan(&id)
		//fmt.Println("__________________________________ID")
		//fmt.Println(id)
		//fmt.Println("check after scan")
		if err != nil {
			funcname := services.GetFunctionName()
			log.Printf("Function: %s, Error: %v, while scaning", funcname, err)
			return []int{}, err
			//return []int{},[]time.Time{}, err
		}
		result = append(result, id)
		//timeresult = append(timeresult, t)
	}
	//fmt.Println("RESULT IDS CREATEPOSTS",result)
	if err != nil {
		return []int{}, err
		//return []int{},[]time.Time{}, err
	}
	return result, nil
	//return result,timeresult,nil
}

func GetPosts(db *sql.DB, querystr string, args []interface{}) ([]*Post, error) {
	rows, err := db.Query(querystr, args...)
	defer rows.Close()
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return []*Post{}, err
	}
	result := make([]*Post, 0)
	for rows.Next() {
		post := &Post{}
		var pathstring string
		err = rows.Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &pathstring)
		if err != nil {
			funcname := services.GetFunctionName()
			log.Printf("[SCAN] Function: %s, Error: %v", funcname, err)
			return []*Post{}, err
		}
		IDs := strings.Split(pathstring[1:len(pathstring)-1], ",")
		for _, val := range IDs {
			item, _ := strconv.Atoi(val)
			post.Path = append(post.Path, item)
		}
		result = append(result, post)
	}
	return result, nil
}

func GetPostDetailsByID(db *sql.DB, id string) (*PostDetails, error) {
	//rows, err := db.Query(`SELECT * FROM posts WHERE id = $1`, id)
	rows, err := db.Query(`select * from posts as p join users as u on u.nickname = p.author
join forums as f on p.forum = f.slug
join threads thread on p.thread = thread.id
where p.id = $1`, id)
	defer rows.Close()
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return nil, err
	}
	post := &Post{}
	user := &User{}
	forum := &Forum{}
	thread := &Thread{}
	pathstring := ""
	if rows.Next() == true {
		var slug sql.NullString
		err := rows.Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &pathstring,
			&user.About, &user.Email, &user.Fullname, &user.Nickname,
			&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.Author,
			&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &slug, &thread.Title, &thread.Votes)
		if slug.String != "" {
			thread.Slug = slug.String
		}
		if err != nil {
			funcname := services.GetFunctionName()
			log.Printf("[SCAN] Function: %s, Error: %v", funcname, err)
			return nil, err
		}
		for rows.Next() {
			var slug sql.NullString
			err := rows.Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &pathstring,
				&user.About, &user.Email, &user.Fullname, &user.Nickname,
				&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.Author,
				&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &slug, &thread.Title, &thread.Votes)
			if slug.String != "" {
				thread.Slug = slug.String
			}
			if err != nil {
				funcname := services.GetFunctionName()
				log.Printf("[SCAN] Function: %s, Error: %v", funcname, err)
				return nil, err
			}

		}
	} else {
		return nil, sql.ErrNoRows
	}
	pd := &PostDetails{}
	pd.Post = post
	pd.Author = user
	pd.Forum = forum
	pd.Thread = thread
	return pd, nil
}

func UpdatePosts(db *sql.DB, message string, id string) (*Post, error) {
	tx, err := db.Begin()
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return nil, err
	}
	defer tx.Commit()
	pathstring := ""
	post := &Post{}
	err = tx.QueryRow("select * from posts where id = $1", id).Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &pathstring)
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return nil, err
	}
	if message != "" && post.Message != message {
		row := tx.QueryRow("update posts set message = $1, isEdited = true where id = $2 and message <> $1 returning *", message, id)
		err = row.Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &pathstring)
		if err != nil {
			funcname := services.GetFunctionName()
			log.Printf("Function: %s, Error: %v", funcname, err)
			return nil, err
		}
		return post, nil
	} else {
		return post, err
	}
	return post, nil
}
