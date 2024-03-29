package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/mohitkumar/finch/coordinator"
	"github.com/mohitkumar/finch/httpserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	loggerMgr := initZapLog()
	zap.ReplaceGlobals(loggerMgr)
	defer loggerMgr.Sync()

	cli := &cli{}

	// START_HIGHLIGHT
	cmd := &cobra.Command{
		Use:     "proglog",
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

type cli struct {
	cfg cfg
}

type cfg struct {
	coordinator.Config
	startHttp bool
	httpPort  int
}

func setupFlags(cmd *cobra.Command) error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Flags().String("config-file", "", "Path to config file.")

	dataDir := path.Join(os.TempDir(), "coord")
	cmd.Flags().String("data-dir",
		dataDir,
		"Directory to store coordinator data and Raft data.")
	cmd.Flags().String("node-name", hostname, "Unique server ID.")

	cmd.Flags().String("bind-addr",
		"127.0.0.1:8401",
		"Address to bind Serf on.")
	cmd.Flags().Int("rpc-port",
		8400,
		"Port for RPC clients (and Raft) connections.")
	cmd.Flags().Int("http-port",
		8000,
		"Port for http request.")
	cmd.Flags().StringSlice("start-join-addrs",
		nil,
		"Serf addresses to join.")
	cmd.Flags().Bool("bootstrap", false, "Bootstrap the cluster.")
	cmd.Flags().Bool("startHttp", false, "whehter to start http server")

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

	c.cfg.Dir = viper.GetString("data-dir")
	c.cfg.NodeName = viper.GetString("node-name")
	c.cfg.BindAddr = viper.GetString("bind-addr")
	c.cfg.RPCPort = viper.GetInt("rpc-port")
	c.cfg.StartJoinAddrs = viper.GetStringSlice("start-join-addrs")
	c.cfg.Bootstrap = viper.GetBool("bootstrap")
	c.cfg.startHttp = viper.GetBool("startHttp")
	c.cfg.httpPort = viper.GetInt("http-port")
	return nil
}

func (c *cli) run(cmd *cobra.Command, args []string) error {
	var err error
	coord, err := coordinator.New(c.cfg.Config)
	if err != nil {
		return err
	}
	if c.cfg.startHttp {
		cfg := &httpserver.Config{
			Port:         c.cfg.httpPort,
			CoordRpcPort: c.cfg.RPCPort,
		}
		http, err := httpserver.NewServer(*cfg)
		if err != nil {
			return err
		}
		err = http.Start()
		if err != nil {
			return err
		}
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	return coord.Shutdown()
}

func initZapLog() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	return logger
}
