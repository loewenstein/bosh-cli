package director

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type RuntimeConfig struct {
	Properties string
}

type RuntimeConfigDiffResponse struct {
	Diff [][]interface{} `json:"diff"`
}

type RuntimeConfigDiff struct {
	Diff [][]interface{}
}

func NewRuntimeConfigDiff(diff [][]interface{}) RuntimeConfigDiff {
	return RuntimeConfigDiff{
		Diff: diff,
	}
}

func (d DirectorImpl) LatestRuntimeConfig() (RuntimeConfig, error) {
	resps, err := d.client.RuntimeConfigs()
	if err != nil {
		return RuntimeConfig{}, err
	}

	if len(resps) == 0 {
		return RuntimeConfig{}, bosherr.Error("No runtime config")
	}

	return resps[0], nil
}

func (d DirectorImpl) UpdateRuntimeConfig(manifest []byte) error {
	return d.client.UpdateRuntimeConfig(manifest)
}

func (d DirectorImpl) DiffRuntimeConfig(manifest []byte) (RuntimeConfigDiff, error) {
	resps, err := d.client.DiffRuntimeConfig(manifest)
	if err != nil {
		return RuntimeConfigDiff{}, err
	}

	return NewRuntimeConfigDiff(resps.Diff), nil
}

func (c Client) RuntimeConfigs() ([]RuntimeConfig, error) {
	var resps []RuntimeConfig

	err := c.clientRequest.Get("/runtime_configs?limit=1", &resps)
	if err != nil {
		return resps, bosherr.WrapErrorf(err, "Finding runtime configs")
	}

	return resps, nil
}

func (c Client) UpdateRuntimeConfig(manifest []byte) error {
	path := "/runtime_configs"

	setHeaders := func(req *http.Request) {
		req.Header.Add("Content-Type", "text/yaml")
	}

	_, _, err := c.clientRequest.RawPost(path, manifest, setHeaders)
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating runtime config")
	}

	return nil
}

func (c Client) DiffRuntimeConfig(manifest []byte) (RuntimeConfigDiffResponse, error) {
	var resps RuntimeConfigDiffResponse

	path := "/diff"

	setHeaders := func(req *http.Request) {
		req.Header.Add("Content-Type", "text/yaml")
	}

	err := c.clientRequest.Post(path, manifest, setHeaders, &resps)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Calculation runtime config diff")
	}

	return resps, nil
}
