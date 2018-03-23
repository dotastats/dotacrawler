package config

import (
	"dotamaster/cmd"
	"dotamaster/config"
	"sync"
)

var (
	conf Config
	once sync.Once
)

type Config struct {
	App App
	config.Config
	VpGame VpGame
	Log    Log
}

type App struct {
	Host      string
	Port      int
	Debug     bool
	Whitelist []string
}

type Log struct {
	Prefix     string
	Dir        string
	LevelDebug bool
}

type VpGame struct {
	UrlCrawlMatchVpgame  string
	UrlCrawlSeriesVpgame string
	UrlCdnVpgame         string
	IntervalSecondCrawl  int
}

func load() {
	once.Do(func() {
		// load master
		mconf := config.Load()
		// load child
		if err := cmd.GetViper().Unmarshal(&conf); err != nil {
			panic(err)
		}
		conf.Config = mconf
	})
}

func Load() Config {
	load()
	return conf
}

func Get() Config {
	load()
	return conf
}

func GetVpGame() VpGame { return conf.GetVpGame() }
func (r Config) GetVpGame() VpGame {
	return r.VpGame
}
