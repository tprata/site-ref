package main


import (
        "net/http"
        "github.com/zenazn/goji/web"
        "github.com/asaskevich/govalidator"
        "strings"
        "strconv"
)

var NUMBERS_IMAGE = 5



//get user's birth_date
func GetBirthDate(userId int) string {
    rows, err := gest.Db.Query("Select birth_date from more_information where id_user=?", userId)
    LogFatalError(err)
    var birthDate string
    for rows.Next() {
        rows.Scan(&birthDate)
    }
    return birthDate
}

//get user's bio
func GetBio(userId int) string {
    rows, err := gest.Db.Query("Select bio from more_information where id_user=?", userId)
    LogFatalError(err)
    var bio string
    for rows.Next() {
        rows.Scan(&bio)
    }
    return bio
}


//get user's email
func GetEmail(userId int) string {
    rows, err := gest.Db.Query("Select email from account where id=?", userId)
    LogFatalError(err)
    var email string
    for rows.Next() {
        rows.Scan(&email)
    }
    return email
}

//get user's name
func GetName(userId int) string {
    rows, err := gest.Db.Query("Select name from account where id=?", userId)
    LogFatalError(err)
    var name string
    for rows.Next() {
        rows.Scan(&name)
    }
    return name
}

//get user's login
func GetLogin(userId int) string {
    rows, err := gest.Db.Query("Select login from account where id=?", userId)
    LogFatalError(err)
    var login string
    for rows.Next() {
        rows.Scan(&login)
    }
    return login
}

//get user's orientation
func GetOrientation(userId int) string {
    rows, err := gest.Db.Query("Select orientation from more_information where id_user=?", userId)
    LogFatalError(err)
    var orientation string
    for rows.Next() {
        rows.Scan(&orientation)
    }
    return orientation
}

// get user's sexe
func GetSex(userId int) bool {
    rows, err := gest.Db.Query("Select sex from more_information where id_user=?", userId)
    LogFatalError(err)
    var sex bool
    for rows.Next() {
        rows.Scan(&sex)
    }
    return sex
}

//return a full profile (only use for permit to user to access his own informations)
func GetMyProfile(userId int) ProfileMe {
    
    profile := ProfileMe{}

    profile.Bio = GetBio(userId)
    profile.BirthDate = GetBirthDate(userId)
    profile.Email = GetEmail(userId)
    profile.Id = userId
    profile.Images = GetListImage(userId)
    profile.Matchs = GetListMatchs(userId)
    profile.Name = GetName(userId)
    profile.Orientation = GetOrientation(userId)
    profile.PeopleILike = GetListPeopleILike(userId)
    profile.PeopleLikesMe = GetListPeopleWhoLikeMe(userId)
    profile.Sexe = GetSex(userId)
    profile.Tags = GetListTag(userId)
    profile.Score = getScore(userId)
    profile.Online = UserIsOnline(userId)
    profile.Local = GetLastLocalisation(userId)
    profile.Status = "OK"

    return profile
}



func IndexShit(c web.C, w http.ResponseWriter, r *http.Request) {
    RenderJSON(w, RenderStructOk(), http.StatusOK)
}


//METHOD GET : {{url}}/user/profile/me
//render JSON to profile with login
func MyProfile(c web.C, w http.ResponseWriter, r *http.Request) {
    
    //get in param url login :  //http://url/user/profile/login 
    var userId int 
    //Get session token in param of url
    //get token passed in url param
    var token = r.FormValue("token")
    
    //init a Profile struct empty
    
    //Var who contains user_id
    //If session token in param isn't exist return with render error in json
    if CheckToken(token, w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(token, w); userId == -1 {
        return
    }
    
    //set localisation
    SetLocalisation(userId, r)
    
    RenderJSON(w, GetMyProfile(userId), http.StatusOK)
}

//METHOD GET : {{url}}/user/profile/:login
//render JSON to profile with login
func GetProfile(c web.C, w http.ResponseWriter, r *http.Request) {
    
    //get in param url login :  //http://url/user/profile/login 
    UserToShow := c.URLParams["login"]
    strings.Replace(UserToShow, " ", "", -1)
    
    var userId int 
    //Get session token in param of url
    //get token passed in url param
    var token = r.FormValue("token")
    
    //init a Profile struct empty
    userProfile := Profile{}
    
    //Var who contains user_id
    //If session token in param isn't exist return with render error in json
    if CheckToken(token, w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
   /* if userId = GetUserIdWithToken(token, w); userId == -1 {
        return
    }*/
    
    userId = GetIdByLogin(UserToShow, w)
    if (userId == -1) {
        return 
    }
    userIdByToken := GetUserIdWithToken(token, w)
    
    //set localisation
    SetLocalisation(userIdByToken, r)
    
    if IfIHaveBlockedThisUser(userId, userIdByToken) == true {
        RenderJSON(w, RenderStructError("Block", "this user has blocked you."), http.StatusOK)
        return
    }
    if IfIHaveBlockedThisUser(userIdByToken, userId) == true {
        RenderJSON(w, RenderStructError("Block", "You have blocked this user."), http.StatusOK)
        return
    }
    if userId == userIdByToken {
        RenderJSON(w, GetMyProfile(userId), http.StatusOK)
    } else {
        userProfile.Bio = GetBio(userId)
        userProfile.BirthDate = GetBirthDate(userId)
        userProfile.Id = userId
        userProfile.Images = GetListImage(userId)
        userProfile.Name = GetName(userId)
        userProfile.Orientation = GetOrientation(userId)
        userProfile.Sexe = GetSex(userId)
        userProfile.Tags = GetListTag(userId)
        userProfile.Score = getScore(userId)
        userProfile.Status = "OK"
        userProfile.Local = GetLastLocalisation(userId)
        userProfile.Online = UserIsOnline(userId)
        SetVisite(userId, userIdByToken)
        RenderJSON(w, userProfile, http.StatusOK)
    }
}

//this function is called when user sign-up
func CreateProfile(c web.C, w http.ResponseWriter, r *http.Request) {
    
    w.Header().Set("Content-Type", "application/json")
    //init an Inscri and Errors empty struct
    inscri := &Inscri{}
    errors := &Errors{}
    
    //this map [string]string contains all errors may
    some_errors := make(map[string]string)
 
 
    //get variable which were posted and put them in vars of Inscri's struct
    //using my function NormalizeString to delete all double space in a string
    inscri.Email = NormalizeString(r.FormValue("email"))
    inscri.Login = NormalizeString(r.FormValue("login"))
    inscri.Name = NormalizeString(r.FormValue("name"))
    inscri.Password = []byte(r.FormValue("password1"))
        
    //check len of password  ( p >= 6 && p <= 255)
    if len(r.FormValue("password1")) < 6 {
        some_errors["password"] = "password are too short, choose a password with minimum 6 characters."
    } else if len(r.FormValue("password1")) > 255 {
        some_errors["password"] = "password are too long"
    } else { //else if two password was posted are identics 
        if strings.Compare(r.FormValue("password1"), r.FormValue("password2")) != 0 {
            some_errors["password"] = "two passwords are not identical."
        } else { // else if two password are identic get a string of hash byte[] with my function newCryptPasswd
            inscri.Password = NewCryptPasswd(inscri.Password)
        }
    }
    
    //if login has a bad format
    if IsAConformLogin(inscri.Login) == false {
        if len(inscri.Login) < 4 {
            some_errors["login"] = "login you have choose is too short, choose a login with minimum 4 characters."
        } else if len (inscri.Login) > 16 {
            some_errors["login"] = "login you have choose is too long, chose a login with maximum 16 characters."
        } else {
            some_errors["login"] = "please use only only alphanumeric characters and _ or -"
        }
    } else { //else if login exist in DB
        if IsLoginExist(inscri.Login) == true {
            some_errors["login"] = "login you have choose already exist, please choose another login"
        }
    }
        
    //if Name has a bad format
    if IsAConformName(inscri.Name) == false {
        if len(inscri.Name) < 3 {
            some_errors["name"] = "your name is too short."
        } else if len(inscri.Name) > 50 {
            some_errors["name"] = "your name is too long, please reduce." 
        } else { 
            some_errors["name"] = "please don't use special characters"
        }
    }
    
    //if email has a bad format : govalidator 
     if govalidator.IsEmail(inscri.Email) == false {
         some_errors["email"] = "email chosen is not the right format."
     } else {
         // if email's format is ok, check if email already exist in db
         if IsEmailExist(inscri.Email) == true {
             some_errors["email"] = "email you have choose already exist, please choose another email"
            }
     }
     
     // if there are no errors 
     if len(some_errors) == 0 {
        _, err := gest.Db.Exec("Insert into account (login, name, password, email, create_at) values (?, ?, ?, ?, ?)", inscri.Login, inscri.Name, string(inscri.Password), inscri.Email, TimeNowString())
        LogFatalError(err)
        rows, err := gest.Db.Query("SELECT id FROM account ORDER BY id DESC LIMIT 1")
        var userID int
        for rows.Next() {
            err = rows.Scan(&userID)
        }
        _, err = gest.Db.Exec("Insert into more_information (bio, orientation, id_user, create_at) values (?, ?, ?, ?)", "", "", userID, TimeNowString())
        LogFatalError(err)
        
        //render struct errors to json with status OK
        RenderJSON(w, RenderStructOk(), http.StatusOK)
     } else {
         errors.List_Errors = some_errors
         errors.Status = "KO"

         //render struct errors to json with status KO
         RenderJSON(w, errors, http.StatusOK)
     }
}


//edit password
func EditPassword(c web.C, w http.ResponseWriter, r *http.Request) {
    
    var hash []byte
    
    //get parameters send in url
    var actualPassword = r.FormValue("lastpassword")
    var passwd = r.FormValue("password")
    var passwd2 = r.FormValue("password2")
    var userId int

    //check if token exist
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    
    //select actual user's password in DB
    rows, err := gest.Db.Query("SELECT password FROM account WHERE id=?", userId)
    LogFatalError(err)
    var passwordDB string
    for rows.Next() {
        
        //put actual user's password get in DB on var passwordDB
        err = rows.Scan(&passwordDB)
    }
    //compare actual password hash and clear password send by user, if return false, they don't correspond
    if CheckPassword([]byte(actualPassword), []byte(passwordDB)) == false {
        RenderJSON(w, RenderStructError("password", "password doesn't match with actual password."), http.StatusBadRequest)
        return
    }
    
    //if two password are not equals 
    if strings.Compare(passwd, passwd2) != 0 {
        RenderJSON(w, RenderStructError("password", "password doesn't match with repeat password."), http.StatusBadRequest)
        return
    }
    
    //if password are too short or too long
    if len(passwd) < 6 || len(passwd) > 50 {
        RenderJSON(w, RenderStructError("password", "password choose a password who have a len between 6 and 50."), http.StatusBadRequest)
        return
    }
    
    //if all is good, we hash password, and update actual user's password in DB
    hash = NewCryptPasswd([]byte(passwd))
    
    //update password
    _, err = gest.Db.Exec("UPDATE account SET password=? WHERE id=?", string(hash), userId)
    LogFatalError(err)
    
    //render Status OK
    RenderJSON(w, RenderStructOk(), http.StatusOK)
}

//edit login
func EditLogin(c web.C, w http.ResponseWriter, r *http.Request) {

    var userId int
    
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }

    var login = r.FormValue("login")
    
    //delete all space in login
    strings.Replace(login, " ", "", -1)

    //check if login's syntax is correct
    if IsAConformLogin(login) == false {
        RenderJSON(w, RenderStructError("login", "Please had a login with len between 4 and 16, and Alphanumeric caracters or _ and -."), http.StatusBadRequest)
        return
    }
    
    //check if login is already used
    if IsLoginExist(login) == true {
        RenderJSON(w, RenderStructError("login", "login is already used"), http.StatusBadRequest)
        return
    }
    
    //query update login
    _, err := gest.Db.Exec("UPDATE account SET login=? WHERE id=?", login, userId)
    LogFatalError(err)
    RenderJSON(w, RenderStructOk(), http.StatusOK)   
}


//edit email
func EditEmail(c web.C, w http.ResponseWriter, r *http.Request) {
       
    var userId int
    
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    
    var email = r.FormValue("email")
    
    //delete all space in var email
    strings.Replace(email, " ", "", -1)
    
    //check if email's syntax is correct
    if govalidator.IsEmail(email) == false {
        RenderJSON(w, RenderStructError("email", "your email haven't a good syntax"), http.StatusBadRequest)
        return
    }
    
    //check if email is already used
    if IsEmailExist(email) == true {
        RenderJSON(w, RenderStructError("email", "email is already used"), http.StatusBadRequest)
        return
    }
    
    //query update email
    _, err := gest.Db.Exec("UPDATE account SET email=? WHERE id=?", email, userId)
    LogFatalError(err)
    RenderJSON(w, RenderStructOk(), http.StatusOK)
}

//edit name
func EditName(c web.C, w http.ResponseWriter, r *http.Request) {
    
    var userId int
    
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    
    var name = NormalizeString(r.FormValue("name"))
    
    //check if name has a good syntax
    if IsAConformName(name) == false {
        RenderJSON(w, RenderStructError("name", "please don't use specials caracters for name"), http.StatusBadRequest)
        return
    }
    
    //query update name
    _, err := gest.Db.Exec("UPDATE account SET name=? WHERE id=?", name, userId)
    LogFatalError(err)
    RenderJSON(w, RenderStructOk(), http.StatusOK)
}

//edit bio
func EditBio(c web.C, w http.ResponseWriter, r *http.Request) {
    
    var userId int
    
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }

    var text = r.FormValue("text")
    
    //query update name
    _, err := gest.Db.Exec("UPDATE more_information SET bio=? WHERE id_user=?", text, userId)
    LogFatalError(err)
    RenderJSON(w, RenderStructOk(), http.StatusOK)
}

func EditBirthDate(c web.C, w http.ResponseWriter, r *http.Request) {

    var userId int
    
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    
    var day = r.FormValue("day")
    var month = r.FormValue("month")
    var year = r.FormValue("year")
    
    day_int, err1 := strconv.Atoi(day)
    month_int, err2 := strconv.Atoi(month)
    year_int, err3 := strconv.Atoi(year)
    
    if (err1 != nil || err2 != nil || err3 != nil) {
        RenderJSON(w, RenderStructError("Value", "Parse value error"), http.StatusOK)
        return;
    }
    if day_int < 1 || day_int > 31 || month_int < 1 || month_int > 12 || year_int < 1900 || year_int > 1998 {
        RenderJSON(w, RenderStructError("Value", "Parse value error"), http.StatusOK)
        return;      
    }
    var put = day + "-" + month + "-" + year
    _, err := gest.Db.Exec("UPDATE more_information SET birth_date=? WHERE id_user=?", put, userId)
    LogFatalError(err)
    
    RenderJSON(w, RenderStructOk(), http.StatusOK)
}

//edit orientation
func EditOrientation(c web.C, w http.ResponseWriter, r *http.Request) {
    
    var userId int
    
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    
    var orientation = r.FormValue("orientation")
    
    //check orientation len for hackers
    if len(orientation) > 50 {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    _, err := gest.Db.Exec("UPDATE more_information SET orientation=? WHERE id_user=?", orientation, userId)
    LogFatalError(err)
    RenderJSON(w, RenderStructOk(), http.StatusOK)
}

//return popularity score of user
func getScore(userId int) (int) {
    
    rows, err := gest.Db.Query("Select id from likes where id_user=?", userId)
    LogFatalError(err)
    var i int = 0
    for rows.Next() {
        i++
    }
    return i
}