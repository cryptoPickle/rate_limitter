package bucket

type Bucket map[string]int

func New() Bucket {
	bucket := make(Bucket)
	return bucket
}

func (b Bucket) Set(key string, value int) error {
	b[key] = value
	return nil
}

func (b Bucket) Get(key string) (int, error) {
	return b[key], nil
}

func (b Bucket) Has(key string) bool {
	_, ok := b[key]
	return ok
}

func (b Bucket) DecrementAll() error {
	for key := range b {
		b[key]--
	}
	return nil
}
