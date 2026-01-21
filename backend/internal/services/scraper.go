package services

type Scraper struct{}

func NewScraper() *Scraper {
	return &Scraper{}
}

func (s *Scraper) FetchMetadata(url string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
