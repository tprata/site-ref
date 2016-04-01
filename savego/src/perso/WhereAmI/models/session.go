package models 

import (
    "github.com/pborman/uuid"
    "time"
    "strings"
    
)

type Session struct {
    UserId int
    Token string
    Status string
}


//Cette fonction verifie que le token n'est pas expiré
//return string is token if there, return bool = false if not there and true if there
func (session *Session) GetSessionToken() {
    
    //format de la date dans la Db
    layout := "2006-01-02 15:04:05"
    
    //Recuperation du token et de sa date d'expiration appartenant a l'Id de l'utilisateur contenu dans userId de session
    rows, err := gest.db.Query("SELECT token, expire_at FROM session WHERE user_id=$1", session.UserId)
    LogFatalError(err)
    
    var expireDate string
    
    //pour le token recupéré
    for rows.Next() {
        
        //scan du retour de la requete
        err = rows.Scan(&session.Token, &expireDate)
        LogFatalError(err)
        
        // now = heure actuel
        now := time.Now()
        
        // on formate la date recuperé au bon format
        expireAt, err := time.Parse(layout, expireDate)
        LogFatalError(err)
        
        //si la date d'expiration est passé alors on met le token a vide pour qu'il detecte qu'il faudra en generer un nouveau
        if !expireAt.After(now) {
            session.Token = ""
        }
    }
}

func CreateSessionInDb(userId int, token, ipAddr string) {
    _, err := gest.db.Exec("Insert into session(user_id, token, ip_address) values ($1, $2, $3)", userId, token, ipAddr)
    LogFatalError(err)
}

//Creer une nouvelle session si besoin est, sinon retourne le token
func (session *Session) CreateSession(ipAddr string) string {
    //Si la date d'expiration de l'ancien token est valide ou bien que c'est une premiere connexion, on genere un token.
    if session.GetSessionToken(); strings.Compare(session.Token, "") == 0 {
        session.Token = uuid.New()
        CreateSessionInDb(session.UserId, session.Token, ipAddr)
    }
    return session.Token
}