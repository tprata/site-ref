package models

import (
       "database/sql"
       "fmt"
       "encoding/json"
       "os"
)

var gest GestionDB
var Conf Configuration


type GestionDB struct {
    db *sql.DB
}

type Configuration struct {
    EmailCheck    bool
    Status        string
    Url           string
    NameApp       string
    DbUser        string
    DbName        string
    DbPassword    string
}

const (
    DB_USER     = "fantasim"
    DB_PASSWORD = "aqw123"
    DB_NAME     = "whereami"
)

//cette fonction decode le fichier json en une structure (ici Configuration)
func ConfApp(pathFile string)  Configuration{
    file, _ := os.Open(pathFile)
    decoder := json.NewDecoder(file)
    cnf := Configuration{}
    err := decoder.Decode(&cnf)
    if err != nil {
        fmt.Println("error:", err)
    }
    return cnf
}

//This function create all table we need
func DbCreateUser(db *sql.DB) {
        _, err := db.Exec("CREATE TABLE IF NOT EXISTS account(id SERIAL, login varchar(17) NOT NULL, password varchar(60) NOT NULL, email varchar(255) NOT NULL, create_at varchar(20) NOT NULL DEFAULT to_char(now(), 'YYYY-MM-DD HH:MM:SS'), verified boolean DEFAULT false, last_request varchar(20) DEFAULT NULL)")
        LogFatalError(err)
        _, err = db.Exec("CREATE TABLE IF NOT EXISTS session(id SERIAL, token varchar(60) NOT NULL, ip_address varchar(45), user_id int NOT NULL, create_at varchar(20) NOT NULL DEFAULT to_char(now(), 'YYYY-MM-DD HH:MM:SS'), expire_at varchar(20) NOT NULL DEFAULT to_char(now() + interval '15 day', 'YYYY-MM-DD HH:MM:SS'))")
        LogFatalError(err)
        _, err = db.Exec("CREATE TABLE IF NOT EXISTS verifiedSession(id SERIAL, token varchar(60) NOT NULL, user_id int NOT NULL, create_at varchar(20) NOT NULL DEFAULT to_char(now(), 'YYYY-MM-DD HH:MM:SS'), expire_at varchar(20) NOT NULL DEFAULT to_char(now() + interval '7 day', 'YYYY-MM-DD HH:MM:SS'))")
        _, err = db.Exec("CREATE TABLE IF NOT EXISTS post(id SERIAL, user_id int NOT NULL, content text, title varchar(255) NOT NULL, important int NOT NULL, create_at varchar(20) NOT NULL DEFAULT to_char(now(), 'YYYY-MM-DD HH:MM:SS'))")
        LogFatalError(err)
}


//Initialise la database mais egalement la configuration du fichier app.json
func InitDB() GestionDB{
    Conf = ConfApp("./app.json")
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        Conf.DbUser, Conf.DbPassword, Conf.DbName)
    db, err := sql.Open("postgres", dbinfo)
    LogFatalError(err)
    err = db.Ping()
    LogFatalError(err)
    gest.db = db
    DbCreateUser(gest.db)
    return gest
}