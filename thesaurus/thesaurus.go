package thesaurus

type Thesaurus interface {
	Synonyms(string) ([]string, error)
}
