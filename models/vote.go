package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Toringol/forum/services"
)

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice, integer"`
	Thread   int    `json:"thread"`
}

func CreateVote(db *sql.DB, vote *Vote) error {
	_, err := db.Exec("insert into votes (nickname, voice, thread) values ($1,$2,$3)", vote.Nickname, vote.Voice, vote.Thread)
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return err
	}
	return nil
}

func UpdateVote(db *sql.DB, vote *Vote) (int, error) {
	var voice int = 0
	rows := db.QueryRow("update votes set voice = $1 where nickname = $2 and thread = $3 and voice <> $1 returning voice", vote.Voice, vote.Nickname, vote.Thread)
	rows.Scan(&voice)
	fmt.Println(voice)
	return voice, nil
}
