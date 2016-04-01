package controllers

import (
    "3YFram/models"
    "3YFram/helpers"
    
)

type Session struct {
    Token string
    Status string
}


//cette fonction creer une nouvelle session.
func (session *Session) CreateSession(emailOrLogin, password, ipAddr string) models.ListError {
    
    //initialisation du tableau d'erreur
    var errr []string
    //initialisation de notre structure user 
    userModel := models.User{}
    
    //initialisation des variables dans la structure user avec celle qui ont été envoyé via la method post
    userModel.Login = helpers.NormalizeString(emailOrLogin)
    userModel.Email = userModel.Login
    userModel.Password = []byte(password)
    
    //if l'identifiant envoyé exist (cette identifiant peut etre un email ou login)
    if userModel.IsUserEmailExist() || userModel.IsUserLoginExist() {
        
        //si l'identifiant match avec le password associé
        if userModel.IsLoginOrEmailMatchWithPassword() {
            
            //on initialise une structure session du MODEL
            sessionModel := models.Session{}
            //on recuperer l'id de l'utilisateur et on le met dans session
            sessionModel.UserId = userModel.GetUserByLoginOrEmail()
           
           //si l'utilisateur est verifié 
            if userModel.IsUserIsVerified() == true {
                //creation d'un nouveau token qui sera renvoyé en response json et egalement ajouté a la base de donnée
                 session.Token = sessionModel.CreateSession(ipAddr)
            } else {
                helpers.AddStringInArray(&errr, "Votre compte n'est pas verifié, vous avez certainement reçu un email de comfirmation, verifiez dans vos spams ou demandez à en recevoir un nouveau.")
            }
        } else {
            helpers.AddStringInArray(&errr, "Le mot de passe ne correspond pas à celui du compte.")
        }
    } else {
        helpers.AddStringInArray(&errr, "Il n'existe aucun compte lié a cette identifiant")
    }
    //si la variable de la structure Session n'est pas vide (Si il est rentré dans la condition qui check que l'utilisateur est bien verifié)
    if session.Token != "" {
        session.Status = "OK"
    }
    return models.RenderErrorList(errr)
}
