package keyword

type Service struct{}

func NewService() *Service {
    return &Service{}
}

func (s *Service) List() []Keyword {
    return allKeywords
}
