package file

import (
	"github.com/go-ceres/go-ceres/config"
	"github.com/go-ceres/go-ceres/logger"
)

type Option func(f *fileSource)

func Unmarshal(unmarshal string) Option {
	return func(f *fileSource) {
		if _, ok := config.Unmarshals[unmarshal]; ok {
			logger.FrameLogger.Panicd("set unmarshal err: no unmarshal name")
		}
		f.unmarshal = unmarshal
	}
}
