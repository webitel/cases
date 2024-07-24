package model

type SearchOptions struct {
	Page   int
	Size   int
	Search string
	Sort   string
	Fields []string
}
