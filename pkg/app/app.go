package app

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	progressMsg = color.GreenString("==>")
)

type App struct {
	basename string
	name     string
	desc     string
	runFunc  RunFunc
	silence  bool
	options  CliOptions
	commands []*Command
	args     cobra.PositionalArgs
	cmd      *cobra.Command
}

func NewApp(basename, name string, opts ...Option) *App {
	app := &App{
		basename: basename,
		name:     name,
	}

	for _, opt := range opts {
		opt(app)
	}

	app.buildCommand()

	return app
}

type Option func(app *App)

type RunFunc func(name string) error

func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.runFunc = run
	}
}

func WithSilence(silence bool) Option {
	return func(app *App) {
		app.silence = silence
	}
}

func WithDescription(desc string) Option {
	return func(app *App) {
		app.desc = desc
	}
}

func WithArgs(args cobra.PositionalArgs) Option {
	return func(app *App) {
		app.args = args
	}
}

func WithDefaultArgs() Option {
	return func(app *App) {
		app.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q doesn't take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

func (app *App) buildCommand() {
	cmd := &cobra.Command{
		Use:           formatBaseName(app.basename),
		Short:         app.name,
		Long:          app.desc,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          app.args,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	InitFlags(cmd.Flags())

	if len(app.commands) > 0 {
		for _, command := range app.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(formatBaseName(app.basename)))
	}

	cmd.RunE = app.runCommand

	var namedFlagSets NamedFlagSets
	if app.options != nil {
		namedFlagSets = app.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	// add new global flagset to cmd FlagSet
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	app.cmd = cmd
}

func (app App) runCommand(cmd *cobra.Command, args []string) error {
	if !app.silence {
		fmt.Printf("%s starting %s ...", progressMsg, app.name)
	}

	if app.runFunc != nil {
		return app.runFunc(app.basename)
	}
	return nil
}

// AddCommand adds sub command to the application.
func (app *App) AddCommand(cmd *Command) {
	app.commands = append(app.commands, cmd)
}

// AddCommands adds multiple sub commands to the application.
func (app *App) AddCommands(cmds ...*Command) {
	app.commands = append(app.commands, cmds...)
}

func formatBaseName(basename string) string {
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}

func printWorkDir() {
	wd, _ := os.Getwd()
	// todo
	fmt.Printf("%s working directory: %s", progressMsg, wd)
}
