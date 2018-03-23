package entity

import (
	"dotacrawler/models/vpgame/raw"
	"dotacrawler/service"
	"dotamaster/models"
	"dotamaster/repo"
	"dotamaster/utils"
	"dotamaster/utils/uerror"
	"dotamaster/utils/ulog"
	"encoding/json"
	"fmt"
	"net/url"
)

type vpGame struct {
}

type VpGame interface {
	CrawlClosedMatches() (err error)
	CrawlOpenMatches() (err error)
	CrawlLiveMatches() (err error)
}

func newVpGame() vpGame {
	return vpGame{}
}

func (r vpGame) CrawlLiveMatches() error {
	vpParams := raw.VPGameAPIParams{
		Page:     "1",
		Status:   "start",
		Limit:    "200",
		Category: []string{"csgo", "dota"},
		Lang:     "en_US",
	}

	return r.crawlWithParam(vpParams)
}

func (r vpGame) CrawlOpenMatches() error {
	vpParams := raw.VPGameAPIParams{
		Page:     "1",
		Status:   "open",
		Limit:    "200",
		Category: []string{"csgo", "dota"},
		Lang:     "en_US",
	}

	return r.crawlWithParam(vpParams)
}

func (r vpGame) CrawlClosedMatches() error {
	vpParams := raw.VPGameAPIParams{
		Page:     "1",
		Status:   "close",
		Limit:    "200",
		Category: []string{"csgo", "dota"},
		Lang:     "en_US",
	}

	return r.crawlWithParam(vpParams)
}

func (r vpGame) crawlWithParam(vpParams raw.VPGameAPIParams) error {
	matches, err := r.getMatches(vpParams)
	if err != nil {
		return err
	}

	// Save matches to db
	err = r.saveMatches(matches)
	if err != nil {
		return err
	}
	return nil
}

func (r vpGame) saveMatches(matches []models.VpMatch) (err error) {
	matchIds := make([]int, len(matches))
	for i, match := range matches {
		matchIds[i] = match.MatchID
	}

	matchIdsExists, err := repo.VpMatch.GetIdsExistsIn(matchIds)
	if err != nil {
		return uerror.StackTrace(err)
	}
	for _, match := range matches {
		if utils.IsExistedInt(match.MatchID, matchIdsExists) {
			continue
		}
		err = repo.VpMatch.Create(&match)
		if err != nil {
			err = uerror.StackTrace(err)
			ulog.Logger().LogErrorObjectManual(err, "Can't create vpgame match", match)
			continue
		}
	}
	return nil
}

func (r vpGame) getMatches(vpParams raw.VPGameAPIParams) (ret []models.VpMatch, err error) {
	q := url.Values{}
	q.Set("page", vpParams.Page)
	q.Set("status", vpParams.Status)
	q.Set("limit", vpParams.Limit)
	for _, cate := range vpParams.Category {
		q.Add("category", cate)
	}
	q.Set("lang", vpParams.Lang)

	url := fmt.Sprintf("%s?%s", confVpGame.UrlCrawlMatchVpgame, q.Encode())

	body, _, err := service.HttpReq.CrawlByURL("GET", url)
	if err != nil {
		err = uerror.StackTrace(err)
		return
	}
	defer body.Close()

	var vpgameResult raw.VPgameAPIResult
	err = json.NewDecoder(body).Decode(&vpgameResult)
	if err != nil {
		err = uerror.StackTrace(err)
		return
	}

	// TODO: refactor
	for _, match := range vpgameResult.Body {
		matchBase := match.CreateBaseMatch(confVpGame.UrlCdnVpgame)

		var seriesResult raw.VPgameAPIResult
		seriesParam := raw.VPGameAPIParams{TID: match.SeriesID}
		seriesBody, _, err := service.HttpReq.CrawlByURL("GET", url)
		if err != nil {
			err = uerror.StackTrace(err)
			ulog.Logger().LogErrorObjectManual(err, "Can't get vpgame series", seriesParam)
			continue
		}
		defer seriesBody.Close()

		err = json.NewDecoder(seriesBody).Decode(&seriesResult)
		if err != nil {
			err = uerror.StackTrace(err)
			ulog.Logger().LogErrorObjectManual(err, "Can't decode vpgame series repsonse", seriesParam)
			continue
		}
		for _, match := range seriesResult.Body {
			matchFinal := match.ConvertFromBase(matchBase)
			ret = append(ret, matchFinal)
		}
	}

	return
}
