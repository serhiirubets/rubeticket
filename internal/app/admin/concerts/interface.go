package concerts

// IConcertService определяет интерфейс для сервиса концертов.
type IConcertService interface {
	Create(request *CreateConcertRequest) (*ConcertResponse, error)
	Update(id uint, request *UpdateConcertRequest) (*ConcertResponse, error)
	Delete(id uint) error
	GetByID(id uint) (*ConcertResponse, error)
	List(page, pageSize int) (*ListConcertsResponse, error)
}
