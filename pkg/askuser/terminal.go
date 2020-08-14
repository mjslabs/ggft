package askuser

import (
	"fmt"
	"os"
)

// Terminal - get input from terminal
func Terminal(prompt, defaultText string) string {
	fmt.Print(prompt)
	return getInput(os.Stdin)
}
