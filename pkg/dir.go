package pkg

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

func Touch(name string) (string, error) {
	dir := fmt.Sprintf("./problem/%s", name)
	if err := os.Mkdir(dir, 0755); err != nil {
		// 不是文件已存在错误
		if !os.IsExist(err) {
			log.Error().Msgf("Touch error: %v", err)
			return "", err
		}
	}
	return dir, nil
}
