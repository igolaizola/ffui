package webcobra

import (
	"github.com/igolaizola/webcli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Parse(cmds []*cobra.Command) []*webcli.Command {
	var wcmds []*webcli.Command
	for _, cmd := range cmds {
		if cmd.Hidden {
			continue
		}
		wcmds = append(wcmds, toCommand(cmd))
	}
	return wcmds
}

func toCommand(c *cobra.Command) *webcli.Command {
	var subs []*webcli.Command
	for _, sub := range c.Commands() {
		subs = append(subs, toCommand(sub))
	}
	return &webcli.Command{
		Fields:      toFields(c.Flags(), c.LocalFlags(), c.PersistentFlags()),
		Name:        c.Name(),
		Description: c.Short + "\n" + c.Long,
		Subcommands: subs,
	}
}

func toFields(fss ...*pflag.FlagSet) []*webcli.Field {
	var fields []*webcli.Field
	lookup := map[string]string{}
	for _, fs := range fss {
		if fs == nil {
			continue
		}
		fs.VisitAll(func(f *pflag.Flag) {
			if _, ok := lookup[f.Name]; ok {
				return
			}
			lookup[f.Name] = f.Name
			fields = append(fields, &webcli.Field{
				Name:        f.Name,
				Default:     f.Value.String(),
				Description: f.Usage,
				Type:        toType(f),
			})
		})
	}
	return fields
}

func toType(f *pflag.Flag) webcli.FieldType {
	switch f.Value.Type() {
	case "bool":
		return webcli.Boolean
	case "duration":
		return webcli.Text
	case "float64", "int", "int64", "uint", "uint64":
		return webcli.Number
	case "string":
		return webcli.Text
	default:
		return webcli.Text
	}
}

type Config struct {
	App      string
	Commands []*cobra.Command
	Address  string
}

func New(cfg *Config) (*webcli.Server, error) {
	webcliConfig := &webcli.Config{
		App:      cfg.App,
		Commands: Parse(cfg.Commands),
		Address:  cfg.Address,
	}
	return webcli.New(webcliConfig)
}