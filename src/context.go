package src

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

var CurrentContext = Context{}

type Context struct {
	Chart      Chart
	Runner     Runner
	RunnerArgs string
	Verbose    bool
}
type Chart struct {
	ApiVersion string
	Name       string
	Version    string

	Scripts Scripts
}

func PrepareContext() (string, []string, error) {
	// parse flags and get actual args
	args := parseCommand()

	// read the Chart file
	chartDir := args[0]

	if CurrentContext.Verbose {
		fmt.Printf("#   Chart: %v\n", chartDir)
		fmt.Printf("#   Script: %v\n", args[1])
		fmt.Printf("#   Script Arguments: %s\n", strings.Join(args[2:], " "))
	}

	chartFile, err := os.ReadFile(chartDir + "/Chart.yaml")
	if err != nil {
		return "", nil, err
	}

	yaml.RegisterCustomUnmarshaler[Script](func(s *Script, bytes []byte) error {
		var node ast.Node
		err := yaml.Unmarshal(bytes, &node)
		if err != nil {
			return err
		}

		switch node.Type() {
		case ast.StringType:
			*s = SimpleScript(node.String())
			return nil
		case ast.MappingType:
			var rscript RunnerScript
			err = yaml.Unmarshal(bytes, &rscript)
			*s = rscript
			return err
		default:
			return errors.New("unknown script type, must either be a string or a map of runners to scripts")
		}
	})

	err = yaml.Unmarshal(chartFile, &CurrentContext.Chart)
	if err != nil {
		return "", nil, err
	}

	// return the script name and arguments
	return args[1], args[2:], nil
}

func parseCommand() []string {
	help := flag.Bool("help", false, "Show this help message")
	runner := flag.String("runner", defaultRunner(), "Runner to use for the script, should be one of: "+availableRunners())
	flag.BoolVar(&CurrentContext.Verbose, "verbose", false, "Show verbose output")
	flag.StringVar(&CurrentContext.RunnerArgs, "runner-args", "", "Arguments passed to the runner instead of the script")
	flag.Parse()

	if CurrentContext.Verbose {
		fmt.Printf("# helm-run\n")
		fmt.Printf("#   Verbose: %v\n", CurrentContext.Verbose)
		fmt.Printf("#   Help: %v\n", *help)
		fmt.Printf("#   Runner: %v\n", *runner)
		fmt.Printf("#   Runner Args: %v\n", CurrentContext.RunnerArgs)
	}

	if *help {
		printUsage()
		os.Exit(0)
	}

	CurrentContext.Runner = Runner(*runner)

	args := flag.Args()
	if len(args) < 2 {
		printUsage()
		os.Exit(0)
	}

	return args
}

func availableRunners() string {
	var names []string
	for _, runner := range Runners {
		names = append(names, string(runner))
	}

	return strings.Join(names, ", ")
}

func defaultRunner() string {
	switch runtime.GOOS {
	case "windows":
		return WindowsPowershell
	case "darwin":
		fallthrough
	case "linux":
		_, parentShell := path.Split(os.Getenv("SHELL"))
		if slices.Contains(Runners, Runner(parentShell)) {
			return parentShell
		}
	}

	return ""
}

func printUsage() {
	fmt.Printf("Usage: helm-run [flags] <chart> <script> [script-args]\n\n")
	fmt.Printf("Flags:\n")
	flag.PrintDefaults()
}
