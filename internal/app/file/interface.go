package file

type IFileRepository interface {
	Create(file *File) (*File, error)
	GetById(id string) (*File, error)
	CreateWithStorage(file *File) (*File, error)
}
