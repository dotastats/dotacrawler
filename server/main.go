package main

import (
	"dotacrawler/config"
	"dotacrawler/cronjob"
	"dotamaster/cmd"
	"dotamaster/infra"
)

func init() {
	cmd.Root().Use = "dotacrawler"
	cmd.Root().Short = "dotacrawler"
	cmd.Root().Long = "dotacrawler"

	cmd.SetRunFunc(run)
}

func main() {
	cmd.Execute()
}

func run() {
	setup()
	defer cleanup()

	conf := config.Get()
	cronjob.NewCronJob(conf.GetVpGame().IntervalSecondCrawl)

	cronjob.Run()
}

func setup() {
	infra.InitPostgreSQL()
}

func cleanup() {
	infra.ClosePostgreSql()
}
