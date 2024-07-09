package types

type Nexter interface {
	Next()
	Get(string) (any, bool)
}

type IRateLimitter interface {
	Start(Nexter) error
}

type Bucketer interface {
	Set(string, int) error
	Get(string) (int, error)
	Has(string) bool
	DecrementAll() error
}
