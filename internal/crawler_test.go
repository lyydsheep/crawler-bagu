package internal

import "testing"

func TestRun(t *testing.T) {
	svc := NewCrawlerService(NewDestinations(), NewCollector(), NewRL(1))
	svc.Run()
}
