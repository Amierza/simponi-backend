package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Amierza/simponi-backend/config/platforms"
	"github.com/Amierza/simponi-backend/dto"
	"go.uber.org/zap"
)

type (
	IMLService interface {
		PredictTags(ctx context.Context, text string) ([]string, error)
	}

	mlService struct {
		logger *zap.Logger
		config *platforms.MLConfig
		client *http.Client
	}

	mlPredictRequest struct {
		Text string `json:"text"`
	}

	mlPredictResponse struct {
		Text      string             `json:"text"`
		Tags      []string           `json:"tags"`
		Scores    map[string]float64 `json:"scores"`
		Threshold float64            `json:"threshold"`
	}
)

func NewMLService(logger *zap.Logger, config *platforms.MLConfig) *mlService {
	return &mlService{
		logger: logger,
		config: config,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (ms *mlService) PredictTags(ctx context.Context, text string) ([]string, error) {
	endpoint := fmt.Sprintf("%s/predict", ms.config.BaseURL)

	payload, err := json.Marshal(mlPredictRequest{Text: text})
	if err != nil {
		ms.logger.Error("failed marshal ml payload", zap.Error(err))
		return nil, fmt.Errorf("failed marshal payload: %w", dto.ErrInternal)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		ms.logger.Error("failed create ml request", zap.Error(err))
		return nil, fmt.Errorf("failed create request: %w", dto.ErrInternal)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ms.client.Do(req)
	if err != nil {
		ms.logger.Error("failed request ml api", zap.Error(err))
		return nil, fmt.Errorf("failed request ml api: %w", dto.ErrInternal)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ms.logger.Error("failed read ml response", zap.Error(err))
		return nil, fmt.Errorf("failed read response: %w", dto.ErrInternal)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		ms.logger.Error("ml api returned non-2xx status", zap.Int("status", resp.StatusCode), zap.ByteString("response", responseBody))
		return nil, fmt.Errorf("ml api returned status %d: %w", resp.StatusCode, dto.ErrInternal)
	}

	var result mlPredictResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		ms.logger.Error("failed unmarshal ml response", zap.Error(err), zap.ByteString("response", responseBody))
		return nil, fmt.Errorf("failed parse response: %w", dto.ErrInternal)
	}

	ms.logger.Info("success predict tags", zap.Strings("tags", result.Tags))

	return result.Tags, nil
}
