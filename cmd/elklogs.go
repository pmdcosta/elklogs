package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/pmdcosta/elklogs/internal/domain"
	"github.com/pmdcosta/elklogs/internal/elasticconn"
	"github.com/pmdcosta/elklogs/internal/tail"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "elklogs",
	Short:   "elklogs query and tail ELK logs from the terminal",
	Version: version,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least 1 argument.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		run(args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// rootConfig holds the application global config
var rootConfig struct {
	logger *logrus.Entry
	debug  bool
}

// logsConfig holds the configs for the logs cmd
var logsConfig struct {
	// auth
	user     string
	password string

	// query
	after          string
	before         string
	indexPattern   string
	query          string
	format         string
	timestampField string

	// behavior
	reverse  bool
	follow   bool
	entries  int
	refresh  time.Duration
	showTime bool
}

func init() {
	cobra.OnInitialize(initLogger)

	// persistent flags
	rootCmd.PersistentFlags().BoolVar(&rootConfig.debug, "debug", false, "Enable debug logs")

	// behavior flags
	rootCmd.Flags().BoolVarP(&logsConfig.follow, "follow", "f", false, "Follow log output")
	rootCmd.Flags().IntVarP(&logsConfig.entries, "entries", "n", 50, "Number of lines to show from the end of the logs")

	// connection flags
	rootCmd.Flags().StringVarP(&logsConfig.user, "user", "u", "", "Elastic search basic auth user")
	rootCmd.Flags().StringVarP(&logsConfig.password, "password", "p", "", "Elastic search basic auth password")

	// query flags
	rootCmd.Flags().StringVarP(&logsConfig.after, "after", "a", "", `Get logs after specified date (example: -a "2016-06-17T15:00")`)
	rootCmd.Flags().StringVarP(&logsConfig.before, "before", "b", "", `Get logs before specified date (example: -a "2016-06-17T15:00")`)
	rootCmd.Flags().StringVar(&logsConfig.indexPattern, "index-pattern", "logstash-[0-9].*", "Only log indices that match the pattern will be retrieved")
	rootCmd.Flags().BoolVarP(&logsConfig.reverse, "reverse", "r", false, "Show the newest entries first")
	rootCmd.Flags().StringVarP(&logsConfig.query, "query", "q", "", `Elastic query string search (example: -q "host:myhost.example.com AND level:error")`)
	rootCmd.Flags().DurationVar(&logsConfig.refresh, "refresh", 1*time.Second, `Refresh interval (example: --refresh 1s)`)
	rootCmd.Flags().StringVarP(&logsConfig.format, "output", "o", "", `Output format (example: -o "%timestamp: %log")`)
	rootCmd.Flags().StringVar(&logsConfig.timestampField, "timestamp-field", "@timestamp", `Timestamp field name in the database`)
	rootCmd.Flags().BoolVarP(&logsConfig.showTime, "timestamp", "t", false, "Show timestamp before the log")
}

// initLogger sets up the application logger
func initLogger() {
	rootConfig.logger = logrus.WithFields(nil)
	if rootConfig.debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func run(args []string) {
	// parse query time filters
	var after *time.Time
	var before *time.Time
	if logsConfig.after != "" {
		a, err := time.Parse("2006-01-02T15:04", logsConfig.after)
		if err != nil {
			rootConfig.logger.WithFields(logrus.Fields{"err": err, "date": logsConfig.after}).Fatal("invalid after date")
		}
		after = &a
	}
	if logsConfig.before != "" {
		b, err := time.Parse("2006-01-02T15:04", logsConfig.before)
		if err != nil {
			rootConfig.logger.WithFields(logrus.Fields{"err": err, "date": logsConfig.before}).Fatal("invalid before date")
		}
		before = &b
	}

	// parse output format
	fields := tail.GetFields(logsConfig.format)
	if len(fields) == 0 && logsConfig.format != "" {
		rootConfig.logger.WithFields(logrus.Fields{"format": logsConfig.format}).Fatal("invalid output format")
	}

	// set tailing mode
	if !logsConfig.follow {
		logsConfig.refresh = 0
	}

	q := &domain.Query{
		IndexPattern:   logsConfig.indexPattern,
		AfterDateTime:  after,
		BeforeDateTime: before,
		Reverse:        logsConfig.reverse,
		Query:          logsConfig.query,
		Refresh:        logsConfig.refresh,
		Entries:        logsConfig.entries,
		Format:         logsConfig.format,
		FormatFields:   fields,
		TimestampField: logsConfig.timestampField,
		ShowTime:       logsConfig.showTime,
	}

	// create elastic connector
	c, err := elasticconn.New(args[0], logsConfig.user, logsConfig.password)
	if err != nil {
		rootConfig.logger.WithFields(logrus.Fields{"err": err, "url": args[0]}).Fatal("failed to connect to elastic cluster")
	}

	// create tail
	t := tail.New(rootConfig.logger, c)

	// start tailing logs
	if err := t.Start(q); err != nil {
		rootConfig.logger.WithFields(logrus.Fields{"err": err, "url": args[0]}).Fatal("failed to connect to elastic cluster")
	}

}
