package crawler

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

func checkExist(dir string, groupId, problemId int) (string, bool) {
	filename := fmt.Sprintf("%s/%d_%d.json", dir, groupId, problemId)
	if _, err := os.Stat(filename); err == nil {
		log.Info().Msgf("answer already exists, skip %s", filename)
		return "", true
	} else if !os.IsNotExist(err) {
		// 如果错误不是“文件不存在”，则返回错误
		log.Fatal().Msgf("os.Stat error: %v", err)
	}
	return filename, false
}

func saveAnswer(filename string, ans []byte) error {
	err := os.WriteFile(filename, ans, 0644)
	if err != nil {
		log.Error().Msgf("os.WriteFile error: %v", err)
	}
	return err
}
