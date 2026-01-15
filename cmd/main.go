package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/je-sidestuff/orgonization/server"

	"github.com/je-sidestuff/orgonization/templates"
)

func configureLogLevel(levelStr string) slog.Level {
	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		fmt.Fprintf(os.Stderr, "Invalid log level: %s. Using 'info' instead.\n", levelStr)
		level = slog.LevelInfo
	}

	// Set the default logger with the configured level
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))

	return level
}

func handleBuiltinTemplate(builtin string, logLevel slog.Level) {
	switch builtin {
	case "weekly":
		handleWeeklyTemplate(logLevel)
	default:
		fmt.Printf("Unknown builtin template: %s\n", builtin)
		os.Exit(1)
	}
}

func handleWeeklyTemplate(logLevel slog.Level) {

	// Get the builtin weekly template directive
	weeklyDirective := templates.GetWeeklyTemplateDirective()

	// Create the template processor and initialize the filesystem and directive
	templateProcessor := templates.NewTemplateProcessor("PrimaryTemplateProcessor", logLevel)

	filesystemConfiguration, err := templates.GetDefaultFilesystemConfiguration()

	if err != nil {
		panic(err)
	}

	templateProcessor.InitializeFilesystem(filesystemConfiguration)

	templateProcessor.InitializeDirectives([]templates.DirectiveConfiguration{weeklyDirective})

	// Complete initialization then perform our single tick
	templateProcessor.CompleteInitialization()

	// In a future increment we'll fix this up to more dependency-injecty
	templateProcessor.InjectSysout("<WEEKLY_NOTES:>")
	templateProcessor.Tick()
}

func handleOnce(logLevel slog.Level) {

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

	templateProcessor := templates.NewTemplateProcessor("PrimaryTemplateProcessor", logLevel)

	templateProcessor.InitializeFilesystem(filesystemConfiguration)

	templateProcessor.InitializeDirectives(defaultDirectives)

	templateProcessor.CompleteInitialization()
	templateProcessor.Tick()

	// // Get the builtin weekly template directive
	// weeklyDirective := templates.GetWeeklyTemplateDirective()

	// // Create the template processor and initialize the filesystem and directive
	// templateProcessor := templates.NewTemplateProcessor("PrimaryTemplateProcessor", logLevel)
	// templateProcessor.InitializeFilesystem(templates.GetDefaultFilesystemConfiguration())
	// templateProcessor.InitializeDirectives([]templates.DirectiveConfiguration{weeklyDirective})

	// // Complete initialization then perform our single tick
	// templateProcessor.CompleteInitialization()
	// templateProcessor.Tick()
}

func handleAgent(logLevel slog.Level) {

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

	templateProcessor := templates.NewTemplateProcessor("PrimaryTemplateProcessor", logLevel)

	templateProcessor.InitializeFilesystem(filesystemConfiguration)

	templateProcessor.InitializeDirectives(defaultDirectives)

	templateProcessor.CompleteInitialization()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

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

func handleServer(logLevel slog.Level) {
	// Start the web server and get the broadcast channel
	broadcast := server.StartServer()

	// Load filesystem configuration
	filesystemConfiguration, err := templates.GetDefaultFilesystemConfiguration()
	if err != nil {
		panic(err)
	}

	// Get the default template directives
	defaultDirectives := templates.GetDefaultTemplateDirectives()

	// Create and initialize template processor
	templateProcessor := templates.NewTemplateProcessor("PrimaryTemplateProcessor", logLevel)
	templateProcessor.InitializeFilesystem(filesystemConfiguration)
	templateProcessor.InitializeDirectives(defaultDirectives)
	templateProcessor.CompleteInitialization()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	// Run agent loop with time broadcasting
	for {
		currentTimeUs := time.Now().UnixMicro()

		// Broadcast current time through the channel every iteration (10 second sleep below)
		broadcast <- server.Event{Message: fmt.Sprintf("Current time: %s", time.Now().Format(time.RFC3339))}

		templateProcessor.Tick()

		// Print the time delta
		logger.Debug("Time delta",
			slog.Int64("delta", time.Now().UnixMicro()-currentTimeUs))

		time.Sleep(time.Second * 10)
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
	-log-level    set log level (debug, info, warn, error)

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

	// Global log level variable
	var logLevel string

	// Define template subcommand and its flags
	templateCmd := flag.NewFlagSet("template", flag.ExitOnError)

	// TODO - revisit soon and remove dead code

	// var buildTarget string
	// templateCmd.StringVar(&buildTarget, "file", "", "Top level file to run templating on.")
	// var buildVerbose bool
	// templateCmd.BoolVar(&buildVerbose, "verbose", false, "Enable verbose output")
	// var templateFlag1 string
	// templateCmd.StringVar(&templateFlag1, "recurse", "", "Whether to recurse templating to generated files.")

	var builtin string
	templateCmd.StringVar(&builtin, "builtin", "", "Builtin template to use (optional)")
	var once bool
	templateCmd.BoolVar(&once, "once", false, "Run templating for the configured files once and exit")
	templateCmd.StringVar(&logLevel, "log-level", "info", "set log level (debug, info, warn, error)")

	// Define agent subcommand and its flags
	agentCmd := flag.NewFlagSet("agent", flag.ExitOnError)
	agentCmd.StringVar(&logLevel, "log-level", "info", "set log level (debug, info, warn, error)")

	// Define server subcommand and its flags
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	serverCmd.StringVar(&logLevel, "log-level", "info", "set log level (debug, info, warn, error)")

	// Get the subcommand before parsing flags
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
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

		// Configure log level after parsing flags
		level := configureLogLevel(logLevel)

		// Handle builtin argument
		if builtin != "" {
			handleBuiltinTemplate(builtin, level)
		}

		if once {
			fmt.Println("Processing all files once.")
			handleOnce(level)
		}

	case "agent":
		if err := agentCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Configure log level after parsing flags
		level := configureLogLevel(logLevel)

		handleAgent(level)

	case "server":
		if err := serverCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Configure log level after parsing flags
		level := configureLogLevel(logLevel)

		handleServer(level)

	default:
		fmt.Fprintf(os.Stderr, "Invalid subcommand: %s", os.Args[1])
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}
}
