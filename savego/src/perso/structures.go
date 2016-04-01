package main 

import (
    "database/sql"
)

//this is API KEY ! IF ANYBODY KNOW IT, CHANGE IT 
const api_key string = "F378386FB33518DF8BE63D171E733"


//this structute contains information of user when he sign-up
type Inscri struct {
    Name string
    Login string
    Password []byte
    Email string
}

//this structure contains token and status with connexion
type Login struct {
    Token string
    Status string
}

//name is explicit
type Errors struct {
    List_Errors map[string]string
    Status string
}

//this contains *DB 
type Gestion struct {
    Db *sql.DB
    DbLog *sql.DB
}

type MoreInformations struct {
    Sex bool // true = M ; false = F
   	Bio string
    Orientation string
}

type TokenErrors struct {
    Token string
    Status string
}

//struct of an image
type Image2 struct {
    ImagePath string
    Profile bool
    Id int
}

//struct of image list
type ImageList struct {
    Images []Image2
    Status string
}

//struct of a tag 
type Tag struct {
    Id int
    Tag string
}

//struct of tag list
type TagList struct {
    TagArray []Tag
    Status string
}

//PROFILE

//profile current_user can see on other user
type Profile struct {
    Id int
    Name string
    Bio string
    BirthDate string
    Orientation string
    Sexe bool
    Tags []Tag
    Images []Image2
    Online bool
    Score int
    Local Localisation 
    Status string
}

//profile current_user can see for him.
type ProfileMe struct {
    Id int
    Email string
    Name string
    Bio string
    BirthDate string
    Orientation string
    Sexe bool
    Tags []Tag
    Images []Image2
    Matchs []LikePeople
    PeopleLikesMe []LikePeople
    PeopleILike []LikePeople
    Online bool
    Score int
    Local Localisation 
    Status string
}



// LIKES AND MATCH 

//struct of an user like
type LikePeople struct {
    Id int
    Login string
    Name string
    Bio string
    Orientation string
}
//struct of user list like
type ListLikePeople struct {
    LikesList []LikePeople 
    Status string
}

//struct of list of match
type ListMatchPeople struct {
    MatchsList []LikePeople 
    Status string
}

//struct of presence a match when user like a people
type ILike struct {
    Status string
    Match bool
}

//struct contains a list like, current_user to others users and other users to current_user
type AllLikes struct {
    ListPeopleWhoLikesMe []LikePeople
    ListPeopleILike []LikePeople
    Status string
}

//LOGGING

type History struct {
    Id int
    UserId int
    Date string
}
type ListHistory struct {
    List []History
    Status string
}

type Location struct {
	Country string
	Region  string
	City    string
	LatLong []string
    Id int
}

type Localisation struct {
	Country string
	Region  string
	City    string
    Latitude float64
    Longitude float64
}

// BLOCKED

// struct of blocked user
type Blocked struct {
    Id int
    Login string
    Name string
}

// struct list of blocked user
type ListBlocked struct {
    PeopleBlocked []Blocked
}



