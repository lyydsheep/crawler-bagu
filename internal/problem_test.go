package internal

import (
	"crawler-bagu/pkg"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	url := fmt.Sprintf("https://www.bagujing.com/api/problems/%d?groupId=%d", 14488, 23)
	body, err := get(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(body)
}

func TestFetchAndSaveAnswer(t *testing.T) {
	dir, err := pkg.Touch("go-23")
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}
	if _, err = fetchAndSaveAnswer(dir, 23, 14488); err != nil {
		log.Fatal(err)
	}
	log.Println("successful")
}
