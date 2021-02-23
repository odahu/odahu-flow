package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"runtime/pprof"
)

var (
	cfgFile string
	cfg config.ToolsConfig
	profilingFile *os.File
)

var rootCmd = &cobra.Command{
	Use:   "odahu-tools",
	Short: "odahu-tools is a simple command line tool that provides API to the set of ODAHU platform features",
	Long: `odahu-tools provides API to execute the same logic that is used by the ODAHU platform in the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cpuprofile != "" {
			var err error
			profilingFile, err = os.Create(cpuprofile)
			if err != nil {
				zap.S().Fatal("could not create CPU profile: ", err)
			}
			if err := pprof.StartCPUProfile(profilingFile); err != nil {
				zap.S().Fatal("could not start CPU profile: ", err)
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if cpuprofile != "" {
			pprof.StopCPUProfile()
			if profilingFile != nil {
				if err := profilingFile.Close(); err != nil {
					zap.S().Fatal(err)
				}
			}
		}

	},
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Unable to initialize logger")
	}
	defer func() {
		_ = logger.Sync()
	}()

	zap.ReplaceGlobals(logger)
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			zap.S().Errorw("Unable to get home directory", zap.Error(err))
			panic("Unable to get home directory")
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".odahu-tools")
	}

	if err := viper.ReadInConfig(); err != nil {
		zap.S().Errorw("Can't read config", zap.Error(err))
		panic("Can't read config")
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		zap.S().Errorw("Unable to unmarshall config", zap.Error(err))
		panic(err)
	}

}

func init() {
	initLogger()
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.odahu-tools.yaml)",
	)
	rootCmd.PersistentFlags().StringVar(
		&cpuprofile, "cpuprofile", "",
		"if specified as TARGET cpu profiling results will be saved to TARGET file",
	)
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		zap.S().Errorw("Error during executing the root command", zap.Error(err))
		os.Exit(1)
	}
}