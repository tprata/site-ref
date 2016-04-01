package main 

import (
    
    "database/sql"
    "fmt"
)

const (
    DB_USER     = "root"
    DB_PASSWORD = "galipette"
    DB_NAME     = "ary"
)

const (
    DB_LOG_USER     = "root"
    DB_LOG_PASSWORD = "galipette"
    DB_LOG_NAME     = "logging"
)

//                                             ^
//this function init DB USER with value in const ( / \   )
//	                                           |

func InitDbUser() *sql.DB {
    /*
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        DB_USER, DB_PASSWORD, DB_NAME)
    */
    dbinfo := fmt.Sprintf("%s:%s@/%s",
        DB_USER, DB_PASSWORD, DB_NAME)
    db, err := sql.Open("mysql", dbinfo)
    LogFatalError(err)
    err2 := db.Ping()
    LogFatalError(err2)
    return db
}

//                                             ^
//this function init DB LOGGING with value in const ( / \   )
//

func InitDbLog() *sql.DB {
  /*
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        DB_LOG_USER, DB_LOG_PASSWORD, DB_LOG_NAME)
  */
    dbinfo := fmt.Sprintf("%s:%s@/%s",
        DB_LOG_USER, DB_LOG_PASSWORD, DB_LOG_NAME)  
    db, err := sql.Open("mysql", dbinfo)
    LogFatalError(err)
    err2 := db.Ping()
    LogFatalError(err2)
    return db
}


//This function create all table we need
func DbCreateUser(db *sql.DB) {
  _, err := db.Exec("CREATE TABLE IF NOT EXISTS account(id SERIAL, login varchar(17) NOT NULL, name varchar(50) NOT NULL, password varchar(60) NOT NULL, email varchar(255) NOT NULL, create_at varchar(20) NOT NULL, last_request varchar(20) DEFAULT NULL)",)
  LogFatalError(err)
   _, err = db.Exec("CREATE TABLE IF NOT EXISTS more_information(id SERIAL, id_user int NOT NULL, bio text, sex boolean DEFAULT TRUE, orientation varchar(50), birth_date varchar(10) NOT NULL DEFAULT '03-03-2016', visit_nb int DEFAULT 0, like_nb int DEFAULT 0, create_at varchar(20) NOT NULL)",)
  LogFatalError(err)
  _, err = db.Exec("CREATE TABLE IF NOT EXISTS tag(id SERIAL, id_user int NOT NULL, tag varchar(50), create_at varchar(20) NOT NULL)")
  LogFatalError(err)
  _, err = db.Exec("CREATE TABLE IF NOT EXISTS image(id SERIAL, id_user int NOT NULL, image varchar(255), profile boolean NOT NULL, create_at varchar(20) NOT NULL)")
  LogFatalError(err)
  _, err = db.Exec("CREATE TABLE IF NOT EXISTS likes(id SERIAL PRIMARY KEY, liked_by_id_user int NOT NULL, id_user int NOT NULL, create_at varchar(20) NOT NULL)")
 LogFatalError(err)
  _, err = db.Exec("CREATE TABLE IF NOT EXISTS matched(id SERIAL, user_id_1 int NOT NULL, user_id_2 int NOT NULL, user_like_id_1 int NOT NULL REFERENCES likes ON DELETE CASCADE, user_like_id_2 int NOT NULL REFERENCES likes ON DELETE CASCADE, create_at varchar(20) NOT NULL)")
  LogFatalError(err)
  _, err = db.Exec("CREATE TABLE IF NOT EXISTS block(id SERIAL, user_id int NOT NULL, blocked_by_id_user int NOT NULL, create_at varchar(20) NOT NULL)")
  LogFatalError(err)
}

//this function create all table we need
func DbCreateLogging(db *sql.DB) {
    _, err := db.Exec("CREATE TABLE IF NOT EXISTS session(id SERIAL, token varchar(60) NOT NULL, ip_address varchar(45), user_id int NOT NULL, create_at DATE NOT NULL, expire_at DATE NOT NULL)")
    LogFatalError(err)
    //_, err := db.Exec("create trigger `session_create_at` before insert on `session` for each row set new.`create_at` = to_char(now(), 'YYYY-MM-DD HH:MM:SS')")
    //LogFatalError(err)
    //_, err := db.Exec("create trigger `session_expire_at` before insert on `session` for each row set new.`expire_at` = to_char(now() + interval '3 month', 'YYYY-MM-DD HH:MM:SS');")
    //L1ogFatalError(err)
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS visite(id SERIAL, user_visited int NOT NULL, id_user int NOT NULL, create_at varchar(20) NOT NULL)")
    LogFatalError(err)
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS report(id SERIAL, cause text, user_reported int NOT NULL, id_user int NOT NULL, create_at varchar(20) NOT NULL)")
    LogFatalError(err)
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS localisation(id SERIAL, country varchar(50) DEFAULT '', region varchar(50) DEFAULT '', city varchar(50) DEFAULT '', latitude float DEFAULT 0, longitude float DEFAULT 0, id_user int NOT NULL)")
    LogFatalError(err)
}