package genre

type GenreRepository interface {
	Save(genre *Genre) error
	Remove(id int) error
	FindAll() ([]Genre, error)
	Count(id int) (int64, error)
}
