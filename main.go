package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
	"sort"
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
	var leagueId = 592906
	league := &League{}

	var staticData = getStaticFFData()
	getJson(fmt.Sprintf("%s/leagues-classic-standings/%d?phase=1", ffApiUri, leagueId), league)

	breakdown := &Breakdown{}
	for i, team := range league.Table.Teams {
		getJson(fmt.Sprintf("%s/entry/%d/event/11/picks", ffApiUri, team.Id), &league.Table.Teams[i].Squad)
		
		for j, squadPlayer := range league.Table.Teams[i].Squad.Players {
			for _, player := range staticData.Players {
				if squadPlayer.Id == player.Id {
					var name = fmt.Sprintf("%s %s", player.FirstName, player.LastName)
					squadPlayer.Name = name
					league.Table.Teams[i].Squad.Players[j].Name = name
					break
				}
			}

			playerExists := false
			for k, player := range breakdown.Players {
				if player.Id == squadPlayer.Id {
					playerExists = true
					breakdown.Players[k].TotalPick += 1
					if (team.Rank < 6) { breakdown.Players[k].Top5Pick += 1 }
					if (team.Rank < 11) { breakdown.Players[k].Top10Pick += 1 }
					break
				}	
			}

			if !playerExists {
				var playerBreakdown = PlayerBreakdown{
					Id: squadPlayer.Id,
					Name: squadPlayer.Name,
					TotalPick: 1,	
				}			

				if (team.Rank < 6) { playerBreakdown.Top5Pick += 1 }
				if (team.Rank < 11) { playerBreakdown.Top10Pick += 1 }

				breakdown.Players = append(breakdown.Players, playerBreakdown)
			}
		}
	}

	sort.Slice(breakdown.Players, func(i, j int) bool {
		return breakdown.Players[i].Top10Pick > breakdown.Players[j].Top10Pick
	})

	json.NewEncoder(w).Encode(breakdown)
}

func getStaticFFData() *StaticData {
	staticData := &StaticData{}
	err := getJson(fmt.Sprintf("%s/bootstrap-static", ffApiUri), staticData)
	if err != nil {
		panic(err.Error())
	}
	return staticData
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

var ffApiUri = "https://fantasy.premierleague.com/drf"

type League struct {
	Table Table `json:"standings"`
}

type Table struct {
	Teams []Team `json:"results"`
}

type Team struct {
	Id         int `json:"entry"`
	Rank       int `json:"rank"`
	PlayerName string `json:"player_name"`
	Squad 	   Squad
}

type Squad struct {
	Players []Player `json:"picks"`
}

type Player struct {
	Id 		 int `json:"element"`
	Position int `json:"position"`
	Name 	 string 
}

type StaticData struct {
	Players []struct {
		Id 		  int `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"second_name"`
	} `json:"elements"`
}

type Breakdown struct {
	Players []PlayerBreakdown `json:"players"`
}

type PlayerBreakdown struct {
	Id	 		int `json:"id"`
	Name 		string `json:"name"`
	Top5Pick  int `json:"top_5_pick"`
	Top10Pick int `json:"top_10_pick"`
	TotalPick int `json:"total_pick"`
} 