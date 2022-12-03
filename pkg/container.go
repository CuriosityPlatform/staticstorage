package pkg

import (
	"staticstorage/pkg/assetcache"
	"staticstorage/pkg/service"
)

func Service() *service.Service {
	return service.NewService(assetcache.Warmer{})
}
