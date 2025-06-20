package internal

import "crawler-bagu/internal/crawler"

type App struct {
	svc *CrawlerService
}

func InitializeApp() *App {
	// 每秒爬一题目
	svc := NewCrawlerService(NewDestinations(), crawler.NewRod(), NewRL(1))
	return &App{svc}
}

func (app *App) Run() {
	app.svc.Run()
}
