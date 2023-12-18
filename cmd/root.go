/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/muhlemmer/zitadel-data-loader/internal/client"
	"github.com/muhlemmer/zitadel-data-loader/internal/config"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	background context.Context
	stop       context.CancelFunc
	clientConn *grpc.ClientConn

	cfgFile string

	//go:embed defaults.yaml
	defaultConfig []byte
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zitadel-data-loader",
	Short: "Create and sent random data to a Zitadel instance.",
	Long: `A gRPC client to the Zitadal API that sends pseudo-random requests using GoFakit.
	The basic commands performs a single health check to validate the config and client connection.
	Actual generators are in the sub-commands`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if config.Global.PAT != "" {
			md := make(metadata.MD)
			md.Set("Authorization", fmt.Sprintf("Bearer %s", config.Global.PAT))
			if config.Global.OrgID != "" {
				md.Set("x-zitadel-orgid", config.Global.OrgID)
			}
			background = metadata.NewOutgoingContext(background, md)
		}

		clientConn, err = client.Dial(background)
		return err
	},

	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		stop()
		return clientConn.Close()
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		mgmt := management.NewManagementServiceClient(clientConn)
		ctx, cancel := client.TimeoutCTX(background)
		defer cancel()
		_, err := mgmt.Healthz(ctx, &management.HealthzRequest{})
		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	background, stop = signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	out := zerolog.NewConsoleWriter()

	logger := zerolog.New(out).With().
		Caller().
		Timestamp().
		Logger()
	background = logger.WithContext(background)

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zitadel-data-loader.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		zerolog.Ctx(background).Fatal().Err(err).Msg("read default config")
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".zitadel-data-loader" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".zitadel-data-loader")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err = viper.MergeInConfig()
	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		zerolog.Ctx(background).Info().Err(err).Msg("config file load")
	} else if err != nil {
		zerolog.Ctx(background).Fatal().Err(err).Msg("config file load")
	}
	if err := viper.Unmarshal(&config.Global); err != nil {
		zerolog.Ctx(background).Fatal().Err(err).Msg("unmarshal config to struct")
	}
	zerolog.Ctx(background).Debug().Func(func(e *zerolog.Event) {
		e.Interface("config", config.Global).Msg("")
	})
}
