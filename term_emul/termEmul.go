package termemul

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
)

//Read reads and returns terminal input
func ReadInput() (string, error) {
	const location = "termemul.ReadInput->"
	fmt.Print("Enter value here: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", logger.CreateDetails(location, err)
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	input = strings.TrimSuffix(input, "\r")

	return input, nil
}

// ====================================================================================
