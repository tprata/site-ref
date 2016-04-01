package main 

import (
    "net/http"
    "github.com/zenazn/goji/web"
    "strings"
    "regexp" 
    "strconv" 
)

 func AddTag(c web.C, w http.ResponseWriter, r *http.Request) {

    var userId int
    //check if token is valid
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    var tag = r.FormValue("tag")
    
    //delete all space in tag
    strings.Replace(tag, " ", "", -1)
    tag = strings.ToLower(tag)
    //init struct error
    errors := Errors{}
    errors.List_Errors = make(map[string]string)
   
    //check len of tag
    if len(tag) < 3 && len(tag) > 50 {
        errors.List_Errors["tag lenght"] = "hashtag needs to have a length between 3 and 50"
    }
    //check syntax of Hashtag
    var validID = regexp.MustCompile(`\S*#(?:\[[^\]]+\]|\S+)`)
    if validID.MatchString(tag) == false {
	   errors.List_Errors["tag syntax"] = "please use a good syntax for hashtag"
    }

    //We check if tag not exist ever in user's tags
    if len(errors.List_Errors) == 0 {
        rows, err := gest.Db.Query("Select 1 from tag where tag=? AND id_user=?", tag, userId)
        LogFatalError(err)
        if rows.Next() {
            errors.List_Errors["tag"] = "You have already this tag in yours hashtags"
        }
    }
    
    if len(errors.List_Errors) == 0 {
        _, err := gest.Db.Exec("Insert into tag(tag, id_user, create_at) values (?, ?, ?)", tag, userId, TimeNowString())
        LogFatalError(err)
        errors.Status = "OK"
        RenderJSON(w, errors, http.StatusOK)
    } else {
        errors.Status = "KO"
        RenderJSON(w, errors, http.StatusNotAcceptable)
    }
    
}

func GetListTag(userId int) []Tag {
    
    var tagArray []Tag
    //select all tag correspond to userId
    rows, err := gest.Db.Query("Select id, tag from tag where id_user=?", userId)
    LogFatalError(err)
    
     // for each column got
    for rows.Next() {
        
        tag := Tag{}
        
        //put id and content of actual tag in tempory var struct tag
        rows.Scan(&tag.Id, &tag.Tag)
        
        //and append tempory struct tag in an array of struct tag
        tagArray = append(tagArray, tag)
    }
    
    return tagArray
}

//METHOD GET {{url}}/user/tag/:login 

//render a json list of common tag between current_user and user (:login)
func ListTag(c web.C, w http.ResponseWriter, r *http.Request) {

    var userId int

    UserToShow := c.URLParams["login"]
    strings.Replace(UserToShow, " ", "", -1)
    listtag := TagList{} 
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
    userIdByLogin := GetIdByLogin(UserToShow, w)
    
    listtag.TagArray = CommonTag(userId, userIdByLogin)    
    listtag.Status = "OK"
    RenderJSON(w, listtag, http.StatusOK)
}

//if tag is existing in array tag, return true if true
func CheckTagInArrayTag(tag string, arr []Tag) bool {
    
    for i := 0; i < len(arr); i++ {
        if strings.Compare(tag, arr[i].Tag) == 0 {
            return true
        }
    }
    return false
}


//return an array tag of all common tag between 2 users
func CommonTag(userId, userIdByLogin int) []Tag {
    var ret []Tag
    
    ListTagCurrentUser := GetListTag(userId)
    ListTagParamUser := GetListTag(userIdByLogin)
    
    for i := 0; i < len(ListTagCurrentUser); i++ {
        for j := 0; j < len(ListTagParamUser); j++ {
            if strings.Compare(ListTagParamUser[j].Tag, ListTagCurrentUser[i].Tag) == 0 && CheckTagInArrayTag(ListTagParamUser[j].Tag, ret) == false {
                ret = append(ret, ListTagParamUser[j])
            }
        }
    }
    return ret
}


//METHOD DELETE : delete a tag
func DeleteTag(c web.C, w http.ResponseWriter, r *http.Request) {
    
    tmp := c.URLParams["id"]
    IdTag, _ := strconv.Atoi(tmp)
    var userId int
    //check if token is valid
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    //check if the tag matches a user's tag
    if GetUserIdByTagId(IdTag) == -1 {
        RenderJSON(w, RenderStructError("Delete", "You cannot only delete your tags"), http.StatusBadRequest)
        return
    }
    // if no problem, delete tag
    _, err := gest.Db.Exec("Delete from tag where id=?", IdTag)
    LogFatalError(err)
    
    //render status OK
    RenderJSON(w, RenderStructOk(), http.StatusOK)
}