package main

import (
    "net/http"
    "github.com/zenazn/goji/web"
    "strings"
)


//check if userId have blocked thisUser, return true if blocked
func IfIHaveBlockedThisUser(userId, thisUser int) bool {
    
    rows, err := gest.Db.Query("Select 1 from block where blocked_by_id_user=? AND user_id=?", userId, thisUser)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}


//return a []Blocked who contains all users blocked by userId
func GetListUserBlockedById(userId int) []Blocked  {
    var listBlock []Blocked 
    
    rows, err := gest.Db.Query("Select user_id from block where blocked_by_id_user =?", userId)
    LogFatalError(err)
    for rows.Next() {
        var id int
        tmp := Blocked{}
        rows.Scan(&id)
        tmp.Id = id
        tmp.Login = GetLogin(id)
        tmp.Name = GetName(id)
        listBlock = append(listBlock, tmp)
    }
    return listBlock
}

//GET request : return in json a list of user blocked by current user
func ListBlock(c web.C, w http.ResponseWriter, r *http.Request) {
    
    //get token in url
    var token = r.FormValue("token")
    var userId int
    
    ret := ListBlocked{}
    
    //init my empty struct who contains an array string for login user
    
    //check if token exist
    if CheckToken(token, w) == false {
        return
    }
    //check if token matches with an user Id 
    if userId = GetUserIdWithToken(token, w); userId == -1 {
        return
    }
    
    ret.PeopleBlocked = GetListUserBlockedById(userId)
    RenderJSON(w, ret, http.StatusOK)
}

//POST request : Block an user
func BlockUser(c web.C, w http.ResponseWriter, r *http.Request) {
    //i get user login in url
    UserToShow := c.URLParams["login"]
    strings.Replace(UserToShow, " ", "", -1)
    
    var userId int
    //Get session token in param of url
    //get token passed in url param
    
    var token = r.FormValue("token")
    //check if token exisst
    if CheckToken(token, w) == false {
        return
    }
    //check if token matches with an user
    if userId = GetUserIdWithToken(token, w); userId == -1 {
        return
    }
    userIdBlocked := GetIdByLogin(UserToShow, w)

    if IfIHaveBlockedThisUser(userId, userIdBlocked) == false {
        _, err := gest.Db.Exec("Insert into block(blocked_by_id_user, user_id, create_at) values (?, ?, ?)", userId, userIdBlocked, TimeNowString())
        LogFatalError(err)
        _,err = gest.Db.Exec("Delete from likes where liked_by_id_user=? AND id_user=?", userId, userIdBlocked)
        LogFatalError(err)
        _,err = gest.Db.Exec("Delete from likes where id_user=? AND liked_by_id_user=?", userId, userIdBlocked)
        LogFatalError(err)
        RenderJSON(w, RenderStructOk(), http.StatusOK)
    } else {
        RenderJSON(w, RenderStructError("Block", "you have already blocked this user"), http.StatusBadRequest)
    }
}

// Delete request : unblock an user
func UnblockUser(c web.C, w http.ResponseWriter, r *http.Request) {
    //i get user login in url
    UserToShow := c.URLParams["login"]
    strings.Replace(UserToShow, " ", "", -1)
    
    var userId int
    //Get session token in param of url
    //get token passed in url param
    
    var token = r.FormValue("token")
    //check if token exisst
    if CheckToken(token, w) == false {
        return
    }
    //check if token matches with an user
    if userId = GetUserIdWithToken(token, w); userId == -1 {
        return
    }
    userIdUnblocked := GetIdByLogin(UserToShow, w)
    if IfIHaveBlockedThisUser(userId, userIdUnblocked) == true {
        _, err := gest.Db.Exec("Delete from block where blocked_by_id_user=? AND user_id=?", userId, userIdUnblocked)
        LogFatalError(err)
        RenderJSON(w, RenderStructOk(), http.StatusOK)
    } else {
        RenderJSON(w, RenderStructError("Block", "You cannot unblock a people you don't block"), http.StatusBadRequest)
    }   
}

//check if userId have ever reported userIdReported, return true if true
func IfUserHaveEverReportUser(userId, userIdReported int) bool {
    
    rows, err := gest.DbLog.Query("Select 1 from report where id_user=? AND user_reported=?", userId, userIdReported)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}

//Post method : report an user
func ReportUser(c web.C, w http.ResponseWriter, r *http.Request) {
    //i get user login in url
    UserToShow := c.URLParams["login"]
    strings.Replace(UserToShow, " ", "", -1)
    
    var userId int
    //Get session token in param of url
    //get token passed in url param
    var motif = r.FormValue("cause")
    var token = r.FormValue("token")
    //check if token exist
    if CheckToken(token, w) == false {
        return
    }
    //check if token matches with an user
    if userId = GetUserIdWithToken(token, w); userId == -1 {
        return
    }
    userIdReported := GetIdByLogin(UserToShow, w)
    
    if IfUserHaveEverReportUser(userId, userIdReported) == false {
        _, err := gest.DbLog.Exec("Insert into report(cause, id_user, user_reported, create_at) values (?, ?, ?, ?)", motif, userId, userIdReported, TimeNowString())
        LogFatalError(err)
        RenderJSON(w, RenderStructOk(), http.StatusOK)
    } else {
        RenderJSON(w, RenderStructError("Report", "You ever reported this user"), http.StatusBadRequest)
    }
}


