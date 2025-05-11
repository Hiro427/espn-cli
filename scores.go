package main

import (
	"encoding/json"
	// "flag"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	// "strconv"
	"strings"

	"time"
)

// TODO: Handle passing multiple teams to the '-team' flag
// TODO: Add support for other sports

func GetEventResponse(id string) Event {

	resp, err := http.Get(fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard/%s", id))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return event

}

// func (e Event) GetScoreTui(all bool) {
//
// 	team1 := e.Competitions[0].Competitors[0].Team.Abbreviation
// 	team1score := ConvertStringtoAscii(e.Competitions[0].Competitors[0].Score)
// 	team2 := e.Competitions[0].Competitors[1].Team.Abbreviation
// 	team2score := ConvertStringtoAscii(e.Competitions[0].Competitors[1].Score)
// 	GameTime := e.Status.Type.ShortDetail
//
// 	fmt.Printf("%s %s - %s - %s %s", team1, team1score, GameTime, team2, team2score)
// }

func (e Event) GetScoreTui(all bool) {
	team1 := e.Competitions[0].Competitors[0].Team.Abbreviation
	team1score := ConvertStringtoAscii(e.Competitions[0].Competitors[0].Score)
	team2 := e.Competitions[0].Competitors[1].Team.Abbreviation
	team2score := ConvertStringtoAscii(e.Competitions[0].Competitors[1].Score)
	gameTime := e.Status.Type.ShortDetail

	// Split the ASCII scores into lines
	team1ScoreLines := strings.Split(team1score, "\n")
	team2ScoreLines := strings.Split(team2score, "\n")

	// Get the number of lines in the ASCII art
	numLines := len(team1ScoreLines)
	if len(team2ScoreLines) > numLines {
		numLines = len(team2ScoreLines)
	}

	// Calculate the middle line index
	middleLine := numLines / 2

	// Print each line of the output
	for i := 0; i < numLines; i++ {
		if i == middleLine {
			// This is the middle line - include the team names and game time
			fmt.Print(team1 + " ")

			if i < len(team1ScoreLines) {
				fmt.Print(team1ScoreLines[i])
			}

			fmt.Print(" - " + gameTime + " - ")

			fmt.Print(team2 + " ")

			if i < len(team2ScoreLines) {
				fmt.Println(team2ScoreLines[i])
			} else {
				fmt.Println()
			}
		} else {
			// For non-middle lines, add spacing where the team1 name would be
			fmt.Print(strings.Repeat(" ", len(team1)+1))

			// Print team1 score line if available
			if i < len(team1ScoreLines) {
				fmt.Print(team1ScoreLines[i])
			} else {
				// Print spaces to maintain alignment
				fmt.Print(strings.Repeat(" ", len(team1ScoreLines[0])))
			}

			// Add middle spacing (where game time and separators would be)
			middleSpacing := len(" - " + gameTime + " - ")
			fmt.Print(strings.Repeat(" ", middleSpacing))

			// Add team2 name spacing
			fmt.Print(strings.Repeat(" ", len(team2)+1))

			// Print team2 score line if available
			if i < len(team2ScoreLines) {
				fmt.Println(team2ScoreLines[i])
			} else {
				fmt.Println()
			}
		}
	}
}

func (e Event) GetScore(all bool) {
	team1 := e.Competitions[0].Competitors[0].Team.Abbreviation
	team1score := e.Competitions[0].Competitors[0].Score
	team2 := e.Competitions[0].Competitors[1].Team.Abbreviation
	team2score := e.Competitions[0].Competitors[1].Score
	GameTime := e.Status.Type.ShortDetail

	if all {
		fmt.Printf("%-4s %-3s - %-4s %-3s %s\n", team1, team1score, team2, team2score, GameTime)
	} else {
		fmt.Printf("%s %s - %s %s %s\n", team1, team1score, team2, team2score, GameTime)
	}
}

func FetchEventIds(url string, team string) string {
	var e string
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not read response")
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Body")
		os.Exit(1)
	}

	var events ScoreboardResponse
	if err := json.Unmarshal(body, &events); err != nil {
		fmt.Println("Error Parsing Json")
	}

	for _, event := range events.Events {
		if strings.Contains(event.ShortName, strings.ToUpper(team)) {
			if len(event.Competitions) == 0 {
				continue
			}
			e = event.ID
		}
	}
	return e
}

func FetchAll(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not read response")
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Body")
		os.Exit(1)
	}

	var events ScoreboardResponse
	if err := json.Unmarshal(body, &events); err != nil {
		fmt.Println("Error Parsing Json")
	}

	for _, event := range events.Events {
		event_resp := GetEventResponse(event.ID)
		event_resp.GetScore(true)
		time.Sleep(200 * time.Millisecond)

	}
}

func GetBoxScore(id string, url string) BoxScoreResponse {
	box_url := fmt.Sprintf("%s%s", url, id)
	resp, err := http.Get(box_url)
	if err != nil {
		fmt.Println("Could not read response")
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Body")
		os.Exit(1)
	}

	var boxScore BoxScoreResponse
	if err := json.Unmarshal(body, &boxScore); err != nil {
		fmt.Println("Error Parsing Json")
	}
	return boxScore

}
func (b BoxScoreResponse) GetBoxScoreString(team string, id string, active bool) string {
	// Create a buffer to capture output
	var buf bytes.Buffer
	// Redirect output to the buffer temporarily
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the original function
	b.PrintBoxScore(team, id, active)

	// Restore stdout and get the captured output
	w.Close()
	os.Stdout = oldStdout
	io.Copy(&buf, r)

	return buf.String()
}

func (boxScore BoxScoreResponse) PrintBoxScore(team string, id string, active bool) {

	score := GetEventResponse(id)
	score.GetScore(false)

	fmt.Printf("\nPlay-by-Play: %s\n", boxScore.Plays[len(boxScore.Plays)-1].Desc)

	team1 := boxScore.BoxScore.Players[0]
	team2 := boxScore.BoxScore.Players[1]
	team1.ConstructPlayers(active)
	team2.ConstructPlayers(active)

}
func (p Players) ConstructPlayers(active bool) {
	var starters []string
	var bench []string
	var onlyactive []string

	fmt.Printf("\n%s\n", p.TeamInfo.Name)
	label := Displaytext("#cdd6f4", fmt.Sprintf(
		"%-23s %3s %3s %3s %3s %3s %3s %5s %5s %5s %s %2s %3s",
		"Name", "MIN", "PTS", "REB", "AST", "STL", "BLK", "FG", "3PT", "FT", "PF", "TO", "+/-"))
	fmt.Println(label)
	if !active {
		for _, a := range p.Statistics[0].Athletes {
			ath_stats := a.PrintAthleteScore(active)
			if a.Starter {
				starters = append(starters, ath_stats)
			} else {
				bench = append(bench, ath_stats)
			}
		}

		fmt.Printf("Starters\n")
		for _, str := range starters {
			if str == "" {
				continue
			}
			fmt.Println(str)
		}
		fmt.Printf("Bench\n")
		for _, bch := range bench {
			if bch == "" {
				continue
			}
			fmt.Println(bch)
		}
	} else {
		for _, a := range p.Statistics[0].Athletes {
			if a.Active {
				astat := a.PrintAthleteScore(active)
				onlyactive = append(onlyactive, astat)
			}
		}

		for _, str := range onlyactive {
			fmt.Println(str)
		}
	}

}

func (a Athletes) PrintAthleteScore(all bool) string {
	var player string
	if !a.DNP {
		name := a.Athlete.Name
		mins := a.Stats[0]
		fg := a.Stats[1]
		threes := a.Stats[2]
		ft := a.Stats[3]
		reb := a.Stats[6]
		ast := a.Stats[7]
		stl := a.Stats[8]
		blk := a.Stats[9]
		turnov := a.Stats[10]
		pf := a.Stats[11]
		plmi := a.Stats[12]
		pts := a.Stats[13]

		if !all {
			if a.Active {
				player = Displaytext("#cdd6f4", fmt.Sprintf(
					"%-23s %3s %3s %3s %3s %3s %3s %5s %5s %5s %s %2s %4s",
					name, mins, pts, reb, ast, stl, blk, fg, threes, ft, pf, turnov, plmi))
			} else {
				player = Displaytext("#585b70", fmt.Sprintf(
					"%-23s %3s %3s %3s %3s %3s %3s %5s %5s %5s %s %2s %4s",
					name, mins, pts, reb, ast, stl, blk, fg, threes, ft, pf, turnov, plmi))

			}

		} else {
			if a.Active {
				player = Displaytext("#cdd6f4", fmt.Sprintf(
					"%-23s %3s %3s %3s %3s %3s %3s %5s %5s %5s %s %2s %4s",
					name, mins, pts, reb, ast, stl, blk, fg, threes, ft, pf, turnov, plmi))
			}

		}

	}

	return player
}
