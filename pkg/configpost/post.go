package configpost

import (
	"context"
	"gecko/pkg/config"
	"gecko/pkg/search"
	"gecko/server/controller"
	"gecko/server/service"
)

func PostInit(ctx context.Context) {
	search.Init()
	go service.NewService().Init()
	controller.RunServer(ctx, config.Conf.ServerBind)
}
