package models

type Source struct {
	ID        SourceID
	Name      string
	Path      string
	Recursive bool
}
