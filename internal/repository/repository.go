package repository

type Billboard interface {
	GetBillboard() error
}

type PreliminaryResults interface {
}

type FinalResults interface {
}

type Repositories struct {
	Billboard          Billboard
	PreliminaryResults PreliminaryResults
	FinalResults       FinalResults
}

func NewRepositories() Repositories {
	return Repositories{}
}
