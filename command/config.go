package command

import "os"

type KubecolorConfig struct {
	Plain                bool
	DarkBackground       bool
	ForceColor           bool
	ShowKubecolorVersion bool
	KubectlCmd           string
	UseOcCli             bool
}

func ResolveConfig(args []string) ([]string, *KubecolorConfig) {
	args, plainFlagFound := findAndRemoveBoolFlagIfExists(args, "--plain")
	args, lightBackgroundFlagFound := findAndRemoveBoolFlagIfExists(args, "--light-background")
	args, forceColorFlagFound := findAndRemoveBoolFlagIfExists(args, "--force-colors")
	args, kubecolorVersionFlagFound := findAndRemoveBoolFlagIfExists(args, "--kubecolor-version")
	args, useOcCliFlagFound := findAndRemoveBoolFlagIfExists(args, "--use-oc-cli")

	darkBackground := !lightBackgroundFlagFound

	kubectlCmd := "kubectl"
	if useOcCliFlagFound {
		kubectlCmd = "oc"
	} else {
		if kc := os.Getenv("KUBECTL_COMMAND"); kc != "" {
			kubectlCmd = kc
		}
	}

	return args, &KubecolorConfig{
		Plain:                plainFlagFound,
		DarkBackground:       darkBackground,
		ForceColor:           forceColorFlagFound,
		ShowKubecolorVersion: kubecolorVersionFlagFound,
		KubectlCmd:           kubectlCmd,
		UseOcCli:             useOcCliFlagFound,
	}
}

func findAndRemoveBoolFlagIfExists(args []string, key string) ([]string, bool) {
	for i, arg := range args {
		if arg == key {
			return append(args[:i], args[i+1:]...), true
		}
	}

	return args, false
}
