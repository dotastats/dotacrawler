package entity

import (
	"dotacrawler/service"
	"dotamaster/models"
	"dotamaster/repo"
	"dotamaster/utils"
	"dotamaster/utils/uerror"
	"dotamaster/utils/ulog"
	"encoding/json"
)

type vpGame struct {
}
type VpGame interface {
	CrawlMatches() (err error)
}

func newVpGame() vpGame {
	return vpGame{}
}

func (r vpGame) CrawlMatches() (err error) {
	matches, err := r.getMatches()
	if err != nil {
		return err
	}

	// Save matches to db
	err = r.saveMatches(matches)
	if err != nil {
		return err
	}
	return
}

func (r vpGame) saveMatches(matches []models.VpMatch) (err error) {
	matchIds := make([]string, len(matches))
	for i, match := range matches {
		matchIds[i] = match.MatchID
	}

	matchIdsExists, err := repo.VpMatch.GetIdsExistsIn(matchIds)
	if err != nil {
		return uerror.StackTrace(err)
	}
	for _, match := range matches {
		if utils.IsExistedString(match.MatchID, matchIdsExists) {
			continue
		}
		err = repo.VpMatch.Create(&match)
		if err != nil {
			ulog.Logger().LogErrorObjectManual(err, "Can't create vpgame match", match)
			continue
		}
	}
	return nil
}

func (r vpGame) getMatches() (ret []models.VpMatch, err error) {
	url := confVpGame.UrlCrawlVpgame
	body, _, err := service.HttpReq.CrawlByURL("GET", url)
	if err != nil {
		err = uerror.StackTrace(err)
		return
	}
	err = json.NewDecoder(body).Decode(&ret)
	if err != nil {
		err = uerror.StackTrace(err)
		return
	}
	return
}
