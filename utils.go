package main

import (
	"fmt"
	"strings"
)

// Box Score
type BoxScoreResponse struct {
	BoxScore BoxScore `json:"boxscore"`
	Plays    []Plays  `json:"plays"`
}

type Plays struct {
	Desc      string `json:"text"`
	HomeScore int    `json:"homeScore"`
	AwayScore int    `json:"awayScore"`
	Period    Period `json:"period"`
	Clock     Clock  `json:"clock"`
}

type Period struct {
	Value int `json:"number"`
}

type Clock struct {
	Value string `json:"displayValue"`
}

type BoxScore struct {
	Players []Players `json:"players"`
}

type Players struct {
	TeamInfo   TeamInfo     `json:"team"`
	Statistics []Statistics `json:"statistics"`
}

type Statistics struct {
	Athletes []Athletes `json:"athletes"`
}

type TeamInfo struct {
	Name      string `json:"displayName"`
	ShortName string `json:"name"`
}

type Athletes struct {
	Active  bool     `json:"active"`
	Athlete Athlete  `json:"athlete"`
	Starter bool     `json:"starter"`
	DNP     bool     `json:"didNotPlay"`
	Stats   []string `json:"stats"`
}

type Athlete struct {
	Name string `json:"shortName"`
}

// Team Scoring
type ScoreboardResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	ID           string        `json:"id"`
	Date         string        `json:"date"`
	Name         string        `json:"name"`
	ShortName    string        `json:"shortName"`
	Competitions []Competition `json:"competitions"`
	Status       Status        `json:"status"`
}

type Competition struct {
	Competitors []Competitor `json:"competitors"`
}

type Competitor struct {
	Team  Team   `json:"team"`
	Score string `json:"score"`
}

type Team struct {
	ID           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
	DisplayName  string `json:"displayName"`
}

type Status struct {
	Type struct {
		ShortDetail string `json:"shortDetail"`
	} `json:"type"`
}

var ASCIINums = map[string]string{
	"1": `
 ██╗
███║
╚██║
 ██║
 ██║
 ╚═╝`,
	"2": `
██████╗ 
╚════██╗
 █████╔╝
██╔═══╝ 
███████╗
╚══════╝`,
	"3": `
██████╗ 
╚════██╗
 █████╔╝
 ╚═══██╗
██████╔╝
╚═════╝`,
	"4": `
██╗  ██╗
██║  ██║
███████║
╚════██║
     ██║
     ╚═╝`,
	"5": `
    
███████╗
██╔════╝
███████╗
╚════██║
███████║
╚══════╝`,
	"6": `
    
 ██████╗ 
██╔════╝ 
███████╗ 
██╔═══██╗
╚██████╔╝
 ╚═════╝`,
	"7": `
    
███████╗
╚════██║
    ██╔╝
   ██╔╝ 
   ██║  
   ╚═╝`,
	"8": `
 █████╗ 
██╔══██╗
╚█████╔╝
██╔══██╗
╚█████╔╝
 ╚════╝`,
	"9": `
 █████╗ 
██╔══██╗
╚██████║
 ╚═══██║
 █████╔╝
 ╚════╝`,
	"0": `
 
 ██████╗ 
██╔═████╗
██║██╔██║
████╔╝██║
╚██████╔╝
 ╚═════╝`,
}

var NbaTeams = map[string]string{
	"ATL": "Atlanta Hawks",
	"BOS": "Boston Celtics",
	"BKN": "Brooklyn Nets",
	"CHA": "Charlotte Hornets",
	"CHI": "Chicago Bulls",
	"CLE": "Cleveland Cavaliers",
	"DAL": "Dallas Mavericks",
	"DEN": "Denver Nuggets",
	"DET": "Detroit Pistons",
	"GS":  "Golden State Warriors",
	"HOU": "Houston Rockets",
	"IND": "Indiana Pacers",
	"LAC": "Los Angeles Clippers",
	"LAL": "Los Angeles Lakers",
	"MEM": "Memphis Grizzlies",
	"MIA": "Miami Heat",
	"MIL": "Milwaukee Bucks",
	"MIN": "Minnesota Timberwolves",
	"NO":  "New Orleans Pelicans",
	"NY":  "New York Knicks",
	"OKC": "Oklahoma City Thunder",
	"ORL": "Orlando Magic",
	"PHI": "Philadelphia 76ers",
	"PHX": "Phoenix Suns",
	"POR": "Portland Trail Blazers",
	"SAC": "Sacramento Kings",
	"SA":  "San Antonio Spurs",
	"TOR": "Toronto Raptors",
	"UTA": "Utah Jazz",
	"WS":  "Washington Wizards",
}

func Displaytext(hex string, msg string) string {
	var r, g, b int
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	hex_color := fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)

	colored_msg := fmt.Sprintf("%s%s\033[0m", hex_color, msg)
	return colored_msg
}

// func ConvertStringtoAscii(score string) string {
// 	var result string
// 	for _, chr := range score {
// 		asciichar := ASCIINums[string(chr)]
// 		result += asciichar
// 	}
// 	return result
// }

func ConvertStringtoAscii(score string) string {
	// First, determine how many lines each ASCII digit has
	// (assuming all digits have the same number of lines)
	var lineCount int
	if len(score) > 0 && len(ASCIINums[string(score[0])]) > 0 {
		lineCount = len(strings.Split(ASCIINums[string(score[0])], "\n"))
	} else {
		return ""
	}

	// Create a slice to hold each line of the final result
	lines := make([]string, lineCount)

	// For each digit in the score
	for _, chr := range score {
		asciiChar := ASCIINums[string(chr)]

		// Split the ASCII art into lines
		asciiLines := strings.Split(asciiChar, "\n")

		// Add each line of this digit to the corresponding result line
		for i, line := range asciiLines {
			if i < lineCount {
				lines[i] += line
			}
		}
	}

	// Join all lines with newlines to create the final result
	return strings.Join(lines, "\n")
}

func displayASCIINumber(number string) string {
	lines := make([]string, 5)

	for _, digit := range number {
		digitStr := string(digit)
		if asciiArt, ok := ASCIINums[digitStr]; ok {
			artLines := strings.Split(asciiArt, "\n")
			for i, line := range artLines {
				if i < len(lines) {
					lines[i] += line
				}
			}
		}
	}

	return strings.Join(lines, "\n")
}
