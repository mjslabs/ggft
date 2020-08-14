package askuser

import (
	"bufio"
	"io"
	"strings"
)

func getInput(descriptor io.Reader) string {
	reader := bufio.NewReader(descriptor)
	text, _ := reader.ReadString('\n')
	return strings.Replace(text, "\n", "", -1)
}
