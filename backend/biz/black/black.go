package black

import (
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	blackMap     = map[string]bool{}
	whiteMap     = map[string]bool{}
	currBlack    = []string{}
	currBlackMap = map[string]bool{}
)

func Init(blackPath, whitePath string) {
	InitBlack(blackPath)
	InitWhite(whitePath)
}

func InitWhite(whiteFilePath string) {
	if !utils.HasFile(whiteFilePath) {
		log.Printf("WARNING: white file not found")
		//panic("black file not found")
	}
	fileutil.ReadByLine(whiteFilePath, func(input string) error {
		//log.Printf("insert white:%v", input)
		whiteMap[input] = true
		return nil
	})
}

func InitBlack(blackFilePath string) {
	if !utils.HasFile(blackFilePath) {
		log.Printf("WARNING: black file not found")
		//panic("black file not found")
	}
	fileutil.ReadByLine(blackFilePath, func(input string) error {
		//log.Printf("insert black:%v", input)
		blackMap[input] = true
		return nil
	})

	fileutil.ReadByLine("chore/black_record.txt", func(s string) error {
		currBlack = append(currBlack, s)
		return nil
	})

}

func IsWhite(input string) bool {
	for k := range whiteMap {
		if strings.Contains(input, k) {
			log.Printf("hit WHITE: k=%v input=%v", k, input)
			return true
		}
	}
	return false
}

func HitBlack(input, desc string) string {

	if input == "" {
		return ""
	}

	if IsWhite(input) {
		return ""
	}

	for k := range blackMap {
		if strings.Contains(input, k) {
			//log.Printf("hit BLACK: k=%v input=%v", k, input)
			resp := fmt.Sprintf("[black=%v source=%v desc= %v ]", k, input, desc)
			if !currBlackMap[input] {
				currBlackMap[input] = true
				currBlack = append(currBlack, resp)
				if len(currBlack)%10 == 0 {
					fileutil.WriteToFile([]byte(strings.Join(currBlack, "\n")), "chore/black_record.txt")
				}
			}
			return resp
		}
	}

	return ""
}
