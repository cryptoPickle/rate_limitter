package types

type Nexter interface {
	Next()
	Get(string) (any, bool)
}

type IRateLimitter interface {
	Start(Nexter) error
}
