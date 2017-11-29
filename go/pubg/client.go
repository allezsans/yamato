package pubg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

// API holds configuration variables for accessing the API.
type API struct {
	APIKey  string
	BaseURL *url.URL
}

// JSONTime is RFC3339 format time.
type JSONTime struct {
	time.Time
}

// UnmarshalJSON hooks
func (j *JSONTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	s = s[1 : len(s)-1]

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.99", s)
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	j.Time = t.In(jst)
	return
}

// Format for output
func (j *JSONTime) Format() string {
	return j.Time.Format("2006-01-02 15:04")
}

// Player holds the information returned from the API.
type Player struct {
	PubgTrackerID int      `json:"pubgTrackerId"`
	AccountID     string   `json:"accountId"`
	Platform      int      `json:"platform"`
	NickName      string   `json:"nickName"`
	Avatar        string   `json:"avatar"`
	SteamName     string   `json:"steamName"`
	SteamID       string   `json:"steamId"`
	LastUpdated   JSONTime `json:"lastUpdated"`
	TimePlayed    int      `json:"timePlayed"`
	Stats         []Stats  `json:"stats"`
}

type Stats struct {
	Region string   `json:"region"`
	Season string   `json:"season"`
	Mode   string   `json:"mode"`
	Status []Status `json:"stats"`
}

type Status struct {
	Label        string  `json:"label"`
	Field        string  `json:"field"`
	Category     string  `json:"category"`
	ValueDec     float64 `json:"valueDec"`
	Value        string  `json:"value"`
	DisplayValue string  `json:"displayValue"`
}

// GetPlayerStatsFilterBy returns filtered by specific region
// 0:Region 1:Season 2:Mode
func (p *Player) GetPlayerStatsFilteredBy(f func(s Stats) bool) (Stats, error) {
	for _, v := range p.Stats {
		if f(v) {
			return v, nil
		}
	}
	return Stats{}, errors.New("Invalid filter")
}

func SelectLabel(s Stats, label string) string {
	for _, v := range s.Status {
		if v.Label == label {
			return v.DisplayValue
		}
	}
	return ""
}

// MatchHistory holds the information returned from the API.
type MatchHistory []struct {
	ID                   int      `json:"Id"`
	Updated              JSONTime `json:"Updated"`
	UpdatedJS            string   `json:"UpdatedJS"`
	Season               int      `json:"Season"`
	SeasonDisplay        string   `json:"SeasonDisplay"`
	Match                int      `json:"Match"`
	MatchDisplay         string   `json:"MatchDisplay"`
	Region               int      `json:"Region"`
	RegionDisplay        string   `json:"RegionDisplay"`
	Rounds               int      `json:"Rounds"`
	Wins                 int      `json:"Wins"`
	Kills                int      `json:"Kills"`
	Assists              int      `json:"Assists"`
	Top10                int      `json:"Top10"`
	Rating               float64  `json:"Rating"`
	RatingChange         float64  `json:"RatingChange"`
	RatingRank           int      `json:"RatingRank"`
	RatingRankChange     int      `json:"RatingRankChange"`
	Kd                   float64  `json:"Kd"`
	Damage               int      `json:"Damage"`
	TimeSurvived         float64  `json:"TimeSurvived"`
	WinRating            int      `json:"WinRating"`
	WinRank              int      `json:"WinRank"`
	WinRatingChange      int      `json:"WinRatingChange"`
	WinRatingRankChange  int      `json:"WinRatingRankChange"`
	KillRating           int      `json:"KillRating"`
	KillRank             int      `json:"KillRank"`
	KillRatingChange     int      `json:"KillRatingChange"`
	KillRatingRankChange int      `json:"KillRatingRankChange"`
	MoveDistance         float64  `json:"MoveDistance"`
}

func (m MatchHistory) Len() int {
	return len(m)
}

func (m MatchHistory) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MatchHistory) Less(i, j int) bool {
	return m[i].Updated.Time.Unix() > m[j].Updated.Time.Unix()
}

func GetMatchHistoryFilteredBy(f func(v MatchHistory) bool, m []MatchHistory) (MatchHistory, error) {
	for _, s := range m {
		if f(s) {
			return s, nil
		}
	}
	return MatchHistory{}, errors.New("Invalid filter")
}

// SteamInfo holds information returned from GetSteamInfo.
type SteamInfo struct {
	AccountID   string `json:"AccountId"`
	Nickname    string `json:"Nickname"`
	AvatarURL   string `json:"AvatarUrl"`
	SteamID     string `json:"SteamId"`
	SteamName   string `json:"SteamName"`
	State       string `json:"State"`
	InviteAllow string `json:"InviteAllow"`
}

// New creates a new API client.
func New(key string) (*API, error) {
	base, err := url.Parse("https://api.pubgtracker.com/v2/")
	if err != nil {
		return &API{}, err
	}

	return &API{
		APIKey:  key,
		BaseURL: base,
	}, nil
}

// NewRequest creates the GET request to access the API.
func (a *API) NewRequest(endpoint string) (*http.Request, error) {
	end, err := url.Parse(endpoint)
	if err != nil {
		return &http.Request{}, err
	}
	urlStr := a.BaseURL.ResolveReference(end)

	fmt.Printf(urlStr.String() + "\n")

	req, err := http.NewRequest("GET", urlStr.String(), nil)
	if err != nil {
		// Handle error
		return req, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Add("trn-api-key", a.APIKey)

	return req, nil
}

// Do sends out a request to the API and unmarshals the data.
func (a *API) Do(req *http.Request, i interface{}) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &i)
}

// GetPlayer returns a player's stats.
func (a *API) GetPlayer(param []string) (*Player, error) {
	endpoint := "profile/pc/"
	for i, value := range param {
		switch i {
		case 0:
			endpoint += value
		case 1:
			endpoint += fmt.Sprintf("?region=%s", value)
		case 2:
			endpoint += fmt.Sprintf("&season=%s", value)
		case 3:
			endpoint += fmt.Sprintf("&mode=%s", value)
		}
	}
	req, err := a.NewRequest(endpoint)

	if err != nil {
		return &Player{}, err
	}

	var player Player
	err = a.Do(req, &player)

	return &player, err
}

// GetMatchHistory reutnrns specific player match history.
func (a *API) GetMatchHistory(player *Player) (*MatchHistory, error) {
	endpoint := "matches/pc/" + player.AccountID
	fmt.Printf(player.AccountID + "\n")

	req, err := a.NewRequest(endpoint)

	if err != nil {
		return &MatchHistory{}, err
	}

	var matchHistory MatchHistory
	err = a.Do(req, &matchHistory)

	sort.Sort(matchHistory)

	return &matchHistory, err
}

// GetSteamInfo retrieves a player's steam information.
func (a *API) GetSteamInfo(sid string) (*SteamInfo, error) {
	endpoint := "search?steamId=" + sid
	req, err := a.NewRequest(endpoint)

	if err != nil {
		return &SteamInfo{}, err
	}

	var sinfo SteamInfo
	err = a.Do(req, &sinfo)

	return &sinfo, err
}
