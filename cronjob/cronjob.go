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
	cron.Wg.Add(3)
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
				fmt.Println("Crawl live game")
				cron.crawlLiveVpGame()
			}
		}
	}()

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
				fmt.Println("Crawl open game")
				cron.crawlOpenVpGame()
			}
		}
	}()

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
				fmt.Println("Crawl close game")
				cron.crawlCloseVpGame()
			}
		}
	}()
	cron.Wg.Wait()
}

func (r *cronJob) crawlLiveVpGame() {
	err := r.Entity.VpGame.CrawlLiveMatches()
	if err != nil {
		ulog.Logger().LogError("Can't crawl live game from vpgame", ulog.Fields{
			"ERROR": err,
		})
	}
	time.Sleep(1 * time.Second)
}
func (r *cronJob) crawlOpenVpGame() {
	err := r.Entity.VpGame.CrawlOpenMatches()
	if err != nil {
		ulog.Logger().LogError("Can't crawl open game from vpgame", ulog.Fields{
			"ERROR": err,
		})
	}
	time.Sleep(1 * time.Second)
}
func (r *cronJob) crawlCloseVpGame() {
	err := r.Entity.VpGame.CrawlClosedMatches()
	if err != nil {
		ulog.Logger().LogError("Can't crawl close game from vpgame", ulog.Fields{
			"ERROR": err,
		})
	}
	time.Sleep(1 * time.Second)
}
