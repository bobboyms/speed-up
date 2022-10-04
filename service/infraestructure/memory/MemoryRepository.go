package memory

import (
	"fmt"
	"speed-up/service/application/logger"
	repositories "speed-up/service/application/repository"
	"speed-up/service/infraestructure/files"
	"sync"
	"time"
)

type IndexRepository struct {
	Chan  chan map[string]string
	Index sync.Map
}

func NewMemoryIndexRepository(logger logger.Logger) repositories.DataRepository {

	ch := make(chan map[string]string, 100)

	start := time.Now()
	index, count := files.LoadIndex()
	elapsed := time.Since(start)
	logger.Info("LOAD DATA", "Total..:", count, "Time ..: ", elapsed)

	i := IndexRepository{
		Chan:  ch,
		Index: index,
	}

	go files.DumpIndex(ch)

	return &i
}

func (i *IndexRepository) Get(key string) string {
	values, _ := i.Index.Load(key)
	return fmt.Sprintf("%v", values)
}

func (i *IndexRepository) Gets(keys ...string) []string {

	values := make([]string, 0)
	for _, key := range keys {
		values = append(values, i.Get(key))
	}

	return values
}

func (i *IndexRepository) Set(key string, value string) {
	i.Index.Store(key, value)
	i.Chan <- map[string]string{key: value}
}
