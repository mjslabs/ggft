package askuser

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPkgAskUser(t *testing.T) {
	t.Run("GetInput", testGetInput)
}

func testGetInput(t *testing.T) {
	input := "Testing user input string"
	reader := bufio.NewReader(strings.NewReader(input + "\n"))

	value := getInput(reader)
	assert.Equal(t, value, input)

}
