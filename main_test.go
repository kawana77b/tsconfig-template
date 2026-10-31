package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Run the real CLI in a subprocess so its exit status can be checked.
func TestCLIProcess(t *testing.T) {
	if os.Getenv("TSCONFIG_TEMPLATE_TEST_PROCESS") != "1" {
		return
	}
	for i, arg := range os.Args {
		if arg == "--" {
			os.Args = append([]string{"tsconfig-template"}, os.Args[i+1:]...)
			main()
			os.Exit(0)
		}
	}
	t.Fatal("missing CLI argument separator")
}

func cliCommand(t *testing.T, dir string, args ...string) *exec.Cmd {
	t.Helper()
	command := exec.Command(os.Args[0], append([]string{"-test.run=^TestCLIProcess$", "--"}, args...)...)
	command.Env = append(os.Environ(), "TSCONFIG_TEMPLATE_TEST_PROCESS=1")
	command.Dir = dir
	return command
}

func TestDefaultTemplateOutput(t *testing.T) {
	dir := t.TempDir()
	output, err := cliCommand(t, dir, "--default").CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, output)
	}
	got, err := os.ReadFile(filepath.Join(dir, "tsconfig.json"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := fs.ReadFile("bases/bases/recommended.json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("generated template differs from embedded recommended.json")
	}
}

func TestDefaultTemplateWriteFailure(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "tsconfig.json"), 0755); err != nil {
		t.Fatal(err)
	}
	output, err := cliCommand(t, dir, "--default").CombinedOutput()
	if err == nil {
		t.Fatalf("CLI succeeded despite write failure:\n%s", output)
	}
	if !strings.Contains(string(output), "failed to write tsconfig") {
		t.Fatalf("missing write error:\n%s", output)
	}
	if strings.Contains(string(output), "created:") {
		t.Fatalf("CLI reported successful creation after write failure:\n%s", output)
	}
}
