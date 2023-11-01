package models

type Source struct {
	ID        SourceID
	SortIndex int
	Name      string
	Path      string
	Recursive bool
}
