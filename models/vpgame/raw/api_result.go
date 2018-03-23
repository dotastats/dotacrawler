package raw

type VPgameAPIResult struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Body    []VPGameMatch `json:"body"`
}
