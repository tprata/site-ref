package main 


import (
    "net/http"
    "github.com/zenazn/goji/web"
    "strconv"
    "time"
)

//this function create a session in db
func CreateSessionInDb(userId int, token, ipAddr string) {
     
     t := time.Now()
     sessiondate := t.Format("2006-01-02 15:04:05")

     t2 := t.AddDate(0, 3, 0)
     sessionexpire := t2.Format("2006-01-02 15:04:05")


    _, err := gest.DbLog.Exec("Insert into session(user_id, token, ip_address, create_at, expire_at) values (?, ?, ?, ?, ?)", userId, token, ipAddr, sessiondate, sessionexpire)
    LogFatalError(err)
}

//this function return a struct Location to know the locate of userId
func GetLocalisation(userId int, r *http.Request) Location {
	location := Location{r.Header.Get("X-AppEngine-Country"),
		r.Header.Get("X-AppEngine-Region"),
		r.Header.Get("X-AppEngine-City"),
		[]string{r.Header.Get("X-AppEngine-CityLatLong")}, userId}
        return location
}

//this func return last location to know the locate of userId
func GetLastLocalisation(userId int) Localisation {
    
    rows, err := gest.DbLog.Query("Select country, region, city, latitude, longitude from localisation where id_user=? Order by id Desc limit 1", userId)
    LogFatalError(err)
    location := Localisation{}
    for rows.Next() {
        rows.Scan(&location.Country, &location.Region, &location.City, &location.Latitude, &location.Longitude)
    }
    return location
}

//return a list of people have viewed current user's profile
func GetListVisiteHistory(userId int) []History {
    
    var ArrayHistory []History
    rows, err := gest.DbLog.Query("Select id_user, create_at, id from visite where user_visited=? ORDER BY id DESC LIMIT 100", userId)
    LogFatalError(err)
    for rows.Next() {
        history := History{}
        rows.Scan(&history.UserId, &history.Date, &history.Id)
        ArrayHistory = append(ArrayHistory, history)
    }
    return ArrayHistory
}

//METHOD GET : render json list of people have viewed current user's profile
func GetVisiteHistory(c web.C, w http.ResponseWriter, r *http.Request) {
   
   listHistory := ListHistory{}
   var userId int
    //check if token is valid
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    listHistory.List = GetListVisiteHistory(userId)
    listHistory.Status = "OK"
    RenderJSON(w, listHistory, http.StatusOK)
}

//when user get profile of another user, this func his called to add a visit in history visit of user visited
func SetVisite(userVisited, userId int) {
 
    _, err := gest.DbLog.Exec("Insert into visite (user_visited, id_user, create_at) values (?, ?, ?)", userVisited, userId, TimeNowString())
    LogFatalError(err)
}

//set new localisation in db if is possible
func SetLocalisation(userId int, r *http.Request) {
    
    var insert bool
    
    insert = true
    location := GetLocalisation(userId, r)
    
    if len(location.LatLong) != 2 || location.LatLong[0] == "" || location.LatLong[1] == "" {
        return 
    }
    
    lat, err := strconv.ParseFloat(location.LatLong[0], 64)
    long, err := strconv.ParseFloat(location.LatLong[1], 64)
    LogFatalError(err)
    
    rows, err := gest.DbLog.Query("Select latitude, longitude from localisation where id_user=?", userId)
    LogFatalError(err)
    for rows.Next() {
        var latitude float64
        var longitude float64
        rows.Scan(&latitude, &longitude)
        if latitude == lat && longitude == long {
            insert = false
        } 
    }
    
    if insert == true && location.City != "" && location.Country != "" {
        _, err := gest.DbLog.Exec("Insert into localisation(country, region, city, latitude, longitude, id_user) values (?, ?, ?, ?, ?, ?)", location.Country, location.Region, location.City, lat, long, userId)
        LogFatalError(err)
    }
}

//last request
func UpdateLastRequest(userId int) {
    str := TimeNowString()
    _, err := gest.Db.Exec("UPDATE account SET last_request=? WHERE id=?", str, userId)
    LogFatalError(err)
}