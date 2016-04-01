package main

import (
    
    "net/http"
    "encoding/json"
    "html/template"
)

var static_view = "assets/view/"

//render an html web page
func Render_html_page(w http.ResponseWriter, url string) {
    url = static_view + url 
    t, err := template.ParseFiles(url) 
    if err != nil {
        panic (err)
    }
    t.Execute(w, nil)
}

//This function render json // stat = http.status
func RenderJSON(w http.ResponseWriter, v interface{}, stat int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stat)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic(err)
	}
}

//Render struct error 
func RenderStructError(key string, content string) Errors {
    errr := Errors{}
    errr.List_Errors = make(map[string]string)
    errr.List_Errors[key] = content
    errr.Status = "KO"
    return errr
}

//render struct success
func RenderStructOk() Errors {
    errr := Errors{}
    errr.Status = "OK"
    return errr
}