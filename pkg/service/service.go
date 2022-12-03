package service

import (
	"context"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/UsingCoding/fpgo/pkg/slices"

	"staticstorage/pkg/assetcache"
	"staticstorage/pkg/config"
)

func NewService(warmer assetcache.Warmer) *Service {
	return &Service{warmer: warmer}
}

type Service struct {
	warmer assetcache.Warmer
}

func (s Service) WarmUpStorage(ctx context.Context, c config.Config) error {
	err := os.MkdirAll(c.Cache, os.ModeDir)
	if err != nil {
		return errors.Wrapf(err, "failed to create cache folder %s", c.Cache)
	}

	return s.warmer.WarmAssetsCache(ctx, c.ExternalAssets, c.Cache)
}

func (s Service) CreateHandlers(ctx context.Context, c config.Config) ([]Handler, error) {
	if exists, err := pathExists(c.Cache); !exists || err != nil {
		return nil, errors.Wrapf(err, "cache folder %s does not exists: did you warm up cache", c.Cache)
	}

	return slices.MapErr(c.Handlers, func(h config.Handler) (Handler, error) {
		asset := findAsset(c.ExternalAssets, h.Asset)
		if asset == nil {
			return Handler{}, errors.Errorf("asset %s not found for %s handler", h.Asset, h.Path)
		}

		assetPath := path.Join(c.Cache, asset.Name)
		if exists, err := pathExists(assetPath); !exists || err != nil {
			return Handler{}, errors.Errorf("asset %s not found for %s handler in cache folder: did you warm up cache", h.Asset, h.Path)
		}

		return Handler{
			Path:      h.Path,
			AssetPath: assetPath,
		}, nil
	})
}

type Handler struct {
	Path      string
	AssetPath string
}

func findAsset(assets []config.ExternalAsset, assetName string) *config.ExternalAsset {
	for _, asset := range assets {
		if asset.Name == assetName {
			return &asset
		}
	}
	return nil
}

func pathExists(p string) (bool, error) {
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
