package controllers

import (
    
    "3YFram/helpers"
    "3YFram/models"
    "fmt"
)

type Post struct {
    Title string
    Content string
    Important int
    PostModel models.Post
    PostsModel models.PostList
    user User
    UserId int
    Login string
    Token string
    Id int
}


func (post *Post) CheckPostRequest() []string {
    var errr []string
    
    post.user.model.SessionToken = post.Token
    
    if post.user.model.GetUserBySessionToken() == -1 {
        helpers.AddStringInArray(&errr, "Votre token est mauvais ou a expiré, veuillez vous reconnecter.")
    }
    post.UserId = post.user.model.Id
    return errr
}

func (post *Post) GetPostById() models.ListError {
    errr := post.CheckPostRequest()
    if len(errr) == 0 {
        post.PostModel.Id = post.Id
        errr = post.PostModel.GetPostById(post.user.model.Id)
    }
     return models.RenderErrorList(errr) 
}

func (post *Post) GetUserPosts() models.ListError {
    errr := post.CheckPostRequest()
    if len(errr) == 0 {
        fmt.Println(post.Login)
        errr = post.PostsModel.GetUserPosts(post.Login, post.user.model.Login)
    }
    return models.RenderErrorList(errr)
}

func (post *Post) GetOwnPosts() models.ListError {
    errr := post.CheckPostRequest()
    if len(errr) == 0 {
        fmt.Println(post.Login)
        errr = post.PostsModel.GetUserPosts(post.user.model.Login, post.user.model.Login)
    }
    return models.RenderErrorList(errr)
}
 


func (post *Post) CreatePost() models.ListError {
    var errr []string
    
    errr = post.CheckPostRequest()
    if len(errr) != 0 {
        return models.RenderErrorList(errr)
    }
    
    if post.Important < 1 || post.Important > 3 {
        helpers.AddStringInArray(&errr, "Desolé les sujets n'ont que 3 degré d'importance classés entre 1 et 3.")
    } 
    if len(post.Title) < 3 {
        helpers.AddStringInArray(&errr, "Le titre doit contenir au minimum 3 caractères.")
    }
    if len(post.Content) < 15 {
        helpers.AddStringInArray(&errr, "Le contenu du sujet doit contenir minimum 15 caractères")
    }
    if len(errr) == 0 {
        post.PostModel.Title = post.Title
        post.PostModel.Content = post.Content
        post.PostModel.Important = post.Important
        post.PostModel.UserId = post.UserId
        errr2 := post.PostModel.CreatePost()
        return models.RenderErrorList(errr2)
    }
    return models.RenderErrorList(errr)
}

func (post *Post) DeletePost() models.ListError {
    var errr []string

    errr = post.CheckPostRequest()
    if len(errr) != 0 {
        return models.RenderErrorList(errr)
    }
    post.PostModel.Id = post.Id
    post.PostModel.UserId = post.UserId
    return models.RenderErrorList(post.PostModel.DeletePost())
}