package api

import "blockchain/pkg/config"

func NewApi(
	cfg *config.Configuration,
) *api {
	a := &api{
		cfg: cfg,
	}
	return a
}
