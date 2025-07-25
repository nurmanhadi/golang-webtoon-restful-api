package comic

type ComicRepository interface {
	Save(comic *Comic) error
	FindById(id string) (*Comic, error)
	FindAll(page int, size int) ([]Comic, error)
	CountTotal() (int64, error)
}
