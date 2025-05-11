# get team score 

curl -s "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/teams/9" | jq -r '.team.record.items[] | select (.type == "total") | .summary'

curl -s "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard" | jq -r 'select(.events[].date | contains("2025-04-08"))' | jq -r '.events[].id' | sort -u

curl -s "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard/401705712" | jq -r '.competitions[].competitors[1].team.abbreviation'





