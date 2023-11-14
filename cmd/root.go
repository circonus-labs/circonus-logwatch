// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package cmd

import (
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/circonus-labs/circonus-logwatch/internal/agent"
	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/circonus-labs/circonus-logwatch/internal/config/defaults"
	"github.com/circonus-labs/circonus-logwatch/internal/release"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   release.NAME,
	Short: "A small utility to send metrics extracted from log files to Circonus.",
	Long: `Send metrics extracted from log files to Circonus.

Not all useful metrics can be sent directly to a centralized system for analysis
and alerting. Often, there are valuable metrics sequestered in system and application
log files. These logs are not always in a common, easily parsable format.

Using named regular expressions and templates, this utility offers the ability
to extract these metrics from the logs and send them directly to the Circonus
system using one of several methods.

See https://github.com/circonus-labs/circonus-logwatch/etc for example configurations.
`,
	PersistentPreRunE: initLogging,
	Run: func(cmd *cobra.Command, args []string) {
		//
		// show version and exit
		//
		if viper.GetBool(config.KeyShowVersion) {
			fmt.Printf("%s v%s - commit: %s, date: %s, tag: %s\n", release.NAME, release.VERSION, release.COMMIT, release.DATE, release.TAG)
			return
		}

		//
		// show configuration and exit
		//
		if viper.GetString(config.KeyShowConfig) != "" {
			if err := config.ShowConfig(os.Stdout); err != nil {
				log.Fatal().Err(err).Msg("show-config")
			}
			return
		}

		log.Info().
			Int("pid", os.Getpid()).
			Str("name", release.NAME).
			Str("ver", release.VERSION).Msg("Starting")

		a, err := agent.New()
		if err != nil {
			log.Fatal().Err(err).Msg("Initializing")
		}

		_ = config.StatConfig()

		if err := a.Start(); err != nil {
			log.Fatal().Err(err).Msg("Startup")
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func bindFlagError(flag string, err error) {
	if err != nil {
		log.Fatal().Err(err).Str("flag", flag).Msg("binding flag")
	}
}
func bindEnvError(envVar string, err error) {
	if err != nil {
		log.Fatal().Err(err).Str("var", envVar).Msg("binding env var")
	}
}

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zlog := zerolog.New(zerolog.SyncWriter(os.Stderr)).With().Timestamp().Logger()
	log.Logger = zlog

	stdlog.SetFlags(0)
	stdlog.SetOutput(zlog)

	cobra.OnInitialize(initConfig)

	desc := func(desc, env string) string {
		return fmt.Sprintf("[ENV: %s] %s", env, desc)
	}

	//
	// Basic
	//
	{
		var (
			longOpt     = "config"
			shortOpt    = "c"
			description = "config file (default is " + defaults.EtcPath + "/" + release.NAME + ".(json|toml|yaml)"
		)
		RootCmd.PersistentFlags().StringVarP(&cfgFile, longOpt, shortOpt, "", description)
	}

	{
		const (
			key         = config.KeyLogConfDir
			longOpt     = "log-conf-dir"
			shortOpt    = "l"
			envVar      = release.ENVPREFIX + "_PLUGIN_DIR"
			description = "Log configuration directory"
		)

		RootCmd.Flags().StringP(longOpt, shortOpt, defaults.LogConfPath, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaults.LogConfPath)
	}

	//
	// Destination for metrics
	//
	{
		const (
			key         = config.KeyDestType
			longOpt     = "dest"
			envVar      = release.ENVPREFIX + "_DESTINATION"
			description = "Destination[agent|check|log|statsd] type for metrics"
		)

		RootCmd.Flags().String(longOpt, defaults.DestinationType, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaults.DestinationType)
	}
	{
		const (
			key         = config.KeyDestCfgID
			longOpt     = "dest-id"
			envVar      = release.ENVPREFIX + "_DEST_ID"
			description = "Destination[statsd|agent] metric group ID"
		)

		RootCmd.Flags().String(longOpt, release.NAME, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
	}
	{
		const (
			key         = config.KeyDestCfgCID
			longOpt     = "dest-cid"
			envVar      = release.ENVPREFIX + "_DEST_CID"
			description = "Destination[check] Check ID (not check bundle id)"
		)

		RootCmd.Flags().String(longOpt, "", desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
	}
	{
		const (
			key         = config.KeyDestCfgURL
			longOpt     = "dest-url"
			envVar      = release.ENVPREFIX + "_DEST_URL"
			description = "Destination[check] Check Submission URL"
		)

		RootCmd.Flags().String(longOpt, "", desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
	}
	{
		const (
			key         = config.KeyDestCfgPort
			longOpt     = "dest-port"
			envVar      = release.ENVPREFIX + "_DEST_PORT"
			description = "Destination[agent|statsd] port (agent=2609, statsd=8125)"
		)

		RootCmd.Flags().String(longOpt, "", desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
	}
	{
		const (
			key         = config.KeyDestCfgInstanceID
			longOpt     = "dest-instance-id"
			envVar      = release.ENVPREFIX + "_DEST_INSTANCE_ID"
			description = "Destination[check] Check Instance ID"
		)

		RootCmd.Flags().String(longOpt, "", desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
	}
	{
		const (
			key         = config.KeyDestCfgTarget
			longOpt     = "dest-target"
			envVar      = release.ENVPREFIX + "_DEST_TARGET"
			description = "Destination[check] Check target (default hostname)"
		)

		RootCmd.Flags().String(longOpt, "", desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
	}
	{
		const (
			key         = config.KeyDestCfgSearchTag
			longOpt     = "dest-tag"
			envVar      = release.ENVPREFIX + "_DEST_TAG"
			description = "Destination[check] Check search tag"
		)

		RootCmd.Flags().String(longOpt, "", desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
	}
	{
		const (
			key         = config.KeyDestCfgStatsdPrefix
			longOpt     = "dest-statsd-prefix"
			envVar      = release.ENVPREFIX + "_DEST_STATSD_PREFIX"
			description = "Destination[statsd] Prefix prepended to every metric sent to StatsD"
		)

		RootCmd.Flags().String(longOpt, defaults.StatsdPrefix, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaults.StatsdPrefix)
	}
	{
		const (
			key         = config.KeyDestCfgAgentInterval
			longOpt     = "dest-agent-interval"
			envVar      = release.ENVPREFIX + "_DEST_AGENT_INTERVAL"
			description = "Destination[agent] Interval for metric submission to agent"
		)

		RootCmd.Flags().String(longOpt, defaults.AgentInterval, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaults.AgentInterval)
	}

	//
	// API
	//
	{
		const (
			key          = config.KeyAPITokenKey
			longOpt      = "api-key"
			defaultValue = ""
			envVar       = release.ENVPREFIX + "_API_KEY"
			description  = "Circonus API Token key or 'cosi' to use COSI config"
		)
		RootCmd.Flags().String(longOpt, defaultValue, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
	}

	{
		const (
			key         = config.KeyAPITokenApp
			longOpt     = "api-app"
			envVar      = release.ENVPREFIX + "_API_APP"
			description = "Circonus API Token app"
		)

		RootCmd.Flags().String(longOpt, defaults.APIApp, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaults.APIApp)
	}

	{
		const (
			key         = config.KeyAPIURL
			longOpt     = "api-url"
			envVar      = release.ENVPREFIX + "_API_URL"
			description = "Circonus API URL"
		)

		RootCmd.Flags().String(longOpt, defaults.APIURL, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaults.APIURL)
	}

	{
		const (
			key          = config.KeyAPICAFile
			longOpt      = "api-ca-file"
			defaultValue = ""
			envVar       = release.ENVPREFIX + "_API_CA_FILE"
			description  = "Circonus API CA certificate file"
		)

		RootCmd.Flags().String(longOpt, defaultValue, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
	}

	// Miscellaneous

	{
		const (
			key         = config.KeyDebug
			longOpt     = "debug"
			shortOpt    = "d"
			envVar      = release.ENVPREFIX + "_DEBUG"
			description = "Enable debug messages"
		)

		RootCmd.Flags().BoolP(longOpt, shortOpt, defaults.Debug, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaults.Debug)
	}

	{
		const (
			key          = config.KeyDebugCGM
			longOpt      = "debug-cgm"
			defaultValue = false
			envVar       = release.ENVPREFIX + "_DEBUG_CGM"
			description  = "Enable CGM & API debug messages"
		)

		RootCmd.Flags().Bool(longOpt, defaultValue, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = config.KeyDebugTail
			longOpt      = "debug-tail"
			defaultValue = false
			envVar       = release.ENVPREFIX + "_DEBUG_TAIL"
			description  = "Enable log tailing messages"
		)

		RootCmd.Flags().Bool(longOpt, defaultValue, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = config.KeyDebugMetric
			longOpt      = "debug-metric"
			defaultValue = false
			envVar       = release.ENVPREFIX + "_DEBUG_METRIC"
			description  = "Enable metric rule evaluation tracing debug messages"
		)

		RootCmd.Flags().Bool(longOpt, defaultValue, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key         = config.KeyAppStatPort
			longOpt     = "stat-port"
			envVar      = release.ENVPREFIX + "_STAT_PORT"
			description = "Exposes app stats while running"
		)

		RootCmd.Flags().String(longOpt, defaults.AppStatPort, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaults.AppStatPort)
	}

	{
		const (
			key         = config.KeyLogLevel
			longOpt     = "log-level"
			envVar      = release.ENVPREFIX + "_LOG_LEVEL"
			description = "Log level [(panic|fatal|error|warn|info|debug|disabled)]"
		)

		RootCmd.Flags().String(longOpt, defaults.LogLevel, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaults.LogLevel)
	}

	{
		const (
			key         = config.KeyLogPretty
			longOpt     = "log-pretty"
			envVar      = release.ENVPREFIX + "_LOG_PRETTY"
			description = "Output formatted/colored log lines"
		)

		RootCmd.Flags().Bool(longOpt, defaults.LogPretty, desc(description, envVar))
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
		bindEnvError(key, viper.BindEnv(key, envVar))
		viper.SetDefault(key, defaults.LogPretty)
	}

	{
		const (
			key          = config.KeyShowVersion
			longOpt      = "version"
			shortOpt     = "V"
			defaultValue = false
			description  = "Show version and exit"
		)
		RootCmd.Flags().BoolP(longOpt, shortOpt, defaultValue, description)
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
	}

	{
		const (
			key         = config.KeyShowConfig
			longOpt     = "show-config"
			description = "Show config (json|toml|yaml) and exit"
		)

		RootCmd.Flags().String(longOpt, "", description)
		bindFlagError(key, viper.BindPFlag(key, RootCmd.Flags().Lookup(longOpt)))
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(defaults.EtcPath)
		viper.AddConfigPath(".")
		viper.SetConfigName(release.NAME)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		f := viper.ConfigFileUsed()
		if f != "" {
			log.Fatal().Err(err).Str("config_file", f).Msg("Unable to load config file")
		}
	}
}

// initLogging initializes zerolog.
func initLogging(_ *cobra.Command, _ []string) error {
	//
	// Enable formatted output
	//
	if viper.GetBool(config.KeyLogPretty) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	//
	// Enable debug logging, if requested
	// otherwise, default to info level and set custom level, if specified
	//
	if viper.GetBool(config.KeyDebug) {
		viper.Set(config.KeyLogLevel, "debug")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return nil
	}

	if viper.IsSet(config.KeyLogLevel) {
		level := viper.GetString(config.KeyLogLevel)

		switch level {
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "disabled":
			zerolog.SetGlobalLevel(zerolog.Disabled)
		default:
			return fmt.Errorf("Unknown log level (%s)", level) //nolint:goerr113
		}
	}

	return nil
}
