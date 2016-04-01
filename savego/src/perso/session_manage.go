package main 

import (
    "net/http"
    "time"
    "net"
    "github.com/pborman/uuid"
    "github.com/zenazn/goji/web"
    "github.com/asaskevich/govalidator"
    "strings"
)



// this function is called when user sign in
func CreateSession(c web.C, w http.ResponseWriter, r *http.Request) {
    
    w.Header().Set("Content-Type", "application/json")
    
    //INIT TO EMPTY STRUCTURE : RENDER TO JSON -> LOGIN IF IT WORKS OR ERRORS
    login := Login{}
    all_errors := Errors{}

    //INIT AN EMPTY MAP[STRING]STRING
    all_errors.List_Errors = make(map[string]string)


    //remote ip address
    ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
    
    //WE GET LOGIN(IDENTIFIANT) AND PASSWORD || METHOD POST 
    var identifiant = NormalizeString(r.FormValue("id")) // can be an email or username, it is us to check.
    var password = r.FormValue("password")

    //IF LOGIN HAVE BAD FORMAT
    if IsAConformLogin(identifiant) == false && govalidator.IsEmail(identifiant) == false {
        all_errors.List_Errors["login"] = "login or email doesn't exist"
        all_errors.Status = "KO"
    } else { // ELSE IF LOGIN HAVE GOOD FORMAT
       
        //we try to get a raws with selecting login or email (value of identifiant) to check if login or email exist in database
        rows, err := gest.Db.Query("SELECT email FROM account WHERE email=?", identifiant)
        rows2, err2 := gest.Db.Query("SELECT email FROM account WHERE login=?", identifiant)
       
        //this function is a personal function to check if an erreur exist and display it, if is 
        LogFatalError(err) 
        LogFatalError(err2)

        //try to get value in db who are similar to which was POST
        var rowsString string
        for rows.Next() {
            err = rows.Scan(&rowsString)
        }
        if rowsString == "" {
            for rows2.Next() {
            err = rows2.Scan(&rowsString)
            }
        }
        //0. if identifiant exist in db (login or email)
        if rowsString != "" {
            
            //1. if var identifiant is an email
            var passwordDB = ""
            if govalidator.IsEmail(rowsString) == true {

                // get password who correspond to email and put it in var passwordDB
                rows, err := gest.Db.Query("SELECT password FROM account WHERE email=?", rowsString)
                LogFatalError(err)
                for rows.Next() {
                    err = rows.Scan(&passwordDB)
                }
                //2. if password which was post and which in db are similar
                if (CheckPassword([]byte(password), []byte(passwordDB))) {
                    
                    //get user id with login or email, we send request in function
                    userId := GetUserIdByEmailOrLogin("SELECT id FROM account WHERE email=?", rowsString)
                    
                    token, exist := SessionToken(userId)
                    
                    if !exist {
                        login.Token = uuid.New() //generate a random token
                        CreateSessionInDb(userId, login.Token, ipAddress)
                    } else {
                        login.Token = token
                    }
                    //create a new session line in DB login
                    login.Status = "OK"
                } else { //2. else if two password don't correspond
                    all_errors.List_Errors["password"] = "your password is wrong"
                    all_errors.Status = "KO"                    
                }
            } else { //1. else if var is a login
                 
                 // get password who correspond to login and put it in var passwordDB
                rows, err := gest.Db.Query("SELECT password FROM account WHERE login=?", rowsString)
                LogFatalError(err)
                for rows.Next() {
                    err = rows.Scan(&passwordDB)
                }
                
                //3. if two password corresponded 
                if (CheckPassword([]byte(password), []byte(passwordDB))) {
                    
                    userId := GetUserIdByEmailOrLogin("SELECT id FROM account WHERE login=?", rowsString)

                    token, exist := SessionToken(userId)

                    if !exist {
                        login.Token = uuid.New() //generate a random token                    
                        CreateSessionInDb(userId, login.Token, ipAddress)
                    } else {
                        login.Token = token
                    }
                    
                    //create a new session line in DB login               
                    login.Status = "OK"

                } else { //3. else if two password don't corresponded
                    all_errors.List_Errors["password"] = "your password is wrong"
                    all_errors.Status = "KO"
                }
            } 
        } else { //0. else if identifiant don't exist in db
            	    all_errors.List_Errors["login"] = "login or email doesn't exist"
                    all_errors.Status = "KO"
                }
            }
            
            //Finally we have check if we have a login var which is an email or a login.
            //- if it has a bad format for sql injection
            //- if it exist in db
            //- if password which was POST correspond to login's password in db
            
            //if status in struct error is KO
            if strings.Compare(all_errors.Status, "KO") == 0 {
                RenderJSON(w, all_errors, http.StatusOK)
                
             //else if status in struct error is not KO
            } else {
                userIdByToken := GetUserIdWithToken(login.Token, w)
                SetLocalisation(userIdByToken, r)
                RenderJSON(w, login, http.StatusOK)
            }            
}

//this function is called each user call an action that requires a token for security session
func CheckToken(token string, w http.ResponseWriter) bool {
    
    var getExpire string

    //get token in db
    rows, err := gest.DbLog.Query("SELECT expire_at from session where token=?", token)
    LogFatalError(err)
    //format of date in Db
    layout := "2006-01-02 15:04:05"
    //if token exist, it find it and return true
    for rows.Next() {
        err = rows.Scan(&getExpire)
        LogFatalError(err)
        now := time.Now()
        expireAt, err := time.Parse(layout, getExpire)
        LogFatalError(err)
        if (expireAt.After(now)) {
            return true
        }
    } // else if token not exist, that render JSON error and return false
    token_errors := TokenErrors{Token: "Token is not valid, maybe you're Session has expired please sign-in,", Status: "KO"}
    RenderJSON(w, token_errors, http.StatusNotAcceptable)
    return false
}


//return string is token if there, return bool = false if not there and true if there
func SessionToken(userID int) (string, bool) {
    
    //format of date in DB
    layout := "2006-01-02 15:04:05"
    
    //Get all token for an userID in Db
    rows, err := gest.DbLog.Query("SELECT token, expire_at FROM session WHERE user_id=?", userID)
    LogFatalError(err)
    
    var token string
    var expireDate string
    
    //for all token we have
    for rows.Next() {
        
        //get token and expireDate of each line get with request
        err = rows.Scan(&token, &expireDate)
        LogFatalError(err)
        
        // now = current time 
        now := time.Now()
        
        // expireAt = expire date token
        expireAt, err := time.Parse(layout, expireDate)
        LogFatalError(err)
        
        //If a token is not expired exists, return it with true
        if expireAt.After(now) {
            return token, true
        }
    }
    // return an empty string and false if any token is not expired exists
    return "", false
}


//Delete token in param if there in DB, return true if there, false not
func DeleteToken(token string) bool {
    
    //check if token in parameters exists in DB
    rows, err := gest.DbLog.Query("SELECT from session WHERE token=?", token)
    LogFatalError(err)
    if (rows.Next()) {
        
        //if there in db, we delete it and return true
        _, err := gest.DbLog.Exec("DELETE from session WHERE token=?", token)
        LogFatalError(err)
        return true
        
    } //if no token was found
    return false
}

//Render userId with token, if there are errors, func return -1 and json error 
func GetUserIdWithToken(token string, w http.ResponseWriter) int {
   
    var user_id int
    
    rows, err := gest.DbLog.Query("Select user_id from session WHERE token=?", token)
    LogFatalError(err)
    if (rows.Next()) {
        err = rows.Scan(&user_id)
        LogFatalError(err)
        UpdateLastRequest(user_id)
        return user_id
    }
    token_errors := TokenErrors{Token: "Token is not valid, maybe you're Session has expired please sign-in,", Status: "KO"}
    RenderJSON(w, token_errors, http.StatusNotAcceptable)
    return -1
}