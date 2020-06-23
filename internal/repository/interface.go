package repository

//go:generate mockery -name=DataProvider -dir=. -output=./mocks

// DataProvider is the contract for any repository implementations
type DataProvider interface{}
