package main 

import (
        "golang.org/x/crypto/bcrypt" 
        "regexp"
        "time"
)


//check if an user have at least one image
func AtLeastOneImage(userId int) bool {
    rows, err := gest.Db.Query("Select 1 from image where id_user=?", userId)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}


//return true if user is online
func UserIsOnline(userId int) bool {
    
    rows, err := gest.Db.Query("select last_request from account where id=?", userId)
    LogFatalError(err)
    layout := "2006-01-02 15:04:05"
    
    for rows.Next() {
        var tmps string
        rows.Scan(&tmps)
        if tmps == "" {
            return false
        }
        now := time.Now()
        expireAt, err := time.Parse(layout, tmps)
        LogFatalError(err)
        during := time.Minute * 5
        expireAt.Add(during)
        if (expireAt.After(now)) {
            return true
        }
    }
    return false
}


//this function return id of like for like between two userId in function's param
func GetTwoIdLikes(userId, userIdIsLiked int) (int, int) {
    var like1 int
    var like2 int
    rows, err := gest.Db.Query("Select id from likes where liked_by_id_user=? AND id_user=?", userId, userIdIsLiked)
    LogFatalError(err)
    for rows.Next() {
        rows.Scan(&like1)
    }
    rows.Close()
    rows, err = gest.Db.Query("Select id from likes where liked_by_id_user=? AND id_user=?", userIdIsLiked, userId)
    LogFatalError(err)
    for rows.Next() {   
        rows.Scan(&like2)
    }
    return like1, like2
}


//this function check if userId like userIdIsLiked
func IsLikeAPeople(userId, userIdIsLiked int) bool {
    rows, err := gest.Db.Query("Select 1 from likes where liked_by_id_user=? AND id_user=?", userId, userIdIsLiked)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}

//this function check if login is already used
func IsLoginExist(login string) bool {
    
    rows, err := gest.Db.Query("Select 1 from account where login=?", login)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}

//this function check if email is already used
func IsEmailExist(email string) bool {
    rows, err := gest.Db.Query("Select 1 from account where email=?", email)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}



//This function compare two passwd with function of import bcrypt : CompareHashAndPassword
func CheckPassword(password []byte, hashedPassword []byte) bool {
    err := bcrypt.CompareHashAndPassword(hashedPassword, password)
    if err != nil {
        return false
    } else {
        return true
    }
}

//this function check if login has a good syntax
func IsAConformLogin(login string) bool {
    var validID = regexp.MustCompile(`^[a-z0-9_-]{4,16}$`)
    return validID.MatchString(login)
}

//this function check if name has a good syntax
func IsAConformName(name string) bool {
    var validID = regexp.MustCompile(`^[a-zA-Z0-9áàâäãåçéèêëíìîïñóòôöõúùûüýÿæœÁÀÂÄÃÅÇÉÈÊËÍÌÎÏÑÓÒÔÖÕÚÙÛÜÝŸÆŒ\s]{3,50}$`)
    return validID.MatchString(name)
}