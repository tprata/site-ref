package helpers

import (
    "strings"
    "unicode"
    "html/template"
    "encoding/json"
    "fmt"
    "net/http"
)

type Ok struct {
    Status string
}

//this function delete all double space in a string for replace them by one space
func NormalizeString(str string) string {
    str = strings.TrimSpace(str)
    var ret string    
    for i := 0; i < len(str); i++ {
        if (unicode.IsSpace(rune(str[i])) && unicode.IsSpace(rune(str[i + 1]))) {
        } else {
            var c = str[i]
            var s string = fmt.Sprintf("%c", c)
            ret += s
        }
    }
    return ret
}

//join deux tableaux pour n'en faire plus qu'un
func ArrayJoin(array1, array2 []string) []string{
    for i := 0; i < len(array2); i++ {
        array1 = append(array1, array2[i])
    }
    return array1
}

//this function add a string in array string with his pointer, for to avoid return 
func AddStringInArray(array *[]string, str string) {
    *array = append(*array, str)
}


//render an html web page
func RenderHTMLPage(w http.ResponseWriter, url string) {

    var static_view = "static/views/"

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

func RenderOk(w http.ResponseWriter) {
    RenderJSON(w, Ok{Status: "OK"}, http.StatusCreated)
}