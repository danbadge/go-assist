package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"sort"
)

func main() {
	var port = os.Getenv("PORT")

	var router = mux.NewRouter()
	log.Printf("Listening on port %s", port)
	router.HandleFunc("/league/squads", GetLeagueSquadBreakdowns).Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func GetLeagueSquadBreakdowns(w http.ResponseWriter, r *http.Request) {
	apiClient := NewApiClient()

	var leagueId = 592906
	var gameweek = 11

	gameData, err := apiClient.GetStaticGameData()
	if err != nil {
		panic(err.Error())
	}

	teams, err := apiClient.GetLeagueStandings(leagueId)
	if err != nil {
		panic(err.Error())
	}

	breakdown := &Breakdown{}
	for _, team := range teams {
		squad, err := apiClient.GetTeamSquad(team.Id, gameweek)
		if err != nil {
			panic(err.Error())
		}
		
		for _, squadPlayer := range squad.Players {
			for _, player := range gameData.Players {
				if squadPlayer.Id == player.Id {
					var name = fmt.Sprintf("%s %s", player.FirstName, player.LastName)
					squadPlayer.Name = name
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