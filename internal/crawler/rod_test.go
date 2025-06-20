package crawler

import (
	"testing"
)

func TestRod(t *testing.T) {
	r := NewRod()
	r.Consume("https://www.bagujing.com/problem-exercise/6?pid=3172")

	//browser := rod.New().MustConnect().NoDefaultDevice()
	//page := browser.MustPage("https://www.wikipedia.org/").MustWindowFullscreen()
	//
	//page.MustElement("#searchInput").MustInput("earth")
	//page.MustElement("#search-form > fieldset > button").MustClick()
	//
	//page.MustWaitStable().MustScreenshot("a.png")
	//time.Sleep(time.Hour)

}
