package main

import (
	"fmt"
	"github.com/spf13/cobra"
	// "os"
	"strings"
	"time"
)

func main() {
	url := "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard"
	summary := "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/summary?event="

	// Root command
	var rootCmd = &cobra.Command{
		Use:   "scores",
		Short: "A CLI tool using the unofficial ESPN API to fetch NBA scores",
	}

	var ScoreCmd = &cobra.Command{
		Use:   "nba",
		Short: "Get Box Score for a Specific Team",
		Run: func(cmd *cobra.Command, args []string) {
			BoxGame, _ := cmd.Flags().GetString("box")
			GameScore, _ := cmd.Flags().GetString("game")
			AllGameScore, _ := cmd.Flags().GetBool("all-games")
			Tui, _ := cmd.Flags().GetBool("tui")
			Active, _ := cmd.Flags().GetBool("active")
			Interval, _ := cmd.Flags().GetInt("interval")

			if BoxGame != "" && GameScore != "" {
				fmt.Println("Cannot use both --box and --game")
				return
			} else if BoxGame != "" {
				events := FetchEventIds(url, BoxGame)
				boxResponse := GetBoxScore(events, summary)
				if len(boxResponse.BoxScore.Players) == 0 {
					fmt.Println("No Box Score Data Available")
					return
					// os.Exit(1)
				}
				if events == "" {
					fmt.Printf("No %s Games Today\n",
						NbaTeams[strings.ToUpper(BoxGame)])
					return
				}
				if Tui {
					outputFn := func() string {
						// Fetch fresh data each time
						freshEvents := FetchEventIds(url, BoxGame)
						freshBoxResponse := GetBoxScore(freshEvents, summary)

						if freshEvents == "" {
							return fmt.Sprintf("No %s Games Today\n", NbaTeams[strings.ToUpper(BoxGame)])
						}

						// Get the formatted output
						result := freshBoxResponse.GetBoxScoreString(BoxGame, freshEvents, Active)
						return result
					}
					err := RunTUI(outputFn, time.Duration(Interval)*time.Second)
					if err != nil {
						fmt.Println("error using tui")
					}
				} else {
					boxResponse.PrintBoxScore(BoxGame, events, Active)
				}
				// PrintBoxScore(events, summary)
			} else if GameScore != "" {
				event_id := FetchEventIds(url, GameScore)
				event_resp := GetEventResponse(event_id)
				if len(event_resp.Competitions) == 0 {
					fmt.Printf("No %s games today\n", NbaTeams[strings.ToUpper(GameScore)])
					return
				}
				if Tui {
					event_resp.GetScoreTui(false)
				} else {
					event_resp.GetScore(false)
				}
			} else if AllGameScore {
				FetchAll(url)
			} else {
				fmt.Println("Invalid Entry")
				return
			}

		},
	}

	ScoreCmd.Flags().String("box", "", "Name of the team")
	ScoreCmd.Flags().String("game", "", "Name of the team")
	ScoreCmd.Flags().Bool("all-games", false, "Print all scores for this gameday")
	ScoreCmd.Flags().Bool("tui", false, "Run in TUI, only for box score ")
	ScoreCmd.Flags().Bool("active", false, "Only show active, less ouput to accomodate smaller terminal windows")
	ScoreCmd.Flags().Int("interval", 30, "Refresh interval for the TUI (seconds), default 30")

	rootCmd.AddCommand(ScoreCmd)

	// Execute the root command
	rootCmd.Execute()
}
