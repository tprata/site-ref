package main 

import (
       "3YFram/controllers"
       "3YFram/helpers"
       "3YFram/models"
        //_ "github.com/lib/pq"
       _ "github.com/go-sql-driver/mysql"
        "github.com/zenazn/goji/web"
        "github.com/zenazn/goji"
        "net/http"
        "net"
        "strings"
        "strconv"
)


func Signup(c web.C, w http.ResponseWriter, r *http.Request) {
    user := controllers.User{}
    helpers.RenderJSON(w, user.CreateAccount(r.FormValue("login"), r.FormValue("email"), r.FormValue("password"), r.FormValue("password2")), http.StatusOK)
}

func Signin(c web.C, w http.ResponseWriter, r *http.Request) {
    session := controllers.Session{}
    ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
    err := session.CreateSession(r.FormValue("login"), r.FormValue("password"), ipAddress)
    if session.Token != "" {
        helpers.RenderJSON(w, session, http.StatusOK)
    } else {
        helpers.RenderJSON(w, err, http.StatusOK)
    }
}

func Comfirmation(c web.C, w http.ResponseWriter, r *http.Request) {
    user := controllers.User{}
    user.VerifiedToken = r.FormValue("verifiedtoken")
    errr := user.ComfirmAccount()
    if strings.Compare(errr.Status, "KO") == 0 {
        helpers.RenderJSON(w, errr, http.StatusOK)
    } else {
        helpers.RenderOk(w)
    }
}

func CreatePost(c web.C, w http.ResponseWriter, r *http.Request) {
    post := controllers.Post{}
    
    post.Title = r.FormValue("title")
    post.Content = r.FormValue("content")
    post.Token = r.FormValue("token")
    post.Important, _ = strconv.Atoi(r.FormValue("important"))
    helpers.RenderJSON(w, post.CreatePost(), http.StatusOK)
}

func DeletePost(c web.C, w http.ResponseWriter, r *http.Request) {
    post := controllers.Post{}
    
    post.Id,  _ = strconv.Atoi(c.URLParams["id"])
    post.Token = r.FormValue("token")
    helpers.RenderJSON(w, post.DeletePost(), http.StatusOK)
}

func GetPost(c web.C, w http.ResponseWriter, r *http.Request) {
    post := controllers.Post{}
    
    post.Id,  _ = strconv.Atoi(c.URLParams["id"])
    post.Token = r.FormValue("token")
    err := post.GetPostById()
    if err.Status == "KO" {
        helpers.RenderJSON(w, err, http.StatusOK)
    } else {
        post.PostModel.Status = "OK"
        helpers.RenderJSON(w, post.PostModel, http.StatusOK)
    }
}

func GetUserPosts(c web.C, w http.ResponseWriter, r *http.Request) {
    post := controllers.Post{}
    
    post.Login = c.URLParams["login"]
    post.Token = r.FormValue("token")
    err := post.GetUserPosts()
    if err.Status == "KO" {
        helpers.RenderJSON(w, err, http.StatusOK)
    } else {
        post.PostsModel.Status = "OK"
        helpers.RenderJSON(w, post.PostsModel, http.StatusOK)
    }
}

func GetOwnPosts(c web.C, w http.ResponseWriter, r *http.Request) {
    post := controllers.Post{}
    
    post.Token = r.FormValue("token")
    err := post.GetOwnPosts()
    if err.Status == "KO" {
        helpers.RenderJSON(w, err, http.StatusOK)
    } else {
        post.PostsModel.Status = "OK"
        helpers.RenderJSON(w, post.PostsModel, http.StatusOK)
    }
}


func Index(c web.C, w http.ResponseWriter, r *http.Request) {
    helpers.RenderHTMLPage(w, "index.html")
}

func routing() {
    
    ///// Website /////
    
    goji.Get("/", Index)    
    
    
    
    /////API///////
    
    //inscription
    goji.Post("/api/signup", Signup)
    
    //connexion
    goji.Post("/api/signin", Signin)
    
    
    
    
    //creation de poste
    goji.Post("/api/post", CreatePost)
   
    //obtenir ses propres posts
    goji.Get("/api/post", GetOwnPosts)
       
    //obtenir un post par son id
    goji.Get("/api/post/:id", GetPost)

    //obtenir une liste de post par le login    
    goji.Get("/api/:login/post", GetUserPosts)
    

    //suppresion d'un poste
    goji.Delete("/api/post/:id", DeletePost)
    
    
    
    
    //comfirmation de l'email si le fichier app.json est configuré pour
    if models.Conf.Status == "prod" || models.Conf.EmailCheck == true {
        goji.Get("/comfirmation", Comfirmation)
    }
}

func main() {
    models.InitDB()
    routing()
    goji.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    goji.Serve()
}