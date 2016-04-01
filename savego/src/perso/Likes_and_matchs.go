package main 


import (
    "net/http"
    "strings"
    "github.com/zenazn/goji/web"
)


//return list of people liked by user
func GetListPeopleILike(userId int) []LikePeople{
    
    //init array of struct peopleILike who contains for each user : bio, id, login, name, orientation.
    var peopleILike []LikePeople
    
    //query to have all user, current user like
    rows, err := gest.Db.Query("Select id_user from likes where liked_by_id_user=?", userId)
    LogFatalError(err)
    //for each user
    for rows.Next() {
        var id int
        
        //init tmp struct LikePeople to get information of actually user in while
        var tmp LikePeople
        
        //get id
        rows.Scan(&id)
        tmp.Id = id
        
        //get bio of user
        tmp.Bio = GetBio(id)
        
        //get login of user
        tmp.Login = GetLogin(id)
        
        //get name of user
        tmp.Name = GetName(id)
        
        //get orientation of user
        tmp.Orientation = GetOrientation(id)
        
        //append tmp in array struct LikePeople
        peopleILike = append(peopleILike, tmp)
    }
    return peopleILike
}

//return list of people who like current user
func GetListPeopleWhoLikeMe(userId int) []LikePeople {
    
    //init array of struct peopleILike who contains for each user : bio, id, login, name, orientation.
    var peopleILike []LikePeople
    
    //query to have all user, current user like
    rows, err := gest.Db.Query("Select liked_by_id_user from likes where id_user=?", userId)
    LogFatalError(err)
    for rows.Next() {
        var id int
        
         //init tmp struct LikePeople to get information of actually user in while
        var tmp LikePeople
        
        //get id
        rows.Scan(&id)
        
        tmp.Id = id
        
        //get bio of user
        tmp.Bio = GetBio(id)
        
        //get login of user
        tmp.Login = GetLogin(id)
        
        //get name of user
        tmp.Name = GetName(id)
        
        //get orientation of user
        tmp.Orientation = GetOrientation(id)
        
        //append tmp in array struct LikePeople
        peopleILike = append(peopleILike, tmp)
    }
    return peopleILike
}


//METHOD GET : get a list of people liked by user and list of like by people to current user
func Likes(c web.C, w http.ResponseWriter, r *http.Request) {
    
    //get token in url
    var token = r.FormValue("token")
    var userId int
    
    //init my empty struct who contains an array string for login user
    myStruct := AllLikes{}
    
    //check if token exist
    if CheckToken(token, w) == false {
        return
    }
    //check if token matches with an user Id 
    if userId = GetUserIdWithToken(token, w); userId == -1 {
        return
    }
    myStruct.ListPeopleWhoLikesMe = GetListPeopleWhoLikeMe(userId)
    myStruct.ListPeopleILike = GetListPeopleILike(userId)
    myStruct.Status = "OK"
    RenderJSON(w, myStruct, http.StatusOK)
}


//METHOD POST : like an user
func Like(c web.C, w http.ResponseWriter, r *http.Request) {
    
    //i get user login in url
    UserToShow := c.URLParams["login"]
    strings.Replace(UserToShow, " ", "", -1)
    
    var userId int
    var userIdToLike int
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
    
    if AtLeastOneImage(userId) == false {
        RenderJSON(w, RenderStructError("Pictures", "Sorry you need to have publish at least one image to like someone.") , http.StatusNotFound)
        return
    }
    //get user_id by login (get in url)
    userIdToLike = GetIdByLogin(UserToShow, w)
     
    //if user enter in param url doesn't exist   
    if userIdToLike == -1 {
        return 
    }
    //check if current user have blocked user he want like
    if IfIHaveBlockedThisUser(userId, userIdToLike) == true {
        RenderJSON(w, RenderStructError("Block", "this user has blocked you."), http.StatusOK)
        return
    }
    //check if user he want like has blocked current user
    if IfIHaveBlockedThisUser(userIdToLike, userId) == true {
        RenderJSON(w, RenderStructError("Block", "You have blocked this user."), http.StatusOK)
        return
    }
    
    //if user enter in param url doesn't matches with current user || users cannot like themself
    if userIdToLike != userId {
        rows, err := gest.Db.Query("Select 1 from likes where liked_by_id_user=? AND id_user=?", userId, userIdToLike)
        LogFatalError(err)
        //check if user if the user does not already like person
        if rows.Next() {
            RenderJSON(w, RenderStructError("Likes", "Sorry you already like this people.") , http.StatusConflict)
            return
        } else { //add like in DB
            _, err := gest.Db.Exec("Insert into likes(liked_by_id_user, id_user, create_at) values (?, ?, ?)", userId, userIdToLike, TimeNowString())
            LogFatalError(err)
        } // if user try to like himself
    } else {
            RenderJSON(w, RenderStructError("Likes", "Sorry you can not like yourself"), http.StatusConflict)
            return
    }
    
    //create an empty structure for like, is it init by default in STATUS OK because all possibles errors have been checked
    ilike := ILike{}
    ilike.Status = "OK"
    
    //this function get two like id of two users for other, if one of two value like1 and like2 are equals 0, there is no match between two users
    like1, like2 := GetTwoIdLikes(userId, userIdToLike)
    
    //means that a match exist
    if (like1 != 0 && like2 != 0) {
        ilike.Match = true
        //function who create a match
        CreateMatch(userId, userIdToLike, like1, like2)
    } else {
        // else if like1 or like2 equals 0 (don't find in DB) there is no match between two users
        ilike.Match = false
    }
    //render JSON
    RenderJSON(w, ilike, http.StatusOK)
}


//METHOD DELETE : unlike an user
func DelLike(c web.C, w http.ResponseWriter, r *http.Request) {
    
    //i get user login in url
    UserToShow := c.URLParams["login"]
    strings.Replace(UserToShow, " ", "", -1)
    
    var userId int
    var userIdToLike int
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
    //get user_id by login (get in url)
    userIdToLike = GetIdByLogin(UserToShow, w)
    //if user enter in param url doesn't exist   
    if userIdToLike == -1 {
        return
    }
    if IsLikeAPeople(userId, userIdToLike) == false {
        RenderJSON(w, RenderStructError("Like","You cannot unlike this people because actually you don't like him."), http.StatusNotFound)
        return
    }
    _, err := gest.Db.Exec("Delete from likes where liked_by_id_user=? AND id_user=?", userId, userIdToLike)
    LogFatalError(err)
    RenderJSON(w, RenderStructOk(), http.StatusOK)
}


//return an array struct LIKEPeople, it's a list of people with whom the user has a match
func GetListMatchs(userId int) []LikePeople {
    
     var peopleILike []LikePeople
    
    rows, err := gest.Db.Query("Select user_id_1 from matched where user_id_2=?", userId)
    LogFatalError(err)
    
    for rows.Next() {
        var tmp LikePeople
        var id int
        //rows scan id_user of people i have liked
        rows.Scan(&id)
        tmp.Bio = GetBio(id)
        tmp.Id = id
        tmp.Login = GetLogin(id)
        tmp.Name = GetName(id)
        tmp.Orientation = GetOrientation(id)
        peopleILike = append(peopleILike, tmp)
    }
    rows, err = gest.Db.Query("Select user_id_2 from matched where user_id_1=?", userId)
    LogFatalError(err)
    for rows.Next() {
        var tmp LikePeople
        var id int
        //rows scan id_user of people i have liked
        rows.Scan(&id)
        tmp.Bio = GetBio(id)
        tmp.Id = id
        tmp.Login = GetLogin(id)
        tmp.Name = GetName(id)
        tmp.Orientation = GetOrientation(id)
        peopleILike = append(peopleILike, tmp)
    }
    return peopleILike
}

//this function create a match for two users // WARNING it doesn't do any check, you must do it (maybe a function exist in file all_checks.go)
func CreateMatch(userId, userId2, userLikeId, userLikeId2 int) {
	   _, err := gest.Db.Exec("Insert into matched(user_id_1, user_id_2, user_like_id_1, user_like_id_2) values (?, ?, ?, ?)", userId, userId2, userLikeId, userLikeId2)
       LogFatalError(err)
}

//method get return a list of people with whom the user has a match
func MyMatchs(c web.C, w http.ResponseWriter, r *http.Request) {
    
    var matchs ListMatchPeople
    //get token in url
    var token = r.FormValue("token")
    var userId int
    
    //init my empty struct who contains an array string for login user
    
    //check if token exist
    if CheckToken(token, w) == false {
        return
    }
    //check if token matches with an user Id 
    if userId = GetUserIdWithToken(token, w); userId == -1 {
        return
    }
    matchs.MatchsList = GetListMatchs(userId)
    matchs.Status = "OK"
    RenderJSON(w, matchs, http.StatusOK)
}