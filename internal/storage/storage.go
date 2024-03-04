package storage

type URL struct {
	UUID     string `json:"uuid"`
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

type IRepo interface {
	Create(short string, original string) error
	Get(short string) (*URL, error)
}
