package main

import (
    "golang.org/x/crypto/bcrypt" 
    "time"
    "net/http"
    "log"
    "strings"
    "unicode"
    "fmt"
    "encoding/json"
)

//return now time in string to format "2006-01-02 15:04:05"
func TimeNowString() string {
    t := time.Now()
    
    return t.Format("2006-01-02 15:04:05")
}

//return a hash passwd
func NewCryptPasswd(password []byte) []byte {
    hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    return hashedPassword 
}

//Create a cookie with a name, name of value, value and return it
func CreateCookie(nameCookie, nameCookieName string, value string) http.Cookie {
//    cookiee, _ := r.Cookie(nameCookie)
    expiration := time.Now().Add(365 * 24 * time.Hour) 
    cookie := http.Cookie{Name: nameCookieName, Value: value, Expires: expiration}
    //http.SetCookie(w, &cookie)
    return cookie 
}

//This function handles errors // if an error exist website crash, it's better in developement 
func LogFatalError(err error) {
    if err != nil {
    log.Fatal(err)
    }
}

//This function return id of an user with his login or email
//request = "SELECT id FROM account WHERE email=?" | email for email | login for login
func GetUserIdByEmailOrLogin(request, iden string) int {
    var id int
    rows, err := gest.Db.Query(request, iden)            
    LogFatalError(err)
   for rows.Next() {
       err = rows.Scan(&id)
    }
    return id
}

//this function delete all double space in a string for replace them by one space
func NormalizeString(str string) string {
    str = strings.TrimSpace(str)
    var ret string    
    for i := 0; i < len(str); i++ {
        if (unicode.IsSpace(rune(str[i])) && unicode.IsSpace(rune(str[i + 1]))) {
        } else {
            var c = str[i]
            var s string = fmt.Sprintf("%c", c)
            ret += s
        }
    }
    return ret
}

// this function enable to get id_user (int) by login (string)
func GetIdByLogin(login string, w http.ResponseWriter) int {
    
    errr := Errors{}
    errr.List_Errors = make(map[string]string)
    var userId int
    rows, err := gest.Db.Query("SELECT id FROM account WHERE login=?", login)
        for rows.Next() {
            err = rows.Scan(&userId)
            LogFatalError(err)
        }
        if userId == 0 {
            errr.List_Errors["login"] = "login user doesn't exist"
            errr.Status = "KO"
            RenderJSON(w, errr, http.StatusNotFound)
            return -1
        } else {
            return userId
        }
}

//this function enable to get login user (string) by id_user (int)
func GetLoginById(idLogin int) string {
    var login string
    rows, err := gest.Db.Query("Select login from account where id=?", idLogin)
    LogFatalError(err)
    for rows.Next() {
        rows.Scan(&login)  
    }
    return login
}

//return number of images of the user
func GetNumbersImage(userId int) int {

    var cnt int = 0

    rows, err := gest.Db.Query("Select id from image where id_user=?", userId)
    LogFatalError(err)
    for rows.Next() {
        cnt++
    }
    return cnt
}

//get user ID of an image ID
func GetImageUserIdByImage(IdImage int) int {
    rows, err := gest.Db.Query("Select id_user from image where id=?", IdImage)
    LogFatalError(err)
    for rows.Next(){
        var id int
        rows.Scan(&id)
        return id
    }
    return -1
}

//get user ID of a tag ID
func GetUserIdByTagId(IdTag int) int {
    rows, err := gest.Db.Query("Select id_user from tag where id=?", IdTag)
    LogFatalError(err)
    for rows.Next() {
        var id int
        rows.Scan(&id)
        return id
    }
    return -1
}


func Uint8toString(array []uint8) []string {
     var stringArray []string
     json.Unmarshal(array, &stringArray)
     return stringArray
}
