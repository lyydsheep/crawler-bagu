package internal

type App struct {
	svc *CrawlerService
}

func InitializeApp() *App {
	// 每秒爬一题目
	svc := NewCrawlerService(NewDestinations(), NewCollector(), NewRL(1))
	return &App{svc}
}

func (app *App) Run() {
	app.svc.Run()
}
