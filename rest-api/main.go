package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

//APIKEY variable where you have to paste your API key
const APIKEY = "faf7e5bb"

//APIURL for API
const APIURL = "http://www.omdbapi.com/?apikey=" + APIKEY

//DBNAME variable for database name
const DBNAME = "RestGoMovie"

//DBUSERNAME variable for database username
const DBUSERNAME = "root"

//DBPASSWORD variable for database password
const DBPASSWORD = ""

//DBHOST variable for database address
const DBHOST = "127.0.0.1"

//DBPORT variable for database port number
const DBPORT = "3306"

//DATABASEURL variable
const DATABASEURL = DBUSERNAME + ":" + DBPASSWORD + "@tcp(" + DBHOST + ":" + DBPORT + ")/" + DBNAME + "?charset=utf8&parseTime=True&loc=Local"

type movie struct {
	MovieID string `gorm:"type:varchar(12);column:movie_id;primary_key" json:"imdbID"`
	Title   string `json:"Title"`
	Type    string `gorm:"type:varchar(100)" json:"Type"`
	Year    string `gorm:"type:varchar(6)" json:"Year"`
	Image   string `json:"Poster"`
}

type detailMovie struct {
	MovieID  string `gorm:"type:varchar(12);column:movie_id;primary_key" json:"imdbID"`
	Title    string `json:"Title"`
	Year     string `gorm:"type:varchar(6)" json:"Year"`
	Released string `gorm:"type:varchar(100)" json:"Released"`
	Genre    string `gorm:"type:varchar(250)" json:"Genre"`
	Director string `gorm:"type:varchar(100)" json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Image    string `json:"Poster"`
}

type movieList struct {
	List []movie `json:"Search"`
}

func (list movieList) save() {
	db, err := gorm.Open("mysql", DATABASEURL)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	db.Debug().DropTableIfExists(&movie{})
	db.AutoMigrate(&movie{})
	for _, row := range list.List {
		log.Println("ROW ===>> ", &row)
		db.Debug().Create(&row)
	}
}

func (detail detailMovie) save() {
	db, err := gorm.Open("mysql", DATABASEURL)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	db.Debug().DropTableIfExists(&detailMovie{})
	db.AutoMigrate(&detailMovie{})
	db.Debug().Create(&detail)
}

func getSearchOmdb(w http.ResponseWriter, r *http.Request) {
	movieName := r.URL.Query().Get("title")
	if movieName == "" {
		log.Println("title not found")
		w.Write([]byte("title must be provided via QueryParams"))
		return
	}

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	//For receiving API call
	client := http.Client{}
	APIURLReq := APIURL + "&s=" + movieName + "&page=" + page
	movieRequest, err := http.NewRequest(http.MethodGet, APIURLReq, nil)
	if err != nil {
		log.Fatal(err)
	}

	movieResponse, err := client.Do(movieRequest)
	if err != nil {
		log.Fatal(err)
	}

	movieBody, err := ioutil.ReadAll(movieResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	list := movieList{}
	err = json.Unmarshal(movieBody, &list)

	if err != nil {
		log.Fatal(err)
	}
	list.save()

	//For sending API call
	w.Header().Set("Access-Control-Allow-Origin", "*") //This heading is necesary for cross origin
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)

}

func getSearchDetail(w http.ResponseWriter, r *http.Request) {
	imdbID := r.URL.Query().Get("id")
	if imdbID == "" {
		log.Println("ID for imdbID not found")
		w.Write([]byte("ID must be provided via QueryParams"))
	}

	// For receiving API Call
	client := http.Client{}
	APIURLReq := APIURL + "&i=" + imdbID
	movieRequest, err := http.NewRequest(http.MethodGet, APIURLReq, nil)
	if err != nil {
		log.Fatal(err)
	}

	movieResponse, err := client.Do(movieRequest)

	movieBody, err := ioutil.ReadAll(movieResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	detail := detailMovie{}
	err = json.Unmarshal(movieBody, &detail)

	if err != nil {
		log.Fatal(err)
	}
	detail.save()

	// For sending API CALL
	w.Header().Set("Access-Control-Allow-Origin", "*") //This heading is necesary for cross origin
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/search", getSearchOmdb)
	router.HandleFunc("/searchDetail", getSearchDetail)
	fmt.Println("Listening of port 8090")
	http.ListenAndServe(":8090", router)
}
