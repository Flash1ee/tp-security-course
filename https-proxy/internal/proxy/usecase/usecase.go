package usecase

type Usecase interface {
	Handle() error
	Close()
}
