package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"
)

func TestResolveVersion(t *testing.T) {
	tests := []struct {
		name          string
		linkerVersion string
		info          *debug.BuildInfo
		want          string
	}{
		{"release", "1.2.3", &debug.BuildInfo{Main: debug.Module{Version: "v1.0.0"}}, "1.2.3"},
		{"snapshot", "1.2.4-SNAPSHOT", nil, "1.2.4-SNAPSHOT"},
		{"go install", "", &debug.BuildInfo{Main: debug.Module{Version: "v1.2.3"}}, "v1.2.3"},
		{"prerelease", "", &debug.BuildInfo{Main: debug.Module{Version: "v1.2.3-rc.1"}}, "v1.2.3-rc.1"},
		{"local build", "", &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, "dev"},
		{"empty build info", "", &debug.BuildInfo{}, "dev"},
		{"missing build info", "", nil, "dev"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveVersion(tt.linkerVersion, tt.info); got != tt.want {
				t.Fatalf("resolveVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVersionOutput(t *testing.T) {
	info, _ := debug.ReadBuildInfo()
	want := "tsconfig-template version " + resolveVersion(version, info) + "\n"
	output, err := cliCommand(t, t.TempDir(), "--version").CombinedOutput()
	if err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, output)
	}
	if string(output) != want {
		t.Fatalf("--version = %q, want %q", output, want)
	}
}

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

func TestPromptFailure(t *testing.T) {
	dir := t.TempDir()
	// A subprocess without stdin cannot read from an interactive terminal.
	output, err := cliCommand(t, dir).CombinedOutput()
	if err == nil {
		t.Fatalf("CLI succeeded despite prompt failure:\n%s", output)
	}
	if !strings.Contains(string(output), "failed to ask question") {
		t.Fatalf("missing prompt error:\n%s", output)
	}
	if _, err := os.Stat(filepath.Join(dir, "tsconfig.json")); !os.IsNotExist(err) {
		t.Fatalf("expected no output file after prompt failure, got: %v", err)
	}
}
