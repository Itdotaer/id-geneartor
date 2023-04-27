package service

type GeneratorService interface {
	NextId(business string) (int64, error)
}
