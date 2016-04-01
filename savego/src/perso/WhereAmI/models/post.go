package models

import (
    "github.com/jinzhu/now"
    "time"
    "3YFram/helpers"
    "strings"
)

type Post struct {
    Content string
    Title string
    Important int
    UserId int
    CreateDate string
    Id int
    Status string
}

type PostList struct {
    Liste []Post
    Status string
}

func (post *Post) IsPostIdBelongAtUserId() bool {
    rows, err := gest.db.Query("select 1 from post where id=$1 AND user_id=$2", post.Id, post.UserId)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}

func (post *Post) DeletePost() []string {
    var errr []string 
    
    if post.IsPostIdBelongAtUserId() == false {
        helpers.AddStringInArray(&errr, "Vous ne pouvez supprimer que vos posts.")
    }
    
    if len(errr) == 0 {
        _, err := gest.db.Exec("Delete from post where id=$1", post.Id)
        LogFatalError(err)
    }
    return errr
}

func (post *Post) HasEver3PostInLast24Hours() bool {
    
    rows, err := gest.db.Query("Select create_at from post where user_id=$1 order by create_at desc limit 3", post.UserId)
    LogFatalError(err)
    var i = 0
    for rows.Next() {
        var createDate string
        rows.Scan(&createDate)
        
        layout := "2006-01-02 15:04:05"
        timepost, err := time.Parse(layout, createDate)
        LogFatalError(err)
        beginday := now.BeginningOfDay()
        if timepost.After(beginday) {
            i++;
        }
    }
    if i == 3 {
        return true
    }
    return false
    
}

func (post *Post) IsPostExist() bool {
    rows, err := gest.db.Query("Select 1 from post where id=$1", post.Id)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}

func (post *Post) GetInformationByPostId() int {
    
    if post.IsPostExist() == false {
        return -1
    }
    rows, err := gest.db.Query("Select user_id, important, content, title, create_at from post where id=$1", post.Id)
    LogFatalError(err)
    var id int
    var important int
    var content string
    var title string
    var createAt string
    
    for rows.Next() {
        rows.Scan(&id, &important, &content, &title, &createAt)
    }
    post.UserId = id
    post.Content = content
    post.Important = important
    post.Title = title
    post.CreateDate = createAt
    return id
}


func (post *Post) GetPostById(UserIdRequest int)[]string {
    var errr []string
    
    if post.GetInformationByPostId() == -1 {
        helpers.AddStringInArray(&errr, "Ce post n'existe pas.")
    } else {
        if UserIdRequest != post.UserId && post.Important != 1 {
            helpers.AddStringInArray(&errr, "Vous n'avez pas l'autorisation de voir ce post")
        }
    }
    return errr
}

func (postlist *PostList) GetUserPosts(login, LoginUserRequest string) []string {
    var errr []string
    
    user := User{}
    user.Login = login
    if user.IsUserLoginExist() == false {
        helpers.AddStringInArray(&errr, "Cette utilisateur n'existe pas.")
    } else {
        if strings.Compare(login, LoginUserRequest) == 0 {
            user.GetUserByLogin()
            rows, err := gest.db.Query("Select id from post where user_id=$1", user.Id)
            LogFatalError(err)
            for rows.Next() {
                post := Post{}
                rows.Scan(&post.Id)
                if post.GetInformationByPostId() != -1 {
                    postlist.Liste = append(postlist.Liste, post)
                }
            }       
        } else {
            user.GetUserByLogin()
            rows, err := gest.db.Query("Select id from post where user_id=$1 And important=1", user.Id)
            LogFatalError(err)
            for rows.Next() {
                post := Post{}
                rows.Scan(&post.Id)
                if post.GetInformationByPostId() != -1 {
                    postlist.Liste = append(postlist.Liste, post)
                }
            }
        }
    }
    return errr
}



func (post *Post) CreatePost() []string {
    var errr []string
    if post.HasEver3PostInLast24Hours() == true {
        helpers.AddStringInArray(&errr, "Vous avez déjà posté 3 messages depuis hier minuit.")
    } else {
        _, err := gest.db.Exec("Insert into post(content, title, important, user_id) values ($1, $2, $3, $4)", post.Content, post.Title, post.Important, post.UserId)    
        LogFatalError(err)
    }
    return errr
}