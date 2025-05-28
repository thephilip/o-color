package kubectl

import (
	"strings"
)

type CLICommandInfo struct {
	Subcommand   CLICommand
	FormatOption FormatOption
	NoHeader     bool
	Watch        bool
	Help         bool
	Recursive    bool
	Short        bool

	IsKrew bool
	Args   []string
}

type FormatOption int

const (
	None FormatOption = iota
	Wide
	Json
	Yaml
)

type CLICommand int

const (
	Create CLICommand = iota + 1
	Expose
	Run
	Set
	Explain
	Get
	Edit
	Delete
	Rollout
	Scale
	Autoscale
	Certificate
	ClusterInfo
	Top
	Cordon
	Uncordon
	Drain
	Taint
	Describe
	Logs
	Attach
	Exec
	PortForward
	Proxy
	Cp
	Auth
	Diff
	Apply
	Patch
	Replace
	Wait
	Convert
	Kustomize
	Label
	Annotate
	Completion
	APIResources
	APIVersions
	Config
	Plugin
	Version
	Options
	Ctx
	Ns
	Debug
	// oc commands
	Projects
	Status
	NewProject
	NewApp
	Routes
	Policy
)

var strToCLICommand = map[string]CLICommand{
	"create":        Create,
	"expose":        Expose,
	"run":           Run,
	"set":           Set,
	"explain":       Explain,
	"get":           Get,
	"edit":          Edit,
	"delete":        Delete,
	"rollout":       Rollout,
	"scale":         Scale,
	"autoscale":     Autoscale,
	"certificate":   Certificate,
	"cluster-info":  ClusterInfo,
	"top":           Top,
	"cordon":        Cordon,
	"uncordon":      Uncordon,
	"drain":         Drain,
	"taint":         Taint,
	"describe":      Describe,
	"logs":          Logs,
	"attach":        Attach,
	"exec":          Exec,
	"port-forward":  PortForward,
	"proxy":         Proxy,
	"cp":            Cp,
	"auth":          Auth,
	"diff":          Diff,
	"apply":         Apply,
	"patch":         Patch,
	"replace":       Replace,
	"wait":          Wait,
	"convert":       Convert,
	"kustomize":     Kustomize,
	"label":         Label,
	"annotate":      Annotate,
	"completion":    Completion,
	"api-resources": APIResources,
	"api-versions":  APIVersions,
	"config":        Config,
	"plugin":        Plugin,
	"version":       Version,
	"options":       Options,
	"ctx":           Ctx,
	"ns":            Ns,
	"debug":         Debug,
	// oc commands
	"projects":    Projects,
	"status":      Status,
	"new-project": NewProject,
	"new-app":     NewApp,
	"routes":      Routes,
	"policy":      Policy,
}

func InspectCLICommand(command string) (CLICommand, bool) {
	sc, ok := strToCLICommand[command]

	return sc, ok
}

func CollectCommandlineOptions(args []string, info *CLICommandInfo) {
	for i := range args {
		if strings.HasPrefix(args[i], "--output") {
			switch args[i] {
			case "--output=json":
				info.FormatOption = Json
			case "--output=yaml":
				info.FormatOption = Yaml
			case "--output=wide":
				info.FormatOption = Wide
			default:
				if len(args)-1 > i {
					formatOption := args[i+1]
					switch formatOption {
					case "json":
						info.FormatOption = Json
					case "yaml":
						info.FormatOption = Yaml
					case "wide":
						info.FormatOption = Wide
					default:
						// custom-columns, go-template, etc are currently not supported
					}
				}
			}
		} else if strings.HasPrefix(args[i], "-o") {
			switch args[i] {
			// both '-ojson' and '-o=json' works
			case "-ojson", "-o=json":
				info.FormatOption = Json
			case "-oyaml", "-o=yaml":
				info.FormatOption = Yaml
			case "-owide", "-o=wide":
				info.FormatOption = Wide
			default:
				// otherwise, look for next arg because '-o json' also works
				if len(args)-1 > i {
					formatOption := args[i+1]
					switch formatOption {
					case "json":
						info.FormatOption = Json
					case "yaml":
						info.FormatOption = Yaml
					case "wide":
						info.FormatOption = Wide
					default:
						// custom-columns, go-template, etc are currently not supported
					}
				}

			}
		} else if strings.HasPrefix(args[i], "--short") {
			switch args[i] {
			case "--short=true":
				info.Short = true
			case "--short=false":
				info.Short = false
			default:
				info.Short = true
			}
		} else if args[i] == "--no-headers" {
			info.NoHeader = true
		} else if args[i] == "-w" || args[i] == "--watch" {
			info.Watch = true
		} else if args[i] == "--recursive=true" || args[i] == "--recursive" {
			info.Recursive = true
		} else if args[i] == "-h" || args[i] == "--help" {
			info.Help = true
		}
	}
}

// TODO: return shouldColorize = false when the given args is for plugin
func InspectCLICommandInfo(args []string) (*CLICommandInfo, bool) {
	ret := &CLICommandInfo{Args: args} // Store original args

	CollectCommandlineOptions(args, ret)

	for i := range args {
		cmd, ok := InspectCLICommand(args[i])
		if !ok {
			continue
		}

		ret.Subcommand = cmd
		return ret, true
	}

	return ret, false
}
