package services

import (
	"context"
	"github.com/samber/do"
	"go.uber.org/zap"
	"ltp-api/internal/config"
	"ltp-api/internal/dto"
	"ltp-api/pkg/kraken"
	"strings"
)

type KrakenService interface {
	LTP(ctx context.Context) (*dto.LTPResponse, error)
}

type krakenService struct {
	client kraken.Client
	pairs  []string
}

func NewKrakenService(i *do.Injector) (KrakenService, error) {
	cfg := do.MustInvoke[*config.Config](i)
	return &krakenService{
		client: do.MustInvoke[kraken.Client](i),
		pairs:  cfg.Pairs,
	}, nil

}

func (s *krakenService) LTP(ctx context.Context) (*dto.LTPResponse, error) {

	res := &dto.LTPResponse{Ltp: make([]*dto.LTP, 0)}

	for _, pair := range s.pairs {
		data, err := s.client.GetTicker(ctx, strings.ReplaceAll(pair, "/", ""))
		if err != nil {
			zap.S().Error(err)
			continue
		}

		amount := ""

		for _, d := range data.Result {
			if len(d.C[0]) > 0 {
				amount = d.C[0]
				break
			}
		}

		if len(amount) > 0 {
			res.Ltp = append(res.Ltp, &dto.LTP{
				Pair:   pair,
				Amount: amount,
			})
		}
	}

	return res, nil
}
