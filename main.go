package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohitkumar/finch/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cfg struct {
	rest.Config
}
type cli struct {
	cfg cfg
}

func setupFlags(cmd *cobra.Command) error {
	cmd.Flags().String("config-file", "", "Path to config file.")
	cmd.Flags().String("redis-host", "localhost", "redis host name")
	cmd.Flags().Int("redis-port", 6379, "redis port")
	cmd.Flags().String("namespace", "finch", "namespace used in storage")
	cmd.Flags().Int("http-port", 8080, "htt port for rest endpoints")
	cmd.Flags().Int("grpc-port", 8099, "grpc port for worker connection")
	return viper.BindPFlags(cmd.Flags())
}

func (c *cli) setupConfig(cmd *cobra.Command, args []string) error {
	var err error

	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return err
	}
	viper.SetConfigFile(configFile)

	if err = viper.ReadInConfig(); err != nil {
		// it's ok if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	c.cfg.RedisConfig.Host = viper.GetString("redis-host")
	c.cfg.RedisConfig.Port = viper.GetInt("redis-port")
	c.cfg.RedisConfig.Namespace = viper.GetString("namespace")
	c.cfg.Port = viper.GetInt("http-port")

	return nil
}

func (c *cli) run(cmd *cobra.Command, args []string) error {
	var err error
	server, err := rest.NewServer(c.cfg.Config)
	if err != nil {
		panic(err)
	}
	err = server.Start()

	if err != nil {
		fmt.Println("could not start server")
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	return server.Stop()
}

func main() {
	cli := &cli{}

	cmd := &cobra.Command{
		Use:     "finch",
		PreRunE: cli.setupConfig,
		RunE:    cli.run,
	}

	if err := setupFlags(cmd); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
