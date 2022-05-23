package runtime

import (
	"context"
	"time"

	"github.com/open-feature/open-feature-operator/agent/pkg/service"
	"github.com/open-feature/open-feature-operator/agent/pkg/sync"
)

func Start(syncr sync.ISync, server service.IService, ctx context.Context) {

	// This is a very simple example of how the interface can be used for service and sync
	// The service interface will serve requests whilst the sync interface is responsible
	// for refreshing the configuration data
	go server.Serve(

		func(ir service.IServiceRequest) service.IServiceResponse {

			if ir.GetRequestType() == service.SERVICE_REQUEST_ALL_FLAGS {

				return ir.GenerateServiceResponse("{}")
			}
			return nil
		})

	for {

		time.Sleep(time.Second * 10)
	}
}
