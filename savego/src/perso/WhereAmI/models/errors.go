package models 

import (
       "log"
)

type Error struct {
    Err string
}

type ListError struct {
    ListErr []Error
    Status string
}

func LogFatalError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

//transforme un tableau de string en tableau de structure d'erreur (pour rendre en JSON)
func RenderErrorList(err []string) ListError {
    list := ListError{}
    
    for i := 0; i < len(err); i++ {
        var tmp Error
        tmp.Err = err[i]
        list.ListErr = append(list.ListErr, tmp)
    }
    if len(list.ListErr) != 0 {
        list.Status = "KO"
    } else {
        list.Status = "OK"
    }
    return list
}