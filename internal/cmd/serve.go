package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rdoorn/ixxi/internal/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serveCmd() *cobra.Command {
	// Serve
	command := &cobra.Command{
		Use:   "serve",
		Short: "start the mercury loadbalancer",
		Run:   serve(),
	}
	// Serve Flags
	command.PersistentFlags().String("pid-file", "/var/run/mercury.pid", "location of the pid file")
	viper.BindPFlag("pid_file", command.PersistentFlags().Lookup("pid-file"))
	command.PersistentFlags().String("port", "80", "port to listen on")
	viper.BindPFlag("port", command.PersistentFlags().Lookup("port"))
	command.PersistentFlags().String("listener", "127.0.0.1", "ip address to listen on")
	viper.BindPFlag("listener", command.PersistentFlags().Lookup("listener"))
	return command
}

func serve() func(command *cobra.Command, args []string) {
	return func(command *cobra.Command, args []string) {

		var config handler.Config
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("failed to unmarshal config: %s\n", err)
			os.Exit(255)
		}

		if err := config.Verify(); err != nil {
			fmt.Printf("failed to verify config: %s\n", err)
			os.Exit(255)
		}

		// Start the application
		ixxi := handler.New()
		if err := ixxi.Start(&config); err != nil {
			fmt.Printf("failed startup: %s\n", err)
			os.Exit(255)
		}

		// wait for sigint or sigterm for cleanup - note that sigterm cannot be caught
		sigterm := make(chan os.Signal, 10)
		signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)

		for {
			select {
			case <-sigterm:
				ixxi.Warnf("Program killed by signal!")
				ixxi.Stop()
				return
			}
		}

	}
}
