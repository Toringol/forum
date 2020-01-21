package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Toringol/forum/database"
	"github.com/Toringol/forum/models"
	"github.com/astaxie/beego"
)

type PostController struct {
	beego.Controller
}

// @Title Post
// @Description create forum
// @Param id path string true "id"
// @Param related query bool false "related"
// @Success 201 {object} models.Forum
// @Failure 404 no such user
// @Failure 409 already exists
// @router /:id/details [get]
func (p *PostController) Get() {
	db := database.GetDataBase()

	id := p.GetString(":id")
	related := p.Ctx.Input.Query("related")
	infos := strings.Split(related, ",")
	pd, err := models.GetPostDetailsByID(db, id)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if pd == nil {
		p.Ctx.Output.SetStatus(http.StatusNotFound)
		p.Data["json"] = &models.Error{"Can't fild post with id: " + id}
		p.ServeJSON()
		return
	}
	result := &models.PostDetails{}
	for _, v := range infos {
		if v == "user" {
			result.Author = pd.Author
		}
		if v == "forum" {
			result.Forum = pd.Forum
		}
		if v == "thread" {
			result.Thread = pd.Thread
		}
	}
	result.Post = pd.Post
	p.Ctx.Output.SetStatus(http.StatusOK)
	p.Data["json"] = result
	p.ServeJSON()
}

// @Title Post
// @Description create forum
// @Param id path string true "id"
// @Param related query bool false "related"
// @Success 201 {object} models.Forum
// @Failure 404 no such user
// @Failure 409 already exists
// @router /:id/details [post]
func (p *PostController) UpdatePosts() {
	db := database.GetDataBase()
	id := p.GetString(":id")
	body := p.Ctx.Input.RequestBody
	updatepost := &models.Post{}
	json.Unmarshal(body, updatepost)
	//tx, err := db.Begin()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//defer tx.Commit()
	//pathstring:=""
	//post := &models.Post{}
	//err = tx.QueryRow("select * from posts where id = $1", id).Scan(&post.Author,&post.Created,&post.Forum,&post.Id,&post.IsEdited,&post.Message,&post.Parent,&post.Thread,&pathstring)
	//if err != nil {
	//	p.Ctx.Output.SetStatus(http.StatusNotFound)
	//	p.Data["json"] = &models.Error{"can't find post with id: "+ id}
	//	p.ServeJSON()
	//	return
	//}
	//if updatepost.Message != "" && post.Message != updatepost.Message {
	//	post, err = models.UpdatePosts(db, post.Message, id)
	//	if err != nil {
	//		p.Ctx.Output.SetStatus(http.StatusNotFound)
	//		p.Data["json"] = &models.Error{"can't find post with id: "+ id}
	//		p.ServeJSON()
	//		return
	//	}
	//}
	post, err := models.UpdatePosts(db, updatepost.Message, id)
	if post == nil {
		p.Ctx.Output.SetStatus(http.StatusNotFound)
		p.Data["json"] = &models.Error{"can't find post with id: " + id}
		p.ServeJSON()
		return
	}
	if err != nil {
		p.Ctx.Output.SetStatus(http.StatusNotFound)
		p.Data["json"] = &models.Error{"can't find post with id: " + id}
		p.ServeJSON()
		return
	}
	p.Ctx.Output.SetStatus(http.StatusOK)
	p.Data["json"] = post
	p.ServeJSON()
}
