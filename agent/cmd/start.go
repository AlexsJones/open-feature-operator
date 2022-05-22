package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-feature/open-feature-operator/agent/pkg/runtime"
	"github.com/open-feature/open-feature-operator/agent/pkg/service"
	"github.com/open-feature/open-feature-operator/agent/pkg/sync"
	"github.com/spf13/cobra"
)

var (
	serviceProvider    string
	syncProvider       string
	filePath           string
	registeredServices map[string]service.IService = map[string]service.IService{
		"http": &service.HttpService{
			HttpServiceConfiguration: &service.HttpServiceConfiguration{
				Port: int32(8080),
			},
		},
	}
	registeredSync map[string]sync.ISync = map[string]sync.ISync{
		"filepath": &sync.FilePathSync{},
	}
)

func findService(name string) (service.IService, error) {
	if v, ok := registeredServices[name]; !ok {
		return nil, errors.New("no service-provider set")
	} else {
		log.Printf("Using %s service-provider\n", name)
		return v, nil
	}
}

func findSync(name string) (sync.ISync, error) {
	if v, ok := registeredSync[name]; !ok {
		return nil, errors.New("no sync-provider set")
	} else {
		log.Printf("Using %s sync-provider\n", name)
		return v, nil
	}
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the open-feature agent",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Configure service-provider impl------------------------------------------
		var serviceImpl service.IService
		if foundService, err := findService(serviceProvider); err != nil {
			return
		} else {
			serviceImpl = foundService
		}
		// Configure sync-provider impl--------------------------------------------
		var syncImpl sync.ISync
		if foundSync, err := findSync(syncProvider); err != nil {
			return
		} else {
			syncImpl = foundSync
		}

		// Serve ------------------------------------------------------------------
		ctx, cancel := context.WithCancel(context.Background())
		errc := make(chan error)
		go func() {
			errc <- func() error {
				c := make(chan os.Signal)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				return fmt.Errorf("%s", <-c)
			}()
		}()

		go runtime.Start(syncImpl, serviceImpl, ctx)

		err := <-errc
		if err != nil {
			cancel()
			log.Printf(err.Error())
		}
	},
}

func init() {

	startCmd.Flags().StringVarP(&serviceProvider, "service-provider", "s", "http", "Set a serve provider e.g. http or socket")
	startCmd.Flags().StringVarP(&syncProvider, "sync-provider", "y", "filepath", "Set a sync provider e.g. filepath or remote")
	startCmd.Flags().StringVarP(&filePath, "filepath", "f", "", "Set a sync provider filepath to read data from")
	rootCmd.AddCommand(startCmd)

}
