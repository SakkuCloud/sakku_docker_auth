package util

import (
	"bufio"
	"os"
	"strings"
)

const (
	SuccessExitCode = 0
	ErrorExitCode   = 1
)

func ReadStdIn() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
}
