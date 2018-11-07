package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"time"
	"io/ioutil"
)

type ApiClient struct {
	BaseUrl string 
}

func NewApiClient() ApiClient {
	return ApiClient{BaseUrl: "https://fantasy.premierleague.com/drf"}
}

func (apiClient *ApiClient)  GetStaticGameData() (*GameData, error) {
	url := fmt.Sprintf("%s/bootstrap-static", apiClient.BaseUrl)
	responseBody, err := get(url)
	if err != nil {
		return nil, err
	}

	gameData := &GameData{}
	json.Unmarshal(responseBody, gameData)
	return gameData, nil
}

func (apiClient *ApiClient) GetLeagueStandings(leagueId int) ([]Team, error) {
	url := fmt.Sprintf("%s/leagues-classic-standings/%d?phase=1", apiClient.BaseUrl, leagueId)
	responseBody, err := get(url)
	if err != nil {
		return nil, err
	}

	league := &League{}
	json.Unmarshal(responseBody, league)
	return league.Table.Teams, nil
}

func get(url string) ([]byte, error) {
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {		
		return nil, err
	}

	// FF API doesn't return anything without setting a user-agent
	request.Header.Set("User-Agent", "Mozilla/5.0")
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	
	return responseBody, nil
}

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

type GameData struct {
	Players []struct {
		Id 		  int `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"second_name"`
	} `json:"elements"`
}
