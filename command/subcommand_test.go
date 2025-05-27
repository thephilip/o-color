package command

import (
	"testing"

	"github.com/hidetatz/kubecolor/kubectl"
	"github.com/hidetatz/kubecolor/testutil"
)

func Test_ResolveSubcommand(t *testing.T) {
	tests := []struct {
		name                   string
		args                   []string
		conf                   *KubecolorConfig
		isOutputTerminal       func() bool
		expectedShouldColorize bool
		expectedInfo           *kubectl.CLICommandInfo
	}{
		{
			name:             "basic case",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:          false,
				DarkBackground: true,
				ForceColor:     false,
				KubectlCmd:     "kubectl",
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "when plain, it won't colorize",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:          true,
				DarkBackground: true,
				ForceColor:     false,
				KubectlCmd:     "kubectl",
			},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "when help, it will colorize",
			args:             []string{"get", "pods", "-h"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:          false,
				DarkBackground: true,
				ForceColor:     false,
				KubectlCmd:     "kubectl",
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Get, Help: true},
		},
		{
			name:             "when both plain and force, plain is chosen",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:          true,
				DarkBackground: true,
				ForceColor:     true,
				KubectlCmd:     "kubectl",
			},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "when no subcommand is found, it becomes help",
			args:             []string{},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:          false,
				DarkBackground: true,
				ForceColor:     false,
				KubectlCmd:     "kubectl",
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.CLICommandInfo{Help: true},
		},
		{
			name:             "when the internal argument is found, it won't colorize",
			args:             []string{"__completeNoDesc", "get", "pods"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:          false,
				DarkBackground: true,
				ForceColor:     false,
				KubectlCmd:     "kubectl",
			},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "when not tty, it won't colorize",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return false },
			conf: &KubecolorConfig{
				Plain:          false,
				DarkBackground: true,
				ForceColor:     false,
				KubectlCmd:     "kubectl",
			},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "even if not tty, if force, it colorizes",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return false },
			conf: &KubecolorConfig{
				Plain:          false,
				DarkBackground: true,
				ForceColor:     true,
				KubectlCmd:     "kubectl",
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "kubectl edit is unsupported",
			args:             []string{"edit", "deployment"},
			isOutputTerminal: func() bool { return true },
			conf:             &KubecolorConfig{},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Edit},
		},
		{
			name:             "oc projects is supported",
			args:             []string{"projects"},
			isOutputTerminal: func() bool { return true },
			conf:             &KubecolorConfig{},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Projects},
		},
		{
			name:             "oc status is supported",
			args:             []string{"status"},
			isOutputTerminal: func() bool { return true },
			conf:             &KubecolorConfig{},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Status},
		},
		{
			name:             "oc new-project is unsupported",
			args:             []string{"new-project", "myproject"},
			isOutputTerminal: func() bool { return true },
			conf:             &KubecolorConfig{},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.NewProject},
		},
		{
			name:             "oc new-app is unsupported",
			args:             []string{"new-app", "nginx"},
			isOutputTerminal: func() bool { return true },
			conf:             &KubecolorConfig{},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.NewApp},
		},
		{
			name:             "oc routes is supported",
			args:             []string{"routes"},
			isOutputTerminal: func() bool { return true },
			conf:             &KubecolorConfig{},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Routes},
		},
		{
			name:             "oc policy is unsupported",
			args:             []string{"policy", "add-role-to-user", "edit", "user1"},
			isOutputTerminal: func() bool { return true },
			conf:             &KubecolorConfig{},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.CLICommandInfo{Subcommand: kubectl.Policy},
		},
		{
			name:             "when the subcommand is just -h (help), it will colorize",
			args:             []string{"-h"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:          false,
				DarkBackground: true,
				ForceColor:     false,
				KubectlCmd:     "kubectl",
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.CLICommandInfo{Help: true},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			isOutputTerminal = tt.isOutputTerminal
			shouldColorize, info := ResolveSubcommand(tt.args, tt.conf)
			testutil.MustEqual(t, tt.expectedShouldColorize, shouldColorize)
			testutil.MustEqual(t, tt.expectedInfo, info)
		})
	}
}
