package comic

type ComicRepository interface {
	Save(comic *Comic) error
	FindById(id string) (*Comic, error)
	FindAll(page int, size int) ([]Comic, error)
	CountTotal() (int64, error)
	Delete(id string) error
	Search(key string, page int, size int) ([]Comic, error)
	CountTotalByKeyword(key string) (int64, error)
}
