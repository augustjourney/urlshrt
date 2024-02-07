package storage

type URL struct {
	Short    string
	Original string
}

type IRepo interface {
	Create(short string, original string)
	Get(short string) *URL
}
