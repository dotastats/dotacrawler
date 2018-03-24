package raw

import (
	"dotamaster/models"
	"regexp"
	"strconv"
	"strings"
	"time"
	"utilities/ulog"
)

type VPGameMatch struct {
	Id             string `json:"id"`
	Round          string `json:"round"`
	Category       string `json:"category"`
	ModeName       string `json:"mode_name"`
	ModeDesc       string `json:"name"`
	HandicapAmount string `json:"handicap"`
	HandicapTeam   string `json:"handicap_team"`
	SeriesID       string `json:"tournament_schedule_id"`
	GameTime       string `json:"game_time"`
	LeftTeam       string `json:"left_team"`
	RightTeam      string `json:"right_team"`
	// LeftTeamScore  string           `json:"left_team_score"`
	// RightTeamScore string           `json:"right_team_score"`
	Status     string           `json:"status_name"`
	Schedule   VPGameSchedule   `json:"schedule"`
	Odd        VPGameOdd        `json:"odd"`
	Team       VPGameTeam       `json:"team"`
	Tournament VPGameTournament `json:"tournament"`
}

func parseUnixTime(unixString string) (*time.Time, error) {
	i, err := strconv.ParseInt(unixString, 10, 64)
	if err != nil {
		return &time.Time{}, err
	}
	timeParsed := time.Unix(i, 0)
	return &timeParsed, nil
}

func processStatus(status string) string {
	status = strings.ToLower(status)
	if status == "live" {
		return "Live"
	} else if status == "settled" {
		return "Settled"
	} else if status == "canceled" {
		return "Canceled"
	}
	return "Upcoming"
}

func ratioProcess(ratio interface{}, matchStatus string) float64 {
	result, ok := ratio.(float64)
	if ok {
		return result
	}

	re := regexp.MustCompile("\\d+\\.\\d+")
	number := re.FindAllString(ratio.(string), -1)
	if number != nil {
		result, err := strconv.ParseFloat(number[0], 64)
		if err != nil {
			return 0
		}
		return result
	}
	return result
}
func winnerProcess(winner string) string {
	if len(winner) < 10 {
		return "TBD"
	}
	return winner[9:]
}

func (match *VPGameMatch) ConvertFromBase(baseMatch models.VpMatch) models.VpMatch {
	time, _ := parseUnixTime(match.GameTime)

	handicapTeam := match.LeftTeam
	if match.HandicapTeam == "right" {
		handicapTeam = match.RightTeam
	}

	winner := match.LeftTeam
	if match.Odd.Right.Victory == "win" {
		winner = match.RightTeam
	}

	matchId, _ := strconv.Atoi(match.Id)

	return models.VpMatch{
		TeamAID:        baseMatch.TeamAID,
		TeamBID:        baseMatch.TeamBID,
		TeamA:          baseMatch.TeamA,
		TeamB:          baseMatch.TeamB,
		Tournament:     baseMatch.Tournament,
		Game:           baseMatch.Game,
		BestOf:         baseMatch.BestOf,
		TournamentLogo: baseMatch.TournamentLogo,
		LogoA:          baseMatch.LogoA,
		LogoB:          baseMatch.LogoB,
		SeriesID:       baseMatch.SeriesID,

		MatchID:        matchId,
		URL:            "http://www.vpgame.com/match/" + match.Id,
		Time:           time,
		MatchName:      strings.TrimSpace(match.LeftTeam) + " vs " + strings.TrimSpace(match.RightTeam) + ", " + match.ModeName,
		ModeName:       match.ModeName,
		ModeDesc:       match.ModeDesc,
		HandicapTeam:   handicapTeam,
		HandicapAmount: match.HandicapAmount,
		Winner:         winner,
		RatioA:         ratioProcess(match.Odd.Left.Item, match.Status),
		RatioB:         ratioProcess(match.Odd.Right.Item, match.Status),
		Status:         processStatus(match.Status),
		TeamAShort:     match.Team.Left.NameShort,
		TeamBShort:     match.Team.Right.NameShort,
	}

}
func (match *VPGameMatch) CreateBaseMatch(logoURL string) models.VpMatch {
	bestOf := "BO1"
	if match.Round != "" {
		bestOf = match.Round
	}

	seriesId, err := strconv.Atoi(match.SeriesID)
	if err != nil {
		ulog.Logger().LogErrorObjectManual(err, "Can't convert string to integer", seriesId)
	}

	return models.VpMatch{
		TeamAID:        match.Team.Left.ID,
		TeamBID:        match.Team.Right.ID,
		TeamA:          strings.TrimSpace(match.Team.Left.Name),
		TeamB:          strings.TrimSpace(match.Team.Right.Name),
		Tournament:     strings.TrimSpace(match.Tournament.Name),
		Game:           match.Category,
		BestOf:         bestOf,
		TournamentLogo: logoURL + match.Tournament.Logo,
		LogoA:          logoURL + match.Team.Left.Logo,
		LogoB:          logoURL + match.Team.Right.Logo,
		SeriesID:       seriesId,
	}

}
