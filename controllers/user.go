package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Toringol/forum/database"
	"github.com/Toringol/forum/models"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title Post
// @Description create user
// @Param	nickname		path 	string	true	"nickname from uri"
// @Param profile body models.User true "profile"
// @Success 201 {object} models.User
// @Failure 403 :uid is empty
// @router /:nickname/create [post]
func (u *UserController) Post() {
	db := database.GetDataBase()
	body := u.Ctx.Input.RequestBody
	nickname := u.GetString(":nickname")
	user := &models.User{Nickname: nickname}
	json.Unmarshal(body, user)
	result := make([]*models.User, 0)
	result, err := models.GetUserWithEmailOrNickname(db, user.Email, user.Nickname)
	if err != nil {
		log.Printf("PATH: %v, error: %v", u.Ctx.Input.URI(), err)
		return
	}
	if len(result) != 0 {
		u.Ctx.Output.SetStatus(http.StatusConflict)
		u.Data["json"] = result
		u.ServeJSON()
		return
	}
	models.CreateUser(db, user)
	u.Data["json"] = user
	u.Ctx.Output.SetStatus(http.StatusCreated)
	u.ServeJSON()
}

// @Title Post
// @Description user information
// @Param	nickname		path 	string	true	"nickname from uri"
// @Success 200 {object} models.User
// @Failure 404 {object} models.Error
// @router /:nickname/profile [get]
func (u *UserController) ProfileGet() {
	db := database.GetDataBase()
	nickname := u.GetString(":nickname")
	user := &models.User{Nickname: nickname}
	user, err := models.GetUserByNickname(db, nickname)
	if err != nil {
		log.Printf("PATH: %v, error: %v", u.Ctx.Input.URI(), err)
		return
	}
	if user != nil {
		u.Data["json"] = user
		u.Ctx.Output.SetStatus(http.StatusOK)
		u.ServeJSON()
		return
	}
	u.Data["json"] = &models.Error{"Can't find user with nickname " + nickname}
	u.Ctx.Output.SetStatus(http.StatusNotFound)
	u.ServeJSON()
}

// @Title Post
// @Description user information
// @Param	nickname		path 	string	true	"nickname from uri"
// @Param profile body models.User true "profile"
// @Success 200 {object} models.User
// @Failure 404 {object} models.Error
// @router /:nickname/profile [post]
func (u *UserController) ProfilePost() {
	db := database.GetDataBase()
	nickname := u.GetString(":nickname")
	body := u.Ctx.Input.RequestBody
	user, err := models.GetUserByNickname(db, nickname)
	if user == nil {
		u.Ctx.Output.SetStatus(http.StatusNotFound)
		u.Data["json"] = &models.Error{"Can't find user with nickname " + nickname}
		u.ServeJSON()
		return
	}
	oldmail := user.Email

	json.Unmarshal(body, user)
	if err != nil {
		log.Printf("PATH: %v, error: %v", u.Ctx.Input.URI(), err)
		return
	}

	if oldmail != user.Email {
		if checkuser, err := models.GetUserByEmail(db, user.Email); checkuser != nil && err == nil && checkuser.Nickname != nickname {
			u.Ctx.Output.SetStatus(http.StatusConflict)
			u.Data["json"] = &models.Error{"This email is already registered by user: " + checkuser.Nickname}
			u.ServeJSON()
			return
		}
	}
	err = models.UpdateUserByNickname(db, nickname, *user)
	if err != nil {
		log.Printf("PATH: %v, error: %v", u.Ctx.Input.URI(), err)
		return
	}
	u.Data["json"] = user
	u.Ctx.Output.SetStatus(http.StatusOK)
	u.ServeJSON()
}
