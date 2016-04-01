package models 

import (
    "3YFram/helpers"
    "golang.org/x/crypto/bcrypt"
    "github.com/pborman/uuid"
    "strings"
    "net/smtp"
    "log"
    "time"
    "fmt"
)

type User struct {
    Login string
    Email string
    Password []byte
    VerifiedToken string
    SessionToken string
    Id int
}

//Cette fonction ajoute les informations utilisateur a la structure user courante grace a l'ID utilisateur
func (user *User) GetUserById() {
    rows, err := gest.db.Query("Select login, email from account where id=$1", user.Id)
    LogFatalError(err)
    
    var login string
    var email string
    for rows.Next() {
        rows.Scan(&login, &email)
    }
    user.Email = email
    user.Login = login
}

//Cette fonction ajoute les informations utilisateur a la structure user courante grace au login
func (user *User) GetUserByLogin() int {
    
    if !user.IsUserLoginExist() {
        return -1  
    }
    rows, err := gest.db.Query("Select id, email from account where login=$1", user.Login)
    LogFatalError(err)
    var id int
    var email string
    for rows.Next() {
        rows.Scan(&id, &email)
    }
    user.Id = id
    user.Email = email
    
    return id
}
//Cette fonction ajoute les informations utilisateur a la structure user courante grace a l'email
func (user *User) GetUserByEmail() int {
    if !user.IsUserEmailExist() {
        return -1
    }
    rows, err := gest.db.Query("Select id, login from account where email=$1", user.Email) 
    LogFatalError(err)
    var id int 
    var login string
    for rows.Next() {
        rows.Scan(&id, &login)
    }
    user.Id = id
    user.Login = login
    
    return id
}
//Cette fonction ajoute les informations utilisateur a la structure user courante grace a l'identifiant (email ou login)
func (user *User) GetUserByLoginOrEmail() int {
    userIdByEmail := user.GetUserByEmail()
    if userIdByEmail != -1 {
        user.Id = userIdByEmail
        return userIdByEmail
    }
    userIdByLogin := user.GetUserByLogin()
    if userIdByLogin != -1 {
        user.Id = userIdByLogin
        return userIdByLogin
    }
    return -1
}

//Cette fonction ajoute a la structure courante User les informations de l'utilisateur via le token de verifications
func (user *User) GetUserByVerifiedToken() int {
    
    if user.IsVerifiedTokenExpired() == true {
        return -1
    }
    rows, err := gest.db.Query("select user_id from verifiedSession where token=$1", user.VerifiedToken)
    LogFatalError(err)
    var id int
    for rows.Next() {
        rows.Scan(&id)
    }
    user.Id = id
    user.GetUserById()
    return id
}


func (user *User) IsSessionTokenExist() bool {
    rows, err := gest.db.Query("Select 1 from session where token=$1", user.SessionToken)
    LogFatalError(err)
    if rows.Next(){
        return true
    }
    return false
}

//cette fonction verifie si le token de verification est expiré
func (user *User) IsSessionTokenExpired() bool {
    
    if user.IsSessionTokenExist() == false {
        return true
    }
    //format of date in DB
    layout := "2006-01-02 15:04:05"

    rows, err := gest.db.Query("Select expire_at from session where token=$1", user.SessionToken)
    LogFatalError(err)

    var expireDate string
    
    for rows.Next() {
        rows.Scan(&expireDate)
        LogFatalError(err)
        
        // now = current time 
        now := time.Now()     
        
        expireAt, err := time.Parse(layout, expireDate)
        LogFatalError(err)
        
        if expireAt.After(now) {
            return false
        }
    }
    return true
}

//cette fonction ajoute l'email, le login et l'id a la structure User courante grace au token de session
func (user *User) GetUserBySessionToken() int {
    if user.IsSessionTokenExpired() == true {
        return -1
    }
    rows, err := gest.db.Query("select user_id from session where token=$1", user.SessionToken)
    LogFatalError(err)
    var id int
    
    for rows.Next() {
        rows.Scan(&id)
    }
    user.Id = id
    user.GetUserById()
    return id
}

//Cette fonction verifie que l'utilisateur est bel et bien verifié
func (user *User) IsUserIsVerified() bool {
    rows, err := gest.db.Query("Select verified from account where id=$1 OR login=$2 OR email=$3", user.Id, user.Login, user.Email)
    LogFatalError(err)
    var verified bool 
    for rows.Next() {
        rows.Scan(&verified)
    }
    if verified == true {
        return true
    }
    return false
}
//cette fonction verifie si le token de verification est expiré
func (user *User) IsVerifiedTokenExpired() bool {
    
    //format of date in DB
    layout := "2006-01-02 15:04:05"

    rows, err := gest.db.Query("Select expire_at from verifiedSession where token=$1", user.VerifiedToken)
    LogFatalError(err)

    var expireDate string
    
    for rows.Next() {
        rows.Scan(&expireDate)
        LogFatalError(err)
        
        // now = current time 
        now := time.Now()     
        
        expireAt, err := time.Parse(layout, expireDate)
        LogFatalError(err)
        
        if expireAt.After(now) {
            return false
        }
    }
    return true
}

//cette fonction verifie si le token de verification existe bien (pour se protegere des petits malins de merde)
func (user *User) IsExistVerifiedToken() bool {
    
    rows, err := gest.db.Query("Select 1 from verifiedSession where token=$1", user.VerifiedToken)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}

//cette fonction verifie si l'account est verifié
func (user *User) AccountVerified() {
    _, err := gest.db.Exec("update account set verified=$1 where id=$2 OR email=$3 OR login=$4", true, user.Id, user.Email, user.Login)
    LogFatalError(err)
}

//cette fonction verifie si le login existe
func (user *User) IsUserLoginExist() bool {
    rows , err := gest.db.Query("Select 1 from account where login=$1", user.Login)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}

//Cette fonction verifie si l'email existe
func (user *User) IsUserEmailExist() bool {
    rows, err := gest.db.Query("Select 1 from account where email=$1", user.Email)
    LogFatalError(err)
    if rows.Next() {
        return true
    }
    return false
}

//Cette fonction verifie un mot de passe hashé et un mot de passé non hashé et verifie si ils sont les memes
func (user *User) IsPasswordsMatchs(password []byte, hashedPassword []byte) bool {
    err := bcrypt.CompareHashAndPassword(hashedPassword, password)
    if err != nil {
        return false
    } else {
        return true
    }
}

//cette fonction verifie si l'identifiant et le mot de passe ajouté dans la structure User match bien avec celui de la DB
func (user *User) IsLoginOrEmailMatchWithPassword() bool {
    var password string
    rows, err := gest.db.Query("Select password from account where login=$1 OR email=$1", user.Login)
    LogFatalError(err)
    if rows.Next() {
            rows.Scan(&password)
    }
    if user.IsPasswordsMatchs(user.Password, []byte(password)) {
        return true
    }
    return false
}


//Fonction d'envoie d'email parametrer la fonction SendComfirmationMail pour modifié le rendu du mail
func send(body, title, email string) {
	from := "fantasim.dev@gmail.com" // pour l'instant fonctionne avec le smtp de gmail
	pass := "bahuzzqkloftxaop" // il faut generer ce code avec google.

	msg := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: "+ title + "\n\n" + 
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{email}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
    fmt.Println("Email envoyé à ", email)
}


func (user *User) SendComfirmationMail() {
    title := "Comfirmation de votre compte " + Conf.NameApp
    content := "Pour valider votre compte veuillez cliquer sur ce lien : " + Conf.Url + "/comfirmation?verifiedtoken=" + user.VerifiedToken
    go send(content, title, user.Email)
}


//Cette fonction generere un token de verification, l'ajoute dans la variable VerifiedToken de la structure User et l'ajoute egalement dans la DB
func (user *User) CreateVerifiedToken() {
    user.VerifiedToken = uuid.New()
    user.GetUserByLogin()
    
    _, err := gest.db.Exec("Insert into verifiedSession(token, user_id) values ($1, $2)", user.VerifiedToken, user.Id)
    LogFatalError(err)
}


//Creation d'un utilisateur
func (user *User) CreateUser() []string {
    //initialisation d'un tableau de string pour les erreurs
    var errs []string
    
    //On verifie que l'email et le nom d'utilisateur ne soit pas déjà pris.
    if user.IsUserEmailExist() == true {
        helpers.AddStringInArray(&errs, "Un compte utilise déjà cette addresse email.")
    }
    if user.IsUserLoginExist() == true {
        helpers.AddStringInArray(&errs, "Ce nom d'identifiant est déjà utilisé.")
    }
    
    //Si il n'y a pas d'erreurs
    if len(errs) == 0 {
        //on creer le compte dans la DB
        _, err := gest.db.Exec("Insert into account (login, email, password) values ($1, $2, $3)", user.Login, user.Email, string(user.Password))
        LogFatalError(err)
        //Si notre fichier de configuration est en prod OU la variable CheckEmail est initialisé a true 
        if Conf.EmailCheck == true || strings.Compare(Conf.Status, "prod") == 0 {
            
            //On creer un token de verification 
            user.CreateVerifiedToken()
            
            //On envoie un lien par mail à l'utilisateur pour qu'il fasse validé son compte.
            user.SendComfirmationMail()
        } else {
            //Sinon si la verification de mail pour une raison ou une autre est desactivé, le compte est verifié dès sa création.
           _, err := gest.db.Exec("update account set verified=$1 where login=$2", true, user.Login)
           LogFatalError(err)
        }
    }
    return errs
}
