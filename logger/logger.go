package logger

import (
	"github.com/utilyre/role/config"
	"go.uber.org/zap"
)

func New(c config.Config) (*zap.Logger, error) {
	switch c.Mode {
	case config.ModeDev:
		return zap.NewDevelopment()
	case config.ModeProd:
		return zap.NewProduction()
	}

	return nil, nil
}
