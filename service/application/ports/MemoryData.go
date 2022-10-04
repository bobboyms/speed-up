package ports

type MemoryData interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Gets(key ...string) ([]string, error)
}
