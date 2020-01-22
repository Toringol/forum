package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Toringol/forum/database"
	"github.com/Toringol/forum/models"
	"github.com/Toringol/forum/services"
	"github.com/astaxie/beego"
	"github.com/lib/pq"
)

// Operations about Users
type ForumController struct {
	beego.Controller
}

// @Title Post
// @Description create forum
// @Param forum body models.Forum true "profile"
// @Success 201 {object} models.Forum
// @Failure 404 no such user
// @Failure 409 already exists
// @router /create [post]
func (f *ForumController) Post() {
	db := database.GetDataBase()
	body := f.Ctx.Input.RequestBody
	forum := &models.Forum{}
	json.Unmarshal(body, forum)
	user, _ := models.GetUserByNickname(db, forum.Author)
	if user == nil {
		f.Ctx.Output.SetStatus(http.StatusNotFound)
		f.Data["json"] = &models.Error{"Can'f find user with nickname " + forum.Author}
		f.ServeJSON()
		return
	}
	forum.Author = user.Nickname
	err := models.CreateForum(db, forum)
	if pgerr, ok := err.(*pq.Error); ok {
		if pgerr.Code == "23505" {
			f.Ctx.Output.SetStatus(http.StatusConflict)
			forum, err := models.GetForumBySlug(db, forum.Slug)
			if err != nil {
				return
			}
			f.Data["json"] = forum
			f.ServeJSON()
			return
		}
	}
	f.Data["json"] = forum
	f.Ctx.Output.SetStatus(http.StatusCreated)
	f.ServeJSON()
}

// @Title Get
// @Description get forum
// @Param slug path string true "identificator"
// @Success 200 {object} models.Forum
// @Failure 404 no such forum
// @router /:slug/details [Get]
func (f *ForumController) Details() {
	db := database.GetDataBase()
	slug := f.GetString(":slug")
	forum, _ := models.GetForumBySlug(db, slug)
	if forum == nil {
		f.Ctx.Output.SetStatus(http.StatusNotFound)
		f.Data["json"] = &models.Error{"Can'f find forum by slug " + slug}
		f.ServeJSON()
		return
	}

	f.Ctx.Output.SetStatus(http.StatusOK)
	f.Data["json"] = forum
	f.ServeJSON()
}

// @Title Get
// @Description get forum
// @Param slug path string true "identificator"
// @Param slug body models.Thread true "thread"
// @Success 201 {object} models.Thread
// @Failure 404 no such user or forum
// @Failure 409 already exists
// @router /:slug/create [Post]
func (f *ForumController) Create() {
	db := database.GetDataBase()
	forumslug := f.GetString(":slug")
	thread := &models.Thread{Forum: forumslug}
	body := f.Ctx.Input.RequestBody
	json.Unmarshal(body, thread)
	forum, _ := models.GetForumBySlug(db, thread.Forum)
	if forum == nil {
		f.Ctx.Output.SetStatus(http.StatusNotFound)
		f.Data["json"] = &models.Error{"Can'f find forum with slug: " + thread.Forum}
		f.ServeJSON()
		return
	}
	user, _ := models.GetUserByNickname(db, thread.Author)
	if user == nil {
		f.Ctx.Output.SetStatus(http.StatusNotFound)
		f.Data["json"] = &models.Error{"Can'f find user with nickname: " + thread.Author}
		f.ServeJSON()
		return
	}
	thread.Author = user.Nickname
	thread.Forum = forum.Slug
	err := models.CreateThread(db, thread)
	if pgerr, ok := err.(*pq.Error); ok {
		if pgerr.Code == "23505" {
			f.Ctx.Output.SetStatus(http.StatusConflict)
			thr, err := models.GetThreadBySlug(db, thread.Slug)
			if err != nil {
				return
			}
			f.Data["json"] = thr
			f.ServeJSON()
			return
		}
	}
	f.Ctx.Output.SetStatus(http.StatusCreated)
	f.Data["json"] = thread
	f.ServeJSON()
}

// @Title Get
// @Description get forum
// @Param slug path string true "identificator"
// @Param limit query number false "max count threads"
// @Param since query string false "time"
// @Param desc query bool false "sort"
// @Success 200 {object} models.Thread
// @Failure 404 no such forum
// @router /:slug/threads [Get]
func (f *ForumController) Threads() {
	db := database.GetDataBase()
	slug := f.GetString(":slug")
	limit := f.Ctx.Input.Query("limit")
	since := f.Ctx.Input.Query("since")
	desc := f.Ctx.Input.Query("desc")
	forum, _ := models.GetForumBySlug(db, slug)
	if forum == nil {
		f.Ctx.Output.SetStatus(http.StatusNotFound)
		f.Data["json"] = &models.Error{"Can'f find forum by slug: " + slug}
		f.ServeJSON()
		return
	}
	threads, _ := models.GetThreadsByForum(db, slug, limit, since, desc)
	if len(threads) == 0 {
		f.Ctx.Output.SetStatus(http.StatusOK)
		f.Data["json"] = threads
		f.ServeJSON()
		return
	}
	f.Ctx.Output.SetStatus(http.StatusOK)
	f.Data["json"] = threads
	f.ServeJSON()
}

// @Title Get
// @Description get forum
// @Param slug path string true "identificator"
// @Param limit query number false "max count threads"
// @Param since query string false "time"
// @Param desc query bool false "sort"
// @Success 200 {object} models.Thread
// @Failure 404 no such forum
// @router /:slug/users [Get]
func (f *ForumController) Users() {
	db := database.GetDataBase()
	slug := f.GetString(":slug")
	limit := f.Ctx.Input.Query("limit")
	since := f.Ctx.Input.Query("since")
	desc := f.Ctx.Input.Query("desc")
	forum, err := models.GetForumBySlug(db, slug)
	if err != nil {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return
	}
	if forum == nil {
		f.Ctx.Output.SetStatus(http.StatusNotFound)
		f.Data["json"] = &models.Error{"Can'f find forum with slug: " + slug}
		f.ServeJSON()
		return
	}

	result := make([]*models.User, 0)
	args := make([]interface{}, 0, 4)
	args = append(args, forum.Slug)
	lastIndex := 2
	addLimit := ""
	cmp := ""
	addSince := ""
	if desc == "true" {
		desc = "DESC"
		cmp = "<"
	} else {
		desc = "ASC"
		cmp = ">"
	}
	if since != "" {
		addSince = fmt.Sprintf("AND u.nickname %s $%d", cmp, lastIndex)
		args = append(args, since)
		lastIndex++
	}
	if limit != "" {
		addLimit = fmt.Sprintf("limit $%d", lastIndex)
		lastIndex++
		args = append(args, limit)
	}
	queryrow := fmt.Sprintf(`SELECT u.about, u.email, u.fullname,u.nickname
			FROM boost b
			JOIN users u ON b.username = u.nickname
WHERE b.slug = $1 %[1]s ORDER BY nickname %[2]s %[3]s`, addSince, desc, addLimit)
	result, err = models.GetUsers(db, queryrow, args)
	if err != nil && err != sql.ErrNoRows {
		funcname := services.GetFunctionName()
		log.Printf("Function: %s, Error: %v", funcname, err)
		return
	}
	f.Ctx.Output.SetStatus(http.StatusOK)
	f.Data["json"] = result
	f.ServeJSON()
	return
}
