package controllers

import (
       "3YFram/models"
        "3YFram/helpers"
        "golang.org/x/crypto/bcrypt"
        "strings"
        "regexp"
        "github.com/asaskevich/govalidator"
)


type User struct {
    hashedPassword []byte
    password string
    password2 string
    email string
    login string
    VerifiedToken string
    model models.User
}

//Hash le password envoyé en parametre et le met dans la variable hashedPassword dans la structure
func (user *User) newCryptPasswd(password []byte) {
    hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    user.hashedPassword = hashedPassword
}


//Cette fonction verifie les deux passwords entrée lors de l'inscription par exemple.
func (user *User) checkPasswords() []string {
    
    var errors []string
        
    if strings.Compare(user.password, user.password2) != 0 {
        helpers.AddStringInArray(&errors, "Les deux mots de passe ne sont pas identiques.")
    } else if len(user.password) < 6 {
        helpers.AddStringInArray(&errors, "Le mot de passe est choisit est trop court, minimum 6 caractères.")
    } else {
        user.newCryptPasswd([]byte(user.password))
    }
    return errors
}

//Cette fonction est appelé lorsque une personne tente d'activer son compte.
func (user *User) ComfirmAccount() models.ListError {
    
    var errr []string
    
    //on ajoute dans la structure user du mode le token de verification
    user.model.VerifiedToken = user.VerifiedToken
    // si le token n'existe pas
    if user.model.IsExistVerifiedToken() == false {
        helpers.AddStringInArray(&errr, "Le token n'existe pas, veuillez vous créer un compte.")
    } else {
        //on recupere les informations utilisateurs en fonction du token de verification // email, id, login
        //si le token a expiré.
        if user.model.GetUserByVerifiedToken() == -1 {
            helpers.AddStringInArray(&errr, "Le token de validation a expiré, veuillez en demander un nouveau")
        } else {
            //sinon on post pour authentifié l'utilisateur
            user.model.AccountVerified()
        }
    }
    return models.RenderErrorList(errr)
}

//Creer un nouveau compte utilisateur
func (user *User) CreateAccount(login, email, password, password2 string) models.ListError {
   
    //on ajoute les variables envoyé dans la structure User
    user.login = helpers.NormalizeString(login)
    user.email = helpers.NormalizeString(email)
    user.password = password
    user.password2 = password2
    
    //On appel la fonction propre au controller qui verifie la validité du password dans son format.
    err := user.checkPasswords()
    
    // on verifie la syntax du login et de l'email.
    if user.IsAConformLogin() == false {
        helpers.AddStringInArray(&err, "Votre login n'est pas au bon format. Il doit contenir entre 4 et 16 caractères et ne pas avoir de caractères spéciaux.")
    }
    if user.IsAConformEmail() == false {
        helpers.AddStringInArray(&err, "Votre email n'est pas au bon format.")
    }
    
    //Si aucune erreur, on ajoute les variables a la structure du model pour ensuite les poster dans la base de donné
    if len(err) == 0 {
        user.model.Email = user.email
        user.model.Login  = user.login
        user.model.Password = user.hashedPassword
        //fonction qui post les data dans la DB
        errModel := user.model.CreateUser()
        //on join les deux tableaux d'erreus possiblement obtenu dans le controlleur et dans le model pour n'en faire plus qu'un.
        err = helpers.ArrayJoin(err, errModel)
    }
    return models.RenderErrorList(err)
}

//this function check if login has a good syntax
func (user *User) IsAConformLogin() bool {
    var validID = regexp.MustCompile(`^[a-z0-9_-]{4,16}$`)
    return validID.MatchString(user.login)
}

//this function check if email has a good syntax
func (user *User) IsAConformEmail() bool {
    if govalidator.IsEmail(user.email) == false {
        return false
    }
    return true
}