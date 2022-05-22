package runtime

import (
	"context"
	"log"

	"github.com/open-feature/open-feature-operator/agent/pkg/service"
	"github.com/open-feature/open-feature-operator/agent/pkg/sync"
)

func Start(syncr sync.ISync, server service.IService, ctx context.Context) {
	log.Println("Starting run loop")

	server.Serve(

		func(ir service.IServiceRequest) service.IServiceResponse {

			if ir.GetRequestType() == service.SERVICE_REQUEST_ALL_FLAGS {

				return ir.GenerateServiceResponse("{}")
			}
			return nil
		})

}
