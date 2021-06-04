package scan

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os/exec"
)

// GetSHA256 calculates a file's sha256sum
func GetSHA256(name string) string {

	dat, err := ioutil.ReadFile(name)
	if err != nil {
		return ""
	}

	h256 := sha256.New()
	_, err = h256.Write(dat)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h256.Sum(nil))
}

// RunCommand runs cmd on file
func RunCommand(ctx context.Context, cmd string, args ...string) (string, error) {

	var c *exec.Cmd

	if ctx != nil {
		c = exec.CommandContext(ctx, cmd, args...)
	} else {
		c = exec.Command(cmd, args...)
	}

	output, err := c.Output()
	if err != nil {
		return string(output), err
	}

	// check for exec context timeout
	if ctx != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("command %s timed out", cmd)
		}
	}

	return string(output), nil
}
