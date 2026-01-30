package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration_SingleFile_AllMetrics(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	output := runCommand(t, binary, "testdata/simple.txt")

	expected := "      3       8      38 testdata/simple.txt\n"
	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestIntegration_SingleFile_Lines(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	output := runCommand(t, binary, "-l", "testdata/simple.txt")

	expected := "      3 testdata/simple.txt\n"
	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestIntegration_SingleFile_WordsAndBytes(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	output := runCommand(t, binary, "-w", "-c", "testdata/simple.txt")

	expected := "      8      38 testdata/simple.txt\n"
	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestIntegration_MultipleFiles_WithTotal(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	output := runCommand(t, binary, "testdata/simple.txt", "testdata/empty.txt")

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines (2 files + total), got %d", len(lines))
	}

	if !strings.Contains(lines[2], "total") {
		t.Errorf("Expected third line to contain 'total', got: %s", lines[2])
	}

	expected := "      3       8      38 total"
	if lines[2] != expected {
		t.Errorf("Expected total:\n%s\nGot:\n%s", expected, lines[2])
	}
}

func TestIntegration_EmptyFile(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	output := runCommand(t, binary, "testdata/empty.txt")

	expected := "      0       0       0 testdata/empty.txt\n"
	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestIntegration_UTF8_Chars(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	output := runCommand(t, binary, "-m", "testdata/utf8.txt")

	expected := "     21 testdata/utf8.txt\n"
	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestIntegration_UTF8_Bytes(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	output := runCommand(t, binary, "-c", "testdata/utf8.txt")

	expected := "     34 testdata/utf8.txt\n"
	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestIntegration_NonexistentFile(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	cmd := exec.Command(binary, "nonexistent.txt")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Fatal("Expected error for nonexistent file")
	}

	stderrOutput := stderr.String()
	if !strings.Contains(stderrOutput, "nonexistent.txt") {
		t.Errorf("Expected error message to contain filename, got: %s", stderrOutput)
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()

	binary := filepath.Join(t.TempDir(), "wc")
	cmd := exec.Command("go", "build", "-o", binary, ".")
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, output)
	}

	return binary
}

func runCommand(t *testing.T, binary string, args ...string) string {
	t.Helper()

	cmd := exec.Command(binary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
	}

	return stdout.String()
}
