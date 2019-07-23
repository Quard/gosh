package storage

const InitialIdentifier = "aaaa"

type URLStorage interface {
	GetURL(identifier string) (string, error)
	AddURL(url string) (string, error)
}

type URL struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
