package task

import (
	"fmt"
	"os"
	"os/exec"
)

func isFileExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func execCmd(cmd ...string) error {
	binary, lookErr := exec.LookPath(cmd[0])
	if lookErr != nil {
		return fmt.Errorf("command not found: %s", cmd[0])
	}

	c := exec.Command(binary, cmd[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}
