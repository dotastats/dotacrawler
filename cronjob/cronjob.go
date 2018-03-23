package cronjob

import (
	"dotacrawler/config"
	"dotacrawler/entity"
	"dotamaster/utils/ulog"
	"errors"
	"fmt"
	"sync"
	"time"
)

type cronJob struct {
	Sleep         chan bool
	TickerCrawl   *time.Ticker
	TickerProcess *time.Ticker
	TickerCheck   *time.Ticker
	TickerBot     *time.Ticker
	Error         error
	Wg            *sync.WaitGroup
	Entity        entity.Entity
}

var (
	cron       *cronJob
	confVpGame config.VpGame
	GetMethod  = "GET"
)

func NewCronJob(intervalSecond int) *cronJob {
	if cron == nil {
		cron = &cronJob{
			TickerCrawl:   time.NewTicker(time.Duration(intervalSecond) * time.Second),
			TickerProcess: time.NewTicker(time.Duration(intervalSecond/2) * time.Second),
			TickerCheck:   time.NewTicker(time.Duration(intervalSecond/2) * time.Second),
			TickerBot:     time.NewTicker(time.Duration(intervalSecond) * time.Second),
			Wg:            &sync.WaitGroup{},
			Entity:        entity.NewEntity(),
		}
		confVpGame = config.Get().GetVpGame()
	}
	return cron
}

func Run() {
	if cron == nil {
		panic(errors.New("Init cronjob"))
	}
	cron.Wg.Add(1)
	go func() {
		for {
			if cron.Error != nil {
				ulog.Logger().LogError("Can't crawl data", ulog.Fields{
					"ERROR": cron.Error,
				})
				cron.Wg.Done()
				return
			}
			select {
			case <-cron.TickerCrawl.C:
				if time.Now().Hour() >= 22 || time.Now().Hour() <= 7 {
					continue
				}
				fmt.Println("Crawl")
				cron.crawlVpGame()
			}
		}
	}()

	cron.Wg.Wait()
}

func (r *cronJob) crawlVpGame() {
	err := r.Entity.VpGame.CrawlMatches()
	if err != nil {
		ulog.Logger().LogError("Can't crawl data from vpgame", ulog.Fields{
			"ERROR": err,
		})
	}
	time.Sleep(1 * time.Second)
}
