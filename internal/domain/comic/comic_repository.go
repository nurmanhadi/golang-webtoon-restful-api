package comic

type ComicRepository interface {
	Save(comic *Comic) error
	FindById(id string) (*Comic, error)
}
