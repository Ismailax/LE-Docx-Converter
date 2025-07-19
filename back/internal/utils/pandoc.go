package utils

import (
	"os/exec"
	"path/filepath"
)

// RunPandocDocker รัน pandoc ผ่าน docker CLI (mount tmp dir แบบ absolute path)
func RunPandocDocker(input, output, format string) error {
	dir, err := filepath.Abs(filepath.Dir(input))
	if err != nil {
		return err
	}
	inputInContainer := "/data/" + filepath.Base(input)
	outputInContainer := "/data/" + filepath.Base(output)
	cmd := exec.Command(
		"docker", "run", "--rm",
		"-v", dir + ":/data",
		"pandoc/core:3.7.0.2",
		inputInContainer, "-t", format, "-o", outputInContainer,
	)
	return cmd.Run()
}