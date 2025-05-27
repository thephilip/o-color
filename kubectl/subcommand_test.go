package kubectl

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// func TestInspectCLICommandInfo(args []string) (*CLICommandInfo, bool) {
func TestInspectCLICommandInfo(t *testing.T) {
	tests := []struct {
		name       string
		args       string
		expected   *CLICommandInfo
		expectedOK bool
	}{
		{"get pods", "get pods", &CLICommandInfo{Subcommand: Get, Args: []string{"get", "pods"}}, true},
		{"get pod", "get pod", &CLICommandInfo{Subcommand: Get, Args: []string{"get", "pod"}}, true},
		{"get po", "get po", &CLICommandInfo{Subcommand: Get, Args: []string{"get", "po"}}, true},

		{"get pod -o wide", "get pod -o wide", &CLICommandInfo{Subcommand: Get, FormatOption: Wide, Args: []string{"get", "pod", "-o", "wide"}}, true},
		{"get pod -o=wide", "get pod -o=wide", &CLICommandInfo{Subcommand: Get, FormatOption: Wide, Args: []string{"get", "pod", "-o=wide"}}, true},
		{"get pod -owide", "get pod -owide", &CLICommandInfo{Subcommand: Get, FormatOption: Wide, Args: []string{"get", "pod", "-owide"}}, true},

		{"get pod -o json", "get pod -o json", &CLICommandInfo{Subcommand: Get, FormatOption: Json, Args: []string{"get", "pod", "-o", "json"}}, true},
		{"get pod -o=json", "get pod -o=json", &CLICommandInfo{Subcommand: Get, FormatOption: Json, Args: []string{"get", "pod", "-o=json"}}, true},
		{"get pod -ojson", "get pod -ojson", &CLICommandInfo{Subcommand: Get, FormatOption: Json, Args: []string{"get", "pod", "-ojson"}}, true},

		{"get pod -o yaml", "get pod -o yaml", &CLICommandInfo{Subcommand: Get, FormatOption: Yaml, Args: []string{"get", "pod", "-o", "yaml"}}, true},
		{"get pod -o=yaml", "get pod -o=yaml", &CLICommandInfo{Subcommand: Get, FormatOption: Yaml, Args: []string{"get", "pod", "-o=yaml"}}, true},
		{"get pod -oyaml", "get pod -oyaml", &CLICommandInfo{Subcommand: Get, FormatOption: Yaml, Args: []string{"get", "pod", "-oyaml"}}, true},

		{"get pod --output json", "get pod --output json", &CLICommandInfo{Subcommand: Get, FormatOption: Json, Args: []string{"get", "pod", "--output", "json"}}, true},
		{"get pod --output=json", "get pod --output=json", &CLICommandInfo{Subcommand: Get, FormatOption: Json, Args: []string{"get", "pod", "--output=json"}}, true},
		{"get pod --output yaml", "get pod --output yaml", &CLICommandInfo{Subcommand: Get, FormatOption: Yaml, Args: []string{"get", "pod", "--output", "yaml"}}, true},
		{"get pod --output=yaml", "get pod --output=yaml", &CLICommandInfo{Subcommand: Get, FormatOption: Yaml, Args: []string{"get", "pod", "--output=yaml"}}, true},
		{"get pod --output wide", "get pod --output wide", &CLICommandInfo{Subcommand: Get, FormatOption: Wide, Args: []string{"get", "pod", "--output", "wide"}}, true},
		{"get pod --output=wide", "get pod --output=wide", &CLICommandInfo{Subcommand: Get, FormatOption: Wide, Args: []string{"get", "pod", "--output=wide"}}, true},

		{"get pod --no-headers", "get pod --no-headers", &CLICommandInfo{Subcommand: Get, NoHeader: true, Args: []string{"get", "pod", "--no-headers"}}, true},
		{"get pod -w", "get pod -w", &CLICommandInfo{Subcommand: Get, Watch: true, Args: []string{"get", "pod", "-w"}}, true},
		{"get pod --watch", "get pod --watch", &CLICommandInfo{Subcommand: Get, Watch: true, Args: []string{"get", "pod", "--watch"}}, true},
		{"get pod -h", "get pod -h", &CLICommandInfo{Subcommand: Get, Help: true, Args: []string{"get", "pod", "-h"}}, true},
		{"get pod --help", "get pod --help", &CLICommandInfo{Subcommand: Get, Help: true, Args: []string{"get", "pod", "--help"}}, true},

		{"describe pod pod-aaa", "describe pod pod-aaa", &CLICommandInfo{Subcommand: Describe, Args: []string{"describe", "pod", "pod-aaa"}}, true},
		{"top pod", "top pod", &CLICommandInfo{Subcommand: Top, Args: []string{"top", "pod"}}, true},
		{"top pods", "top pods", &CLICommandInfo{Subcommand: Top, Args: []string{"top", "pods"}}, true},

		{"api-versions", "api-versions", &CLICommandInfo{Subcommand: APIVersions, Args: []string{"api-versions"}}, true},

		{"explain pod", "explain pod", &CLICommandInfo{Subcommand: Explain, Args: []string{"explain", "pod"}}, true},
		{"explain pod --recursive=true", "explain pod --recursive=true", &CLICommandInfo{Subcommand: Explain, Recursive: true, Args: []string{"explain", "pod", "--recursive=true"}}, true},
		{"explain pod --recursive", "explain pod --recursive", &CLICommandInfo{Subcommand: Explain, Recursive: true, Args: []string{"explain", "pod", "--recursive"}}, true},

		{"version", "version", &CLICommandInfo{Subcommand: Version, Args: []string{"version"}}, true},
		{"version --client", "version --client", &CLICommandInfo{Subcommand: Version, Args: []string{"version", "--client"}}, true},
		{"version --short", "version --short", &CLICommandInfo{Subcommand: Version, Short: true, Args: []string{"version", "--short"}}, true},
		{"version -o json", "version -o json", &CLICommandInfo{Subcommand: Version, FormatOption: Json, Args: []string{"version", "-o", "json"}}, true},
		{"version -o yaml", "version -o yaml", &CLICommandInfo{Subcommand: Version, FormatOption: Yaml, Args: []string{"version", "-o", "yaml"}}, true},

		{"apply", "apply", &CLICommandInfo{Subcommand: Apply, Args: []string{"apply"}}, true},

		// oc commands
		{"projects", "projects", &CLICommandInfo{Subcommand: Projects, Args: []string{"projects"}}, true},
		{"status", "status", &CLICommandInfo{Subcommand: Status, Args: []string{"status"}}, true},
		{"status -v", "status -v", &CLICommandInfo{Subcommand: Status, Args: []string{"status", "-v"}}, true},
		{"new-project", "new-project myproject", &CLICommandInfo{Subcommand: NewProject, Args: []string{"new-project", "myproject"}}, true},
		{"new-app", "new-app nginx", &CLICommandInfo{Subcommand: NewApp, Args: []string{"new-app", "nginx"}}, true},
		{"routes", "routes", &CLICommandInfo{Subcommand: Routes, Args: []string{"routes"}}, true},
		{"get routes", "get routes", &CLICommandInfo{Subcommand: Get, Args: []string{"get", "routes"}}, true},
		{"policy add-role-to-user", "policy add-role-to-user edit user1", &CLICommandInfo{Subcommand: Policy, Args: []string{"policy", "add-role-to-user", "edit", "user1"}}, true},

		{"empty", "", &CLICommandInfo{Args: []string{""}}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, ok := InspectCLICommandInfo(strings.Split(tt.args, " "))
			if tt.expectedOK != ok {
				t.Errorf("expectedOK got %v, want %v for args: %s", ok, tt.expectedOK, tt.args)
			}

			// Note: Comparing CLICommandInfo.Args directly in the existing table might be complex
			// due to how tt.args is a single string and then split, versus expected.Args being a slice.
			// The cmp.Diff should handle it if expected.Args is correctly populated.
			// All test cases above have been updated to include the Args field.
			if diff := cmp.Diff(tt.expected, s); diff != "" {
				t.Errorf("InspectCLICommandInfo() mismatch (-want +got):\n%s for args: %s", diff, tt.args)
			}
		})
	}
}

func TestCLICommandInfoArgsPopulationSeparate(t *testing.T) {
	tests := []struct {
		name                 string
		args                 []string
		expectedCmd          CLICommand
		expectedArgs         []string
		expectedFormatOption FormatOption
		expectedOK           bool
	}{
		{
			name:           "get pods with namespace",
			args:           []string{"get", "pods", "--namespace=test-ns"},
			expectedCmd:    Get,
			expectedArgs:   []string{"get", "pods", "--namespace=test-ns"},
			expectedOK:     true,
		},
		{
			name:           "version command",
			args:           []string{"version"},
			expectedCmd:    Version,
			expectedArgs:   []string{"version"},
			expectedOK:     true,
		},
		{
			name:           "empty args",
			args:           []string{},
			expectedCmd:    0, // Zero value for CLICommand
			expectedArgs:   []string{},
			expectedOK:     false, // No subcommand found
		},
		{
			name:           "describe pod command with options",
			args:           []string{"describe", "pod", "my-pod", "-n", "test"},
			expectedCmd:    Describe,
			expectedArgs:   []string{"describe", "pod", "my-pod", "-n", "test"},
			expectedOK:     true,
		},
		{
			name:                 "get svc with output format",
			args:                 []string{"get", "svc", "-o", "yaml"},
			expectedCmd:          Get,
			expectedArgs:         []string{"get", "svc", "-o", "yaml"},
			expectedFormatOption: Yaml,
			expectedOK:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, ok := InspectCLICommandInfo(tt.args)

			if ok != tt.expectedOK {
				t.Errorf("For args '%v', expectedOK %v, got %v", tt.args, tt.expectedOK, ok)
			}

			if !reflect.DeepEqual(info.Args, tt.expectedArgs) {
				t.Errorf("For args '%v', expected Args %v, got %v", tt.args, tt.expectedArgs, info.Args)
			}
			if info.Subcommand != tt.expectedCmd {
				t.Errorf("For args '%v', expected Subcommand %v, got %v", tt.args, tt.expectedCmd, info.Subcommand)
			}

			// Check FormatOption only if the command is expected to parse it and it's relevant for the test case
			if tt.expectedFormatOption != 0 { // 0 is None, the default
				if info.FormatOption != tt.expectedFormatOption {
					t.Errorf("For args '%v', expected FormatOption %v, got %v", tt.args, tt.expectedFormatOption, info.FormatOption)
				}
			}
		})
	}
}
