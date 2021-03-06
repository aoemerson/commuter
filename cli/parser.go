package cli

import (
	"flag"

	"github.com/KyleBanks/commuter/cmd"
	"github.com/KyleBanks/commuter/pkg/geo"
)

// ArgParser parses input arguments from the command line.
type ArgParser struct {
	Args []string
}

// NewArgParser initializes and returns an ArgParser.
func NewArgParser(args []string) *ArgParser {
	return &ArgParser{
		Args: args,
	}
}

// Parse attempts to determine which command is being executed,
// parse its flags, and return it.
func (a *ArgParser) Parse(conf *cmd.Configuration, s cmd.StorageProvider) (cmd.RunnerValidator, error) {
	if conf == nil || len(a.Args) == 0 {
		return a.parseConfigureCmd(s)
	}

	switch a.Args[0] {
	case cmdAdd:
		return a.parseAddCmd(s, a.Args[1:])
	case cmdList:
		return a.parseListCmd(s, a.Args[1:])
	}

	return a.parseCommuteCmd(conf, a.Args)
}

// parseConfigureCmd parses and returns a ConfigureCmd.
func (a *ArgParser) parseConfigureCmd(s cmd.StorageProvider) (*cmd.ConfigureCmd, error) {
	return &cmd.ConfigureCmd{
		Input: NewStdin(),
		Store: s,
	}, nil
}

// parseCommuteCmd parses and returns a CommuteCmd from user supplied flags.
func (a *ArgParser) parseCommuteCmd(conf *cmd.Configuration, args []string) (*cmd.CommuteCmd, error) {
	r, err := geo.NewRouter(conf.APIKey)
	if err != nil {
		return nil, err
	}

	c := cmd.CommuteCmd{Durationer: r, Locator: r}

	f := flag.NewFlagSet(cmdCommute, flag.ExitOnError)
	f.StringVar(&c.From, commuteFromParam, cmd.DefaultLocationAlias, commuteFromUsage)
	f.BoolVar(&c.FromCurrent, commuteFromCurrentParam, false, commuteFromCurrentUsage)
	f.StringVar(&c.To, commuteToParam, cmd.DefaultLocationAlias, commuteToUsage)
	f.BoolVar(&c.ToCurrent, commuteToCurrentParam, false, commuteToCurrentUsage)

	f.BoolVar(&c.Drive, commuteDriveParam, false, commuteDriveUsage)
	f.BoolVar(&c.Walk, commuteWalkParam, false, commuteWalkUsage)
	f.BoolVar(&c.Bike, commuteBikeParam, false, commuteBikeUsage)
	f.BoolVar(&c.Transit, commuteTransitParam, false, commuteTransitUsage)
	f.Parse(args)

	// Only default drive to true when no other methods of transport are selected.
	if !c.Walk && !c.Bike && !c.Transit {
		c.Drive = true
	}

	return &c, nil
}

// parseAddCmd parses and returns an AddCmd from user supplied flags.
func (a *ArgParser) parseAddCmd(s cmd.StorageProvider, args []string) (*cmd.AddCmd, error) {
	c := cmd.AddCmd{Store: s}

	f := flag.NewFlagSet(cmdAdd, flag.ExitOnError)
	f.StringVar(&c.Name, addNameParam, "", addNameUsage)
	f.StringVar(&c.Value, addLocationParam, "", addLocationUsage)
	f.Parse(args)

	return &c, nil
}

// parseListCmd parses and returns a ListCmd.
func (a *ArgParser) parseListCmd(s cmd.StorageProvider, args []string) (*cmd.ListCmd, error) {
	return &cmd.ListCmd{}, nil
}
