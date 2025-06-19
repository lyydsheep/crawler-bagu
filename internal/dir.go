package internal

import (
	"fmt"
	"os"
)

func Touch(name string) (string, error) {
	dir := fmt.Sprintf(fmt.Sprintf("../problem/%s", name))
	return dir, os.Mkdir(dir, 0755)
}
