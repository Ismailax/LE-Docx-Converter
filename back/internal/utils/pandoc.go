package utils

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
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
		"-v", dir+":/data",
		"pandoc/core:latest",
		inputInContainer, "-t", format, "-o", outputInContainer,
	)
	return cmd.Run()
}

func RunPandoc(input, output, format string) error {
	// ป้องกัน process ค้าง: timeout 2 นาที (ปรับได้ตามไฟล์ใหญ่/เล็ก)
	const defaultTimeout = 120 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// pandoc <input> -t <format> -o <output>
	cmd := exec.CommandContext(ctx, "pandoc", input, "-t", format, "-o", output)

	// รวม stdout+stderr เพื่อเดบักง่าย
	out, err := cmd.CombinedOutput()

	// แยกเคส timeout อธิบายชัด ๆ
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("pandoc timed out after %s", defaultTimeout)
	}
	if err != nil {
		return fmt.Errorf("pandoc error: %v | output: %s", err, string(out))
	}
	return nil
}
