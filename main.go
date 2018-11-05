package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

// our main function
func main() {
	var port = os.Getenv("PORT")

	var router = mux.NewRouter()
	log.Printf("Listening on port %s", port)
	router.HandleFunc("/league/squads", GetLeagueSquadBreakdowns).Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func GetLeagueSquadBreakdowns(w http.ResponseWriter, r *http.Request) {
	var ffApiUri = "https://fantasy.premierleague.com/drf"
	var leagueId = 592906
	var url = fmt.Sprintf("%s/leagues-classic-standings/%d?phase=1", ffApiUri, leagueId)
	league := &League{}
	getJson(url, league)
	json.NewEncoder(w).Encode(league)
}

func getJson(url string, target interface{}) error {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("User-Agent", "Mozilla/5.0")
	response, err := netClient.Do(request)
	if err != nil {
		return err
	}

    defer response.Body.Close()

    return json.NewDecoder(response.Body).Decode(target)
}

type League struct {
	Table Table `json:"standings"`
}

type Table struct {
	Team []Team `json:"results"`
}

type Team struct {
	Id         int `json:"entry"`
	Rank       int `json:"rank"`
	PlayerName string `json:"player_name"`
}
