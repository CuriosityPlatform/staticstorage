package assetcache

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/pkg/errors"

	"staticstorage/pkg/config"
)

type Warmer struct{}

func (w Warmer) WarmAssetsCache(ctx context.Context, assets []config.ExternalAsset, assetFolder string) error {
	for _, asset := range assets {
		err := w.warmAsset(ctx, asset, assetFolder)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w Warmer) warmAsset(ctx context.Context, asset config.ExternalAsset, assetFolder string) error {
	_, err := url.Parse(asset.URL)
	if err != nil {
		return errors.Wrapf(err, "failed to parse asset %s url", asset.Name)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", asset.URL, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "failed to download asset by %s", asset.URL)
	}
	defer resp.Body.Close()

	out, err := os.Create(path.Join(assetFolder, asset.Name))
	if err != nil {
		return errors.Wrapf(err, "failed to create file for asset %s", asset.Name)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return errors.Wrapf(err, "failed to store asset %s", asset.Name)
}
