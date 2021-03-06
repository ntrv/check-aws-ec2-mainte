package checkawsec2mainte

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"

	"github.com/ntrv/check-aws-ec2-mainte/lib/events"
)

// Set variable from -X option
var (
	version   = "indev"                     // git describe --tags
	buildDate = "1970-01-01 09:00:00+09:00" // date --rfc-3339=seconds
)

// Arguments ... Commandline Options
type Arguments struct {
	Region       string        `short:"r" long:"region" env:"AWS_REGION" description:"AWS Region"`
	CritDuration time.Duration `short:"c" long:"critical-duration" default:"72h" description:"Critical while duration"`
	InstanceIds  []string      `short:"i" long:"instance-id" description:"Filter as EC2 Instance Ids"`
	IsAll        bool          `short:"a" long:"all" description:"Fetch all instances events"`
	Version      func()        `short:"v" long:"version" description:"Print Build Information"`
}

// Cli ...
type Cli struct {
	Args    Arguments
	Command string
	Now     time.Time
}

// NewCli ...
func NewCli(args []string) (*Cli, error) {
	opts := Arguments{}
	opts.Version = Usage

	args, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return nil, err
	}

	return &Cli{
		Args:    opts,
		Command: args[0],
		Now:     time.Now(),
	}, nil
}

// Evaluate ...
func (c Cli) Evaluate(evs events.Events) *checkers.Checker {
	if evs.Len() != 0 {
		msg := evs.String()
		event := evs.GetCloseEvent()

		if event.IsTimeOver(c.Now, c.Args.CritDuration) {
			return checkers.Critical(msg)
		}
		return checkers.Warning(msg)
	}

	return checkers.Ok("Not coming EC2 instance events")
}

// Usage ...
func Usage() {
	fmt.Fprintf(
		os.Stderr,
		"Version: %v\nGoVer: %v\nAwsSDKVer: %v\nBuildDate: %v\n",
		version,
		runtime.Version(),
		aws.SDKVersion,
		buildDate,
	)
	os.Exit(1)
}
