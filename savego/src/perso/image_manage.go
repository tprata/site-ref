package main 

import (
        "net/http"
        "github.com/zenazn/goji/web"
        "path/filepath"
        "os"
        "io"
        "strconv"
        "time"
        "math/rand"
)


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456790")

func GenerateHash(n int) string {
rand.Seed(time.Now().UnixNano())
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

//METHOD POST : send an image
func AddImage(c web.C, w http.ResponseWriter, r *http.Request) {
    
    errr := Errors{}
    
    //check if profile value (string) correspond at a boolean
    if r.FormValue("profile") != "true" && r.FormValue("profile") != "false" { 
        RenderJSON(w, RenderStructError("Argument", "Error value profile (boolean) to determine if pictures is a profil picture or not") , http.StatusBadRequest)
        return
    }
    //get boolean with ParseBool
    var profile, err = strconv.ParseBool(r.FormValue("profile"))
    LogFatalError(err)
    
    var userId int
    //check if token is valid
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    hash := GenerateHash(40)
    var url = hash
    
    if GetNumbersImage(userId) > NUMBERS_IMAGE - 1 {
         RenderJSON(w, RenderStructError("Image's number", "Sorry ! You cannot have more than "+ strconv.Itoa(NUMBERS_IMAGE) +" images") , http.StatusBadRequest)
         return
        }
    //i dont understand this line, i have find this on google and... it works.
     r.ParseMultipartForm(32 << 20)
     
     //get image send in post method 
     file, handler, err := r.FormFile("image")
     
     //get extension of file
     ext := filepath.Ext(handler.Filename)
     //check if file is an image
     if (ext != ".jpg" && ext != ".png" && ext != "bmp") {
         RenderJSON(w, RenderStructError("Image", "Sorry ! only image is authorized") , http.StatusBadRequest)
         return
       }
       //get pwd for rewrite image in good 
       pwd := "/Images/"+url
       //check len total path of image
       if (len(pwd) > 255){
         RenderJSON(w, RenderStructError("Image", "Uknown error please contact us"), http.StatusBadRequest)
         return
       }
       
       
       //open a file with path of image
       f, err := os.OpenFile("./Images/"+url, os.O_WRONLY|os.O_CREATE, 0666)
       LogFatalError(err)
       //and copy content int get file in opened file
       _, err = io.Copy(f, file)
       LogFatalError(err)
       //at end of function close file
       defer f.Close()
       
       //add path of image in db
       if GetNumbersImage(userId) == 0 {
           profile = true
       }
       _, err = gest.Db.Exec("Insert into image (image, profile, id_user, create_at) values (?, ?, ?, ?)",pwd, profile, userId, TimeNowString())
       LogFatalError(err)

       errr.Status = "OK"
       RenderJSON(w, errr, http.StatusOK)
}


//Return list of user's images 
func GetListImage(userId int) []Image2 {
    var arrayImage []Image2
    rows, err := gest.Db.Query("Select id, image, profile from image where id_user=?", userId)
    LogFatalError(err)

    for rows.Next() {
        
        Images := Image2{}
        
        //put id, imagepath and profile var got in column in tempory structure Images2
        rows.Scan(&Images.Id, &Images.ImagePath, &Images.Profile)
        //And append this tempory structure Image2 in an array of structure Image2
        arrayImage = append(arrayImage, Images)
    }
    return arrayImage
}

//METHOD GET : render json list of image
func ListImage(c web.C, w http.ResponseWriter, r *http.Request){
    
    var userId int
    Imagelist := ImageList{}
    //check if token is valid
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    
    //at end put array of structure image2 in structure Arraylist
    Imagelist.Images = GetListImage(userId)
    Imagelist.Status = "OK"
    RenderJSON(w, Imagelist, http.StatusOK)
}

//METHOD DELETE : delete an image
func DelImage(c web.C, w http.ResponseWriter, r *http.Request){
    
    tmp := c.URLParams["id"]
    IdImage, _ := strconv.Atoi(tmp)
    var userId int
    //check if token is valid
    if CheckToken(r.FormValue("token"), w) == false {
        return
    }
    //Get user_id by token render json error if crash ?
    if userId = GetUserIdWithToken(r.FormValue("token"), w); userId == -1 {
        return
    }
    if GetImageUserIdByImage(IdImage) == -1 {
        RenderJSON(w, RenderStructError("Delete", "You cannot only delete your pictures"), http.StatusConflict)
        return
    }
    _, err := gest.Db.Exec("Delete from image where id=?", IdImage)
    LogFatalError(err)
    RenderJSON(w, RenderStructOk(), http.StatusOK)
}