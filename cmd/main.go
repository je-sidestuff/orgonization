package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/je-sidestuff/orgonization/templates"
)

func handleBuiltinTemplate(builtin string) {
	switch builtin {
	case "weekly":
		handleWeeklyTemplate()
	default:
		fmt.Printf("Unknown builtin template: %s\n", builtin)
		os.Exit(1)
	}
}

func handleWeeklyTemplate() {

	// Get the builtin weekly template directive
	weeklyDirective := templates.GetWeeklyTemplateDirective()

	// Create the template processor and initialize the filesystem and directive
	templateProcessor := templates.NewTemplateProcessor("PrimaryTemplateProcessor")
	// We will probably want to change this to 'empty' so we don't get extra directives
	// For now let's not initialize it at all.
	//templateProcessor.InitializeFilesystem(templates.GetDefaultFilesystemConfiguration())

	templateProcessor.InitializeDirectives([]templates.DirectiveConfiguration{weeklyDirective})

	// Complete initialization then perform our single tick
	templateProcessor.CompleteInitialization()
	templateProcessor.Tick()
}

func handleOnce() {

	// Load config file (fake some of this on the first pass)
	// What will this give us?
	// - Locations of the IO files
	// - List of directives

	filesystemConfiguration, err := templates.GetDefaultFilesystemConfiguration()

	if err != nil {
		panic(err)
	}

	// Get the default template directives
	defaultDirectives := templates.GetDefaultTemplateDirectives()

	templateProcessor := templates.NewTemplateProcessor("PrimaryTemplateProcessor")

	templateProcessor.InitializeFilesystem(filesystemConfiguration)

	templateProcessor.InitializeDirectives(defaultDirectives)

	templateProcessor.CompleteInitialization()
	templateProcessor.Tick()

	// // Get the builtin weekly template directive
	// weeklyDirective := templates.GetWeeklyTemplateDirective()

	// // Create the template processor and initialize the filesystem and directive
	// templateProcessor := templates.NewTemplateProcessor("PrimaryTemplateProcessor")
	// templateProcessor.InitializeFilesystem(templates.GetDefaultFilesystemConfiguration())
	// templateProcessor.InitializeDirectives([]templates.DirectiveConfiguration{weeklyDirective})

	// // Complete initialization then perform our single tick
	// templateProcessor.CompleteInitialization()
	// templateProcessor.Tick()
}

func handleAgent() {

	// Load config file (fake some of this on the first pass)
	// What will this give us?
	// - Locations of the IO files
	// - List of directives

	filesystemConfiguration, err := templates.GetDefaultFilesystemConfiguration()

	if err != nil {
		panic(err)
	}

	// Get the default template directives
	defaultDirectives := templates.GetDefaultTemplateDirectives()

	templateProcessor := templates.NewTemplateProcessor("PrimaryTemplateProcessor")

	templateProcessor.InitializeFilesystem(filesystemConfiguration)

	templateProcessor.InitializeDirectives(defaultDirectives)

	templateProcessor.CompleteInitialization()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	for {

		// Get the cullent time in microseconds
		currentTimeUs := time.Now().UnixMicro()

		templateProcessor.Tick()

		// Print the time delta
		logger.Debug("Time delta",
			slog.Int64("delta", time.Now().UnixMicro()-currentTimeUs))

		time.Sleep(time.Millisecond * 100)
	}
}

func main() {

	// Define the main command and subcommands as separate flag sets
	var mainCmd = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Usage comment explaining how to use the CLI tool
	usage := `

	The goorganizethings module is a Go executable capable of running one-off CLI commands, as a local-only agent, or as a server.

	Usage:
	%s <subcommand> [flags]

	Available Subcommands:
	template: Outputs a specified populated template or performs a one-off processing of an IO file tree
	agent:    Runs continuously as an agent, watching an IO file tree and performing updates according to its configuration
	server:   Runs as an agent backed by a server, allowing a client to make updates while IO file tree processing is conducted

	Flags:
	-h, --help    show help message

	**Subcommand-specific flags are available. Run the subcommand with the "-h" flag for details.**

`

	// Help flag
	var help bool
	mainCmd.BoolVar(&help, "h", false, "show help message")
	mainCmd.BoolVar(&help, "help", false, "show help message (shorthand)")

	// Print usage message and exit if help flag is set or arguments are incorrect
	if help || len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}

	// Define template subcommand and its flags
	templateCmd := flag.NewFlagSet("template", flag.ExitOnError)
	var buildTarget string
	templateCmd.StringVar(&buildTarget, "file", "", "Top level file to run templating on.")
	var buildVerbose bool
	templateCmd.BoolVar(&buildVerbose, "verbose", false, "Enable verbose output")
	var templateFlag1 string
	templateCmd.StringVar(&templateFlag1, "recurse", "", "Whether to recurse templating to generated files.")
	var builtin string
	templateCmd.StringVar(&builtin, "builtin", "", "Builtin template to use (optional)")
	var once bool
	templateCmd.BoolVar(&once, "once", false, "Run templating for the configured files once and exit")

	// Parse command line arguments with main command
	if err := mainCmd.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	subcmd := os.Args[1]

	// Parse subcommand specific flags based on chosen subcommand
	switch subcmd {
	case "template":
		if err := templateCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Handle builtin argument
		if builtin != "" {
			handleBuiltinTemplate(builtin)
		}

		if once {
			fmt.Println("Processing all files once.")
			handleOnce()
		}

	case "agent":
		if err := templateCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		handleAgent()

	default:
		fmt.Fprintf(os.Stderr, "Invalid subcommand: %s", os.Args[1])
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}
}
