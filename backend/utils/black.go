package utils

import (
	"fmt"
	"github.com/logxxx/utils/fileutil"
	"os"
	"strings"
)

func WriteBlackUid(uid string, reason string) {
	fileutil.AppendToFile("chore/black_uids.txt", fmt.Sprintf("\n%v //%v", uid, reason))
}

func IsBlackUid(uid string) bool {
	data, _ := os.ReadFile("chore/black_uids.txt")
	return strings.Contains(string(data), uid)
}
