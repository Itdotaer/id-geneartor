package store

type Store interface {
	NextStep(business string) (int64, int64, error)
}
