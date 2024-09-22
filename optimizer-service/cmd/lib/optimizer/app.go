package optimizer

import (
	"io"
	"optimizer-service/cmd/internal/storage"
	"optimizer-service/cmd/internal/utils"
)

type IOptimizer interface {
	Optimize(key string) (io.ReadCloser, error)
	SupportedFormats() []string
}

type Optimizer struct {
	Storage storage.Storage
	Utils   utils.IUtils
}
type OptimizerParams struct {
	Level *string
}

func InitOptimizer(storage storage.Storage, u *utils.Utils) *Optimizer {
	return &Optimizer{
		Storage: storage,
		Utils:   u,
	}
}
