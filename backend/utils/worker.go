package utils

import "os"

func IsWorker() bool {
	return os.Getenv("IS_WORKER") != ""
}

func IsMaster() bool {
	return os.Getenv("IS_MASTER") != ""
}
