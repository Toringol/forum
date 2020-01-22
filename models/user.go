package models

import (
	"database/sql"
	"log"

	"github.com/Toringol/forum/services"
)

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
}

func (user *User) scanUser(rows *sql.Rows) error {
	err := rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		log.Println("Error in scanUser:", err)
		return err
	}
	return nil
}

func CreateUser(db *sql.DB, user *User) error {
	_, err := db.Exec("INSERT INTO users (about, email, fullname, nickname) VALUES ($1, $2, $3, $4)", user.About, user.Email, user.Fullname, user.Nickname)
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, User: %v, Error: %v", funcname, user, err)
		return err
	}
	return nil
}

func GetUserByNickname(db *sql.DB, nick string) (*User, error) {
	row := db.QueryRow("select * from users where nickname = $1", nick)
	user := &User{}
	err := row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return user, nil
	default:
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return nil, err
	}
	return user, nil
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	row := db.QueryRow("select * from users where email = $1", email)
	user := &User{}
	err := row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return user, nil
	default:
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Email: %s, Error: %v", funcname, email, err)
		return nil, err
	}
	return user, nil
}

func GetUserWithEmailOrNickname(db *sql.DB, email, nickname string) ([]*User, error) {
	result := make([]*User, 0)
	rows, err := db.Query("select * from users where email = $1 or nickname = $2", email, nickname)
	defer rows.Close()
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Email: %s, Nickname: %s, Error: %v", funcname, email, nickname, err)
		return result, err
	}
	for rows.Next() {
		user := &User{}
		err := user.scanUser(rows)
		if err != nil {
			return result, err
		}
		result = append(result, user)
	}
	return result, nil
}

func UpdateUserByNickname(db *sql.DB, nickname string, user User) error {
	_, err := db.Exec("update users set about = $1, email = $2, fullname = $3 where nickname =$4", user.About, user.Email, user.Fullname, nickname)
	if err != nil {
		return err
	}

	return nil
}

func GetUsers(db *sql.DB, querystr string, args []interface{}) ([]*User, error) {
	rows, err := db.Query(querystr, args...)
	defer rows.Close()
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return []*User{}, err
	}
	result := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			funcname := services.GetFunctionName()
			log.Printf("[SCAN] Function: %s, Error: %v", funcname, err)
			return []*User{}, err
		}

		result = append(result, user)
	}
	return result, nil
}
