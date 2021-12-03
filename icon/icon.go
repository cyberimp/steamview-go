package icon

import (
	"fmt"
	"io/ioutil"
)

func GetIcon() []byte {
	b, err := ioutil.ReadFile(iconPath)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
