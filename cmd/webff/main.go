package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"

	"github.com/igolaizola/webcli"
	"github.com/igolaizola/webcli/pkg/webff"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

// Build flags
var version = ""
var commit = ""
var date = ""

func main() {
	// Create signal based context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Launch command
	cmd := newCommand()
	if err := cmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func newCommand() *ffcli.Command {
	fs := flag.NewFlagSet("webff", flag.ExitOnError)

	cmds := []*ffcli.Command{
		newVersionCommand(),
		newRunCommand(),
	}
	port := fs.Int("port", 0, "port number")

	return &ffcli.Command{
		ShortUsage: "webff [flags] <subcommand>",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) == 0 {
				s, err := webff.New(cmds, webcli.WithAppName(fs.Name()), webcli.WithAddress(fmt.Sprintf(":%d", *port)))
				if err != nil {
					return err
				}
				return s.Run(ctx)
			}
			return flag.ErrHelp
		},
		Subcommands: cmds,
	}
}

func newVersionCommand() *ffcli.Command {
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "webff version",
		ShortHelp:  "print version",
		Exec: func(ctx context.Context, args []string) error {
			v := version
			if v == "" {
				if buildInfo, ok := debug.ReadBuildInfo(); ok {
					v = buildInfo.Main.Version
				}
			}
			if v == "" {
				v = "dev"
			}
			versionFields := []string{v}
			if commit != "" {
				versionFields = append(versionFields, commit)
			}
			if date != "" {
				versionFields = append(versionFields, date)
			}
			fmt.Println(strings.Join(versionFields, " "))
			return nil
		},
	}
}

func newRunCommand() *ffcli.Command {
	cmd := "run"
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")
	_ = fs.Duration("max-duration", 0, "duration")
	_ = fs.Int("attempts", 0, "int")
	_ = fs.Bool("debug", false, "bool")
	_ = fs.Float64("price", 0, "float64")

	return &ffcli.Command{
		Name:       cmd,
		ShortUsage: fmt.Sprintf("webff %s [flags] <key> <value data...>", cmd),
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("webff"),
		},
		ShortHelp: fmt.Sprintf("webff %s command", cmd),
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			log.Println("running")
			defer log.Println("finished")
			for i := 0; i < 5; i++ {
				select {
				case <-ctx.Done():
				case <-time.After(1 * time.Second):
				}
				fmt.Println("tick", i)
			}
			return nil
		},
		Subcommands: []*ffcli.Command{
			newSubRunCommand(cmd),
		},
	}
}

func newSubRunCommand(parent string) *ffcli.Command {
	cmd := "subrun"
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")
	_ = fs.Duration("max-duration", 0, "duration")
	_ = fs.Int("attempts", 0, "int")
	_ = fs.Bool("debug", false, "bool")
	_ = fs.Float64("price", 0, "float64")

	return &ffcli.Command{
		Name:       cmd,
		ShortUsage: fmt.Sprintf("webff %s %s [flags] <key> <value data...>", parent, cmd),
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("webff"),
		},
		ShortHelp: fmt.Sprintf("webff %s %s command", parent, cmd),
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			log.Println("running")
			defer log.Println("finished")
			for i := 0; i < 5; i++ {
				select {
				case <-ctx.Done():
				case <-time.After(1 * time.Second):
				}
				fmt.Println("tick", i)
			}
			return nil
		},
	}
}
