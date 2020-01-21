package controllers

import (
	"log"
	"net/http"

	"github.com/Toringol/forum/database"
	"github.com/astaxie/beego"
)

type ServiceController struct {
	beego.Controller
}

type status struct {
	Forum  int `json:"forum"`
	Post   int `json:"post"`
	Thread int `json:"thread"`
	User   int `json:"user"`
}

// @Status
// @Get database status
// @Success 200 {object} status
// @router /status [get]
func (s *ServiceController) Status() {
	db := database.GetDataBase()
	status := &status{}
	rows := db.QueryRow(`select count(*) from users`)
	err := rows.Scan(&status.User)
	if err != nil {
		log.Println("[Status] User error:", err)
	}

	rows = db.QueryRow(`select count(*) from forums`)
	err = rows.Scan(&status.Forum)
	if err != nil {
		log.Println("[Status] Forum error:", err)
	}
	rows = db.QueryRow(`select count(*) from threads`)
	err = rows.Scan(&status.Thread)
	if err != nil {
		log.Println("[Status] Thread error:", err)
	}

	rows = db.QueryRow(`select count(*) from posts`)
	err = rows.Scan(&status.Post)
	if err != nil {
		log.Println("[Status] Post error:", err)
	}
	s.Ctx.Output.SetStatus(http.StatusOK)
	s.Data["json"] = status
	s.ServeJSON()
}

// @Status
// @Get database status
// @Success 200 {object} status
// @router /clear [post]
func (s *ServiceController) Clear() {
	db := database.GetDataBase()
	_, err := db.Exec("truncate table users, forums, threads, posts, votes, boost")
	if err != nil {
		log.Println("ERROR CLEARING:", err)
		return
	}
	s.Ctx.Output.SetStatus(http.StatusOK)
}
