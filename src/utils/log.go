package utils

import (
	"os"

	"github.com/op/go-logging"
)

var logger *logging.Logger

func Log() *logging.Logger {
	return logger
}

func InitLogger() {
	logger = logging.MustGetLogger("Walkr")
	format := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfile} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
	)

	stdOutput := logging.NewLogBackend(os.Stderr, "", 0)
	stdOutputFormatter := logging.NewBackendFormatter(stdOutput, format)

	logging.SetBackend(stdOutputFormatter)

}
