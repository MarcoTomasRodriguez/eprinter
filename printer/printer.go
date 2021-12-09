package printer

import (
	"os/exec"
)

// PrintFile prints a file using lp.
// See: https://man7.org/linux/man-pages/man1/lp.1.html and https://www.cups.org/doc/options.html.
func PrintFile(filename string) error {
	return exec.Command("lp", "--", filename).Run()
}
