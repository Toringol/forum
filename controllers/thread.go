package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Toringol/forum/database"
	"github.com/Toringol/forum/models"
	"github.com/Toringol/forum/services"
	"github.com/astaxie/beego"
	"github.com/lib/pq"
)

// custom controller
type ThreadController struct {
	beego.Controller
}

// @Title GetAll
// @Description get Thread from url
// @Success 200 {object} models.Thread
// @router /:slug_or_id/details [post]
func (t *ThreadController) UpdateThread() {
	db := database.GetDataBase()
	slug_or_id := t.GetString(":slug_or_id")
	body := t.Ctx.Input.RequestBody
	thread := &models.Thread{}
	oldthread := &models.Thread{}
	json.Unmarshal(body, thread)
	id, err := strconv.Atoi(slug_or_id)
	if err == nil {
		//thread.ID = id
		oldthread, err = models.GetTreadByID(db, id)
		if oldthread == nil {
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.Data["json"] = &models.Error{"Can't find thread with id: " + strconv.Itoa(id)}
			t.ServeJSON()
			return
		}

	} else {
		//thread.Slug = slug_or_id
		oldthread, err = models.GetThreadBySlug(db, slug_or_id)
		if oldthread == nil {
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.Data["json"] = &models.Error{"Can't find thread with slug: " + slug_or_id}
			t.ServeJSON()
			return
		}
	}
	if thread.Title != "" {
		oldthread.Title = thread.Title
	}
	if thread.Message != "" {
		oldthread.Message = thread.Message
	}
	err = models.UpdateThread(db, oldthread)
	if err != nil {
		return
	}
	t.Ctx.Output.SetStatus(http.StatusOK)
	t.Data["json"] = oldthread
	t.ServeJSON()
}

// @Title GetThread by slug or id
// @Description get Thread from url
// @Success 200 {object} models.Thread
// @router /:slug_or_id/details [get]
func (t *ThreadController) GetThread() {
	db := database.GetDataBase()
	slug_or_id := t.GetString(":slug_or_id")
	thread := &models.Thread{}
	id, err := strconv.Atoi(slug_or_id)
	if err == nil {
		//thread.ID = id
		thread, err = models.GetTreadByID(db, id)
		if thread == nil {
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.Data["json"] = &models.Error{"Can't find thread with id: " + strconv.Itoa(id)}
			t.ServeJSON()
			return
		}

	} else {
		//thread.Slug = slug_or_id
		thread, err = models.GetThreadBySlug(db, slug_or_id)
		if thread == nil {
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.Data["json"] = &models.Error{"Can't find thread with slug: " + slug_or_id}
			t.ServeJSON()
			return
		}
	}
	t.Ctx.Output.SetStatus(http.StatusOK)
	t.Data["json"] = thread
	t.ServeJSON()
}

// @Title GetThread by slug or id
// @Description get Thread from url
// @Success 200 {object} models.Thread
// @router /:slug_or_id/vote [post]
func (t *ThreadController) CreateVote() {
	db := database.GetDataBase()
	slug_or_id := t.GetString(":slug_or_id")
	vote := &models.Vote{}
	body := t.Ctx.Input.RequestBody
	json.Unmarshal(body, vote)
	thread := &models.Thread{}
	id, err := strconv.Atoi(slug_or_id)
	if err == nil {
		//thread.ID = id
		thread, err = models.GetTreadByID(db, id)
		if thread == nil {
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.Data["json"] = &models.Error{"Can't find thread with id: " + strconv.Itoa(id)}
			t.ServeJSON()
			return
		}

	} else {
		//thread.Slug = slug_or_id
		thread, err = models.GetThreadBySlug(db, slug_or_id)
		if thread == nil {
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.Data["json"] = &models.Error{"Can't find thread with slug: " + slug_or_id}
			t.ServeJSON()
			return
		}
	}
	user, err := models.GetUserByNickname(db, vote.Nickname)
	if err != nil {
		log.Printf("PATH: %v, error: %v", t.Ctx.Input.URI(), err)
		return
	}
	if user == nil {
		t.Data["json"] = &models.Error{"Can't find user with nickname " + vote.Nickname}
		t.Ctx.Output.SetStatus(http.StatusNotFound)
		t.ServeJSON()
		return
	}
	vote.Thread = thread.ID
	//fmt.Println("___________________________________")
	//fmt.Println(vote)
	//fmt.Println("vote voice", vote.Voice)
	//fmt.Println("___________________________________")
	err = models.CreateVote(db, vote)
	if pgerr, ok := err.(*pq.Error); ok {
		//fmt.Printf("%v\n", pgerr)
		//fmt.Printf("%#v\n", pgerr.Code)
		if pgerr.Code == "23505" {
			voice, _ := models.UpdateVote(db, vote)
			if voice != 0 {
				thread.Votes += 2 * vote.Voice
			}
			t.Ctx.Output.SetStatus(http.StatusOK)
			t.Data["json"] = thread
			t.ServeJSON()
			return
		}
	}
	thread.Votes += vote.Voice
	t.Ctx.Output.SetStatus(http.StatusOK)
	t.Data["json"] = thread
	t.ServeJSON()
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// @Title GetThread by slug or id
// @Description get Thread from url
// @Success 200 {object} models.Thread
// @router /:slug_or_id/create [post]
func (t *ThreadController) CreatePosts() {
	currentTime := time.Now().Truncate(time.Microsecond)
	//currentTime := time.Now().Round(time.Microsecond)
	//fmt.Println("_____________________________________________________________")
	//fmt.Println("_____________________________________________________________")
	//fmt.Printf("______________________________%v______________________________\n", currentTime)
	//fmt.Println("_____________________________________________________________")
	//fmt.Println("_____________________________________________________________")
	db := database.GetDataBase()
	body := t.Ctx.Input.RequestBody
	slug_or_id := t.GetString(":slug_or_id")
	posts := make([]*models.Post, 0)
	json.Unmarshal(body, &posts)
	id, err := strconv.Atoi(slug_or_id)
	thread := &models.Thread{}
	if err == nil {
		thread, err = models.GetTreadByID(db, id)
		if thread == nil {
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.Data["json"] = &models.Error{"Can't find thread with id: " + strconv.Itoa(id)}
			t.ServeJSON()
			return
		}
	} else {
		thread, err = models.GetThreadBySlug(db, slug_or_id)
		if thread == nil {
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.Data["json"] = &models.Error{"Can't find thread with slug: " + slug_or_id}
			t.ServeJSON()
			return
		}
	}
	ids, err := models.GetPostsIDByThreadID(db, thread.ID)
	//fmt.Println("len posts:",len(posts))
	for _, post := range posts {
		if post.Parent != 0 && !contains(ids, post.Parent) {
			t.Ctx.Output.SetStatus(http.StatusConflict)
			t.Data["json"] = &models.Error{"post parent was created in another thread"}
			t.ServeJSON()
			return
		}

		post.Thread = thread.ID
		post.Forum = thread.Forum
		post.Created = currentTime
		//fmt.Println("post.Thread",post.Thread)
		//fmt.Println("post.Forum",post.Forum)
		user, err := models.GetUserByNickname(db, post.Author)
		if err != nil {
			log.Printf("PATH: %v, error: %v", t.Ctx.Input.URI(), err)
			return
		}
		if user == nil {
			t.Data["json"] = &models.Error{"Can't find user with nickname " + post.Author}
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.ServeJSON()
			return
		}
	}
	//fmt.Println("____________________")
	//fmt.Println("CHECK POSTS")
	//fmt.Println(posts)
	//fmt.Println("____________________")
	if len(posts) == 0 {
		//post := &models.Post{}
		//post.Thread = thread.ID
		//post.Forum = thread.Forum
		//post.Created = currentTime
		//db.QueryRow(`INSERT INTO posts (forum, thread, path) VALUES($1, $2, $3) RETURNING id`, post.Forum, post.Thread, pq.Array(post.Path)).Scan(&post.Id)
	} else {
		ids, err := models.CreatePosts(db, posts)
		if err != nil {
			funcname := services.GetFunctionName()
			//log.Println("_____________________________________")
			//log.Println("_____________________________________")
			//log.Println("_____________________________________")
			//log.Println(t.Ctx.Input.URI())
			log.Printf("Function: %s, Error: %v", funcname, err)
			//log.Println(string(t.Ctx.Input.RequestBody))
			//log.Println("_____________________________________")
			//log.Println("_____________________________________")
			//log.Println("_____________________________________")
		}
		for i, ID := range ids {
			//posts[i].Created = times[i]
			posts[i].Id = ID
		}
	}
	t.Ctx.Output.SetStatus(http.StatusCreated)
	t.Data["json"] = posts
	t.ServeJSON()
}

// @Title GetThread by slug or id
// @Description get Thread from url
// @Success 200 {object} models.Thread
// @router /:slug_or_id/posts [get]
func (t *ThreadController) GetPosts() {
	db := database.GetDataBase()
	slug_or_id := t.GetString(":slug_or_id")
	limit := t.Ctx.Input.Query("limit")
	since := t.Ctx.Input.Query("since")
	sort := t.Ctx.Input.Query("sort")
	desc := t.Ctx.Input.Query("desc")
	id, err := strconv.Atoi(slug_or_id)
	thread := &models.Thread{}
	if err == nil {
		//thread.ID = id
		thread, err = models.GetTreadByID(db, id)
		if thread == nil {
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.Data["json"] = &models.Error{"Can't find thread with id: " + strconv.Itoa(id)}
			t.ServeJSON()
			return
		}
	} else {
		//thread.Slug = slug_or_id
		thread, err = models.GetThreadBySlug(db, slug_or_id)
		if thread == nil {
			t.Ctx.Output.SetStatus(http.StatusNotFound)
			t.Data["json"] = &models.Error{"Can't find thread with slug: " + slug_or_id}
			t.ServeJSON()
			return
		}
	}
	switch {
	case sort == "flat" || sort == "":
		lastIndex := 2
		cmp := ""
		addlimit := ""
		addSince := ""
		args := make([]interface{}, 0, 3)
		args = append(args, thread.ID)
		if desc == "false" || desc == "" {
			desc = "ASC"
			cmp = ">"
		} else {
			desc = "DESC"
			cmp = "<"
		}
		if since != "" {
			addSince = fmt.Sprintf("and id "+cmp+" $%d", lastIndex)
			lastIndex += 1
			args = append(args, since)
		}
		if limit != "" {
			addlimit = fmt.Sprintf("limit $%d", lastIndex)
			lastIndex += 1
			args = append(args, limit)
		}
		querystr := fmt.Sprintf("select * from posts where thread = $1 %[1]s ORDER BY id %[2]s %[3]s", addSince, desc, addlimit)
		//fmt.Println("flat sort querystring :", querystr)
		result, err := models.GetPosts(db, querystr, args)
		if err != nil && err != sql.ErrNoRows {
			return
		}
		t.Ctx.Output.SetStatus(http.StatusOK)
		t.Data["json"] = result
		t.ServeJSON()
		return
	case sort == "tree":
		lastIndex := 1
		cmp := ""
		addlimit := ""
		addSince := ""
		args := make([]interface{}, 0, 4)
		if desc == "false" || desc == "" {
			desc = "ASC"
			cmp = ">"

		} else {
			desc = "DESC"
			cmp = "<"
		}
		if since != "" {
			addSince = fmt.Sprintf("JOIN posts AS p2 ON p1.path %s p2.path AND p2.id = $%d where p1.thread =$%d", cmp, lastIndex, lastIndex+1)
			lastIndex += 2
			args = append(args, since)
			args = append(args, thread.ID)
		} else {
			args = append(args, thread.ID)
			addSince = fmt.Sprintf("where p1.thread = $%d", lastIndex)
			lastIndex += 1
		}
		if limit != "" {
			addlimit = fmt.Sprintf("limit $%d", lastIndex)
			lastIndex += 1
			args = append(args, limit)
		}
		querystr := fmt.Sprintf("select p1.* from posts as p1 %s ORDER BY path %s %s", addSince, desc, addlimit)
		//fmt.Println("tree sort querystring :", querystr)
		//fmt.Println("tree sort args", args)
		result, err := models.GetPosts(db, querystr, args)
		if err != nil && err != sql.ErrNoRows {
			return
		}
		t.Ctx.Output.SetStatus(http.StatusOK)
		t.Data["json"] = result
		t.ServeJSON()
		return
	case sort == "parent_tree":
		//lastIndex := 1
		//cmp := ""
		//addlimit := ""
		addSince := ""

		args := make([]interface{}, 0, 3)
		args = append(args, thread.ID)
		args = append(args, limit)
		if desc == "false" || desc == "" {
			desc = "ASC"
			//cmp = ">"
		} else {
			desc = "DESC"
			//cmp = "<"
		}

		if since == "" {
			addSince = "WHERE p.rank <= $2 ORDER BY p.rank,p.path"
		} else {
			addSince = `JOIN sub SubPosts ON SubPosts.id =$3 
				WHERE p.rank <= $2 +SubPosts.rank 
					AND (p.rank > SubPosts.rank OR p.rank=SubPosts.rank and p.path >SubPosts.path) 
				ORDER BY p.rank, p.path`
			args = append(args, since)
		}
		querystr := fmt.Sprintf(`
		WITH sub AS (
    		SELECT p.author, p.created, t.forum, p.id,p.isEdited, p.message, p.parent, p.thread,p.path,
			dense_rank() over (ORDER BY path [1] %[1]s) AS rank
    		FROM posts p JOIN threads t on p.thread= t.id WHERE t.id= $1
		)
		SELECT p.author, p.created, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread, p.path
		FROM sub p %[2]s
    	`, desc, addSince)
		//fmt.Println("query str:", querystr)
		//fmt.Println("query args:", args)
		result, err := models.GetPosts(db, querystr, args)
		if err != nil && err != sql.ErrNoRows {
			log.Println("oops ErrNoRows?", err)
			return
		}
		t.Ctx.Output.SetStatus(http.StatusOK)
		t.Data["json"] = result
		t.ServeJSON()
		return
	}
}
