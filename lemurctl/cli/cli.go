package cli

import (
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/lemur/lemurctl/utils"
	"github.com/turbonomic/lemur/lemurctl/version"
	"github.com/urfave/cli"
	"os"
	"path"
	"time"
)

var (
	usage       = "lemurctl controls Lemur, powered by Turbonomic."
	description = "lemurctl is a command line utility to access and control Lemur " +
		"powered by Turbonomic. It offers the capability to view real time entities " +
		"discovered by Lemur, and the resource consumer and provider relationships between " +
		"these entities (i.e., the supply chain). Using Lemur, user can easily sort " +
		"entities of certain type (e.g., applications, or containers) based on their " +
		"resource consumption. User can also take a top-down approach to examine the metric " +
		"and resource usage of a specific application and all the other entities along the " +
		"supply chain of that application to quickly determine the resource bottle neck that " +
		"may affect the performance of the application." +
		"\n\nlemurctl interacts with influxdb to retrieve the entity information and build " +
		"the supply chain. Make sure the influxdb service in your Lemur install can be " +
		"accessed from the host where lemurctl is running, and specify the correct " +
		"$INFLUXDB_SERVER environment variable, or the --influxdb command line option."
)

func Run() {
	app := &cli.App{}
	app.Name = path.Base(os.Args[0])
	app.Usage = usage
	app.Description = description
	app.Version = version.Version + " (" + version.GitCommit + ")"
	if buildTime, err := time.Parse(time.RFC1123Z, version.BuildTime); err == nil {
		app.Compiled = buildTime
	} else {
		app.Compiled = time.Time{}
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "influxdb",
			Value:  utils.GetLocalIP() + ":8086",
			Usage:  "specify the endpoint of the InfluxDB server",
			EnvVar: "INFLUXDB_SERVER",
		},
		cli.BoolFlag{
			Name:   "debug,d",
			Usage:  "enable debug mode",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:   "log-level,l",
			Value:  "info",
			Usage:  "specify log level (debug, info, warn, error, fatal, panic)",
			EnvVar: "LOG_LEVEL",
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.SetLevel(level)
		// If a log level wasn't specified and we are running in debug mode,
		// enforce log-level=debug.
		if !c.IsSet("log-level") && !c.IsSet("l") &&
			(c.Bool("debug") || c.Bool("d")) {
			log.SetLevel(log.DebugLevel)
		}
		log.SetFormatter(
			&log.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: timeFormat,
			})
		return nil
	}

	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
