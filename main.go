package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	filename := "README.md"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(os.Stderr, err)
		return
	}
	Parse(Lex(data))
}

/* NOTE(kra53n): wait for parsing
const (
	noteLowerS = "note"
	noteUpperS = "NOTE"

	tipLowerS = "tip"
	tipUpperS = "TIP"

	importantLowerS = "important"
	importantUpperS = "IMPORTANT"

	warningLowerS = "warning"
	warningUpperS = "WARNING"

	cautionLowerS = "caution"
	cautionUpperS = "CAUTION"
)
*/
