package utils

import (
	"fmt"
	"os"
)

func IsWorker() bool {
	fmt.Printf("IS_WORKER env:[%v]", os.Getenv("IS_WORKER"))
	return os.Getenv("IS_WORKER") == "true"
}

func IsMaster() bool {
	fmt.Printf("IS_MASTER env:[%v]", os.Getenv("IS_MASTER"))
	return os.Getenv("IS_MASTER") == "true"
}
