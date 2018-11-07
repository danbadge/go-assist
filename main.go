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
	router.HandleFunc("/", GetLeagueSquadBreakdowns).Methods("GET")

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

	gameDataService := GameDataService{GameData: gameData}

	teams, err := apiClient.GetLeagueStandings(leagueId)
	if err != nil {
		panic(err.Error())
	}

	var players []PlayerBreakdown
	for _, team := range teams {
		squad, err := apiClient.GetTeamSquad(team.Id, gameweek)
		if err != nil {
			panic(err.Error())
		}
		
		for _, squadPlayer := range squad.Players {
			squadPlayer.Name = gameDataService.GetPlayerName(squadPlayer.Id)

			player := findExistingPlayer(players, squadPlayer.Id)

			if player == nil {
				newPlayer := PlayerBreakdown{
					Id: squadPlayer.Id,
					Name: squadPlayer.Name,
					Picks: incrementPicks(Picks{}, team.Rank),
				}

				players = append(players, newPlayer)
			} else {
				player.Picks = incrementPicks(player.Picks, team.Rank)
			}
		}
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Picks.Total > players[j].Picks.Total
	})

	json.NewEncoder(w).Encode(Breakdown{Players: players})
}

type GameDataService struct {
	GameData *GameData
}

func (gameDataService *GameDataService) GetPlayerName(playerId int) string {
	for _, player := range gameDataService.GameData.Players {
		if playerId == player.Id {
			var name = fmt.Sprintf("%s %s", player.FirstName, player.LastName)
			return name
		}
	}
	return ""
}

func findExistingPlayer(players []PlayerBreakdown, id int) *PlayerBreakdown {
	for i, existingPlayer := range players {
		if existingPlayer.Id == id {
			return &players[i]
		}
	}
	return nil
}

func incrementPicks(picks Picks, teamRank int) Picks {
	return Picks{
		Top5: incrementIf(picks.Top5, teamRank < 6),
		Top10: incrementIf(picks.Top10, teamRank < 11),
		Total: picks.Total + 1,
	}
}

func incrementIf(current int, predicate bool) int {
	if (predicate) {
		return current + 1
	} 
	return current
}

type Breakdown struct {
	Players []PlayerBreakdown `json:"players"`
}

type PlayerBreakdown struct {
	Id	 		int `json:"id"`
	Name 		string `json:"name"`
	Picks		Picks `json:"picks"`
}

type Picks struct {
	Top5  int `json:"top_5_pick"`
	Top10 int `json:"top_10_pick"`
	Total int `json:"total_pick"`
}