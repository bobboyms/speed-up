package memory

import (
	"speed-up/service/application/logger"
	repositories "speed-up/service/application/repository"
	"speed-up/service/infraestructure/files"
	"time"
)

type IndexRepository struct {
	Chan  chan map[string]string
	Index map[string]string
}

func NewMemoryIndexRepository(logger logger.Logger) repositories.DataRepository {

	ch := make(chan map[string]string, 100)

	start := time.Now()
	index := files.LoadIndex()
	elapsed := time.Since(start)
	logger.Info("LOAD DATA", "Total..:", len(index), "Time ..: ", elapsed)

	i := IndexRepository{
		Chan:  ch,
		Index: index,
	}

	go files.DumpIndex(ch)

	return &i
}

func (i *IndexRepository) Get(key string) string {
	values := i.Index[key]
	return values
}

func (i *IndexRepository) Set(key string, value string) {
	i.Index[key] = value
	i.Chan <- map[string]string{key: value}
}
