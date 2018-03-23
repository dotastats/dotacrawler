package raw

type VPGameSchedule struct {
	LeftTeamID    string `json:"left_team_id"`
	RightTeamID   string `json:"right_team_id"`
	LeftTeamName  string `json:"left_team_name"`
	RightTeamName string `json:"right_team_name"`
}
