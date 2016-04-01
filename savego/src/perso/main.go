package main

/*.................:::::::::::::::: REMINDYOU .................::::::::::::::::

1) (KO = SOME errors | OK = NOT errors)
2) my english fuck you. (all)

*/

import (
        "github.com/zenazn/goji"
        //_ "github.com/lib/pq"
        _ "github.com/go-sql-driver/mysql"
        "net/http"
        "github.com/zenazn/goji/web"
)

//structure for manage all db
var gest Gestion

func CookieExist(name string, r *http.Request) bool {
    _, err := r.Cookie(name)
    if err != nil {
        return false
    }
    return true
}

func Hihi(c web.C, w http.ResponseWriter, r *http.Request) {
    if CookieExist("token_session", r) {
        Render_html_page(w, "index.html")
    } else {
        Render_html_page(w, "signin_signup.html")
    }
}

func main() {
        //init DB user and DB logging
        gest.Db = InitDbUser()
        gest.DbLog = InitDbLog()
        
        //close all db at end of program 
        defer gest.Db.Close()
        defer gest.DbLog.Close()
       
        DbCreateUser(gest.Db)
        DbCreateLogging(gest.DbLog)
        
        
        // START ROUTING API //
       
        //this is just a try to show it work :)
        goji.Get("/api", IndexShit)
        
        goji.Post("/api/session/create", CreateSession)

        goji.Post("/api/user/create", CreateProfile)


        goji.Get("/api/user/profile/me", MyProfile)
        goji.Get("/api/user/profile/:login", GetProfile)
        goji.Put("/api/user/profile/password", EditPassword)
        goji.Put("/api/user/profile/login", EditLogin)
        goji.Put("/api/user/profile/email", EditEmail)
        goji.Put("/api/user/profile/name", EditName)
        goji.Put("/api/user/profile/bio", EditBio)
        goji.Put("/api/user/profile/orientation", EditOrientation)
        goji.Put("/api/user/profile/birth_date", EditBirthDate)

        goji.Get("/api/user/block", ListBlock)
        goji.Post("/api/user/block/:login", BlockUser)
	    goji.Delete("/api/user/block/:login", UnblockUser)
        
        goji.Post("/api/user/report/:login", ReportUser)
         
        goji.Post("/api/user/tag", AddTag)
        goji.Get("/api/user/tag/:login", ListTag)
        goji.Delete("/api/user/tag/:id", DeleteTag)
        
        goji.Post("/api/user/image", AddImage)
        goji.Get("/api/user/image", ListImage)
        goji.Delete("/api/user/image/:id", DelImage)
        
        goji.Post("/api/user/like/:login", Like)
        goji.Delete("/api/user/like/:login", DelLike)
        goji.Get("/api/user/like", Likes)
        
        goji.Get("/api/user/match", MyMatchs)
        
        goji.Get("/api/user/history", GetVisiteHistory)
        
        //END OF ROUTING API//
        
        
        //Start of web version//
        goji.Get("/", Hihi)
        goji.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
        goji.Get("/assets/view/*", http.StripPrefix("/assets/view/", http.FileServer(http.Dir("./assets/view"))))
        goji.Get("/Images/*", http.StripPrefix("/Images/", http.FileServer(http.Dir("./Images"))))
        goji.Serve()
}