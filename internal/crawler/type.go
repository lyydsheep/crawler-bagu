package crawler

type Crawler interface {
	Consume(string) error
}
