package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Amierza/simponi-backend/config/platforms"
	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/helper"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	IShopeeService interface {
		GenerateAuthURL(ctx context.Context, storeID, platformID *uuid.UUID) (string, error)
		GetAccessToken(ctx context.Context, code, shopID string) (*dto.ShopeeAccessTokenResponse, error)
		RefreshAccessToken(ctx context.Context, refreshToken, shopID string) (*dto.ShopeeRefreshTokenResponse, error)

		// product related
		SyncProducts(ctx context.Context, storePlatform *entity.StorePlatform) error
		GetItemList(ctx context.Context, accessToken string, shopID string) (*dto.ShopeeGetItemListResponse, error)
		GetItemBaseInfo(ctx context.Context, accessToken string, shopID string, itemIDs []int64) (*dto.ShopeeGetItemBaseInfoResponse, error)
	}

	shopeeService struct {
		tx                  repository.ITransaction
		storePlatformRepo   repository.IStorePlatformRepository
		storeCredentialRepo repository.IStoreCredentialRepository
		productRepo         repository.IProductRepository
		externalProductRepo repository.IExternalProductRepository
		logger              *zap.Logger
		config              *platforms.ShopeeConfig
		redis               *redis.Client
		client              *http.Client
	}
)

func NewShopeeService(
	tx repository.ITransaction,
	storePlatformRepo repository.IStorePlatformRepository,
	storeCredentialRepo repository.IStoreCredentialRepository,
	productRepo repository.IProductRepository,
	externalProductRepo repository.IExternalProductRepository,
	logger *zap.Logger,
	config *platforms.ShopeeConfig,
	redis *redis.Client,
) *shopeeService {
	return &shopeeService{
		tx:                  tx,
		storePlatformRepo:   storePlatformRepo,
		storeCredentialRepo: storeCredentialRepo,
		productRepo:         productRepo,
		externalProductRepo: externalProductRepo,
		logger:              logger,
		config:              config,
		redis:               redis,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (ss *shopeeService) doRequest(ctx context.Context, method string, path string, query url.Values, body interface{}, result interface{}) error {
	timestamp := time.Now().Unix()
	sign := helper.GenerateShopeeSign(
		ss.config.PartnerID,
		path,
		timestamp,
		ss.config.PartnerKey,
	)

	if query == nil {
		query = url.Values{}
	}
	query.Set("partner_id", ss.config.PartnerID)
	query.Set("timestamp", fmt.Sprintf("%d", timestamp))
	query.Set("sign", sign)

	endpoint := fmt.Sprintf(
		"%s%s?%s",
		ss.config.BaseURL,
		path,
		query.Encode(),
	)

	var requestBody io.Reader
	if body != nil {
		jsonPayload, err := json.Marshal(body)
		if err != nil {
			ss.logger.Error("failed marshal shopee payload", zap.Error(err))
			return fmt.Errorf("failed marshal payload: %w", dto.ErrInternal)
		}

		requestBody = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		endpoint,
		requestBody,
	)
	if err != nil {
		ss.logger.Error("failed create shopee request", zap.Error(err))
		return fmt.Errorf("failed create request: %w", dto.ErrInternal)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := ss.client.Do(req)
	if err != nil {
		ss.logger.Error("failed request shopee api", zap.Error(err))
		return fmt.Errorf("failed request shopee api: %w", dto.ErrInternal)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ss.logger.Error("failed read shopee response", zap.Error(err))
		return fmt.Errorf("failed read response: %w", dto.ErrInternal)
	}

	if err := json.Unmarshal(responseBody, result); err != nil {
		ss.logger.Error("failed unmarshal shopee response", zap.Error(err), zap.ByteString("response", responseBody))
		return fmt.Errorf("failed parse response: %w", dto.ErrInternal)
	}

	// validate shopee error response
	if baseResp, ok := result.(interface {
		GetError() string
		GetMessage() string
	}); ok {
		if baseResp.GetError() != "" {
			ss.logger.Error("shopee api returned error", zap.String("error", baseResp.GetError()), zap.String("message", baseResp.GetMessage()))
			return fmt.Errorf("shopee api error: %s", baseResp.GetMessage())
		}
	}

	return nil
}

func (ss *shopeeService) GenerateAuthURL(ctx context.Context, storeID, platformID *uuid.UUID) (string, error) {
	path := "/api/v2/shop/auth_partner"
	timestamp := time.Now().Unix()

	sign := helper.GenerateShopeeSign(
		ss.config.PartnerID,
		path,
		timestamp,
		ss.config.PartnerKey,
	)

	// generate oauth state
	state := uuid.NewString()

	// payload yang disimpan ke redis
	payload := dto.ShopeeOAuthState{
		StoreID:    storeID.String(),
		PlatformID: platformID.String(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		ss.logger.Error("failed marshal oauth state", zap.Error(err))
		return "", err
	}

	// simpan ke redis
	redisKey := fmt.Sprintf("oauth:shopee:state:%s", state)
	err = ss.redis.Set(
		ctx,
		redisKey,
		jsonData,
		10*time.Minute,
	).Err()
	if err != nil {
		ss.logger.Error("failed save oauth state to redis", zap.Error(err))
		return "", err
	}

	redirect := url.QueryEscape(ss.config.RedirectURL)
	authURL := fmt.Sprintf(
		"%s%s?partner_id=%s&timestamp=%d&sign=%s&redirect=%s&state=%s",
		ss.config.BaseURL,
		path,
		ss.config.PartnerID,
		timestamp,
		sign,
		redirect,
		state,
	)

	return authURL, nil
}

func (ss *shopeeService) GetAccessToken(ctx context.Context, code, shopID string) (*dto.ShopeeAccessTokenResponse, error) {
	path := "/api/v2/auth/token/get"
	payload := map[string]interface{}{
		"code":       code,
		"shop_id":    shopID,
		"partner_id": ss.config.PartnerID,
	}

	var result dto.ShopeeAccessTokenResponse
	err := ss.doRequest(
		ctx,
		http.MethodPost,
		path,
		nil,
		payload,
		&result,
	)
	if err != nil {
		ss.logger.Error("failed get shopee access token", zap.Error(err), zap.String("shop_id", shopID))
		return nil, err
	}

	ss.logger.Info("success get shopee access token", zap.String("shop_id", shopID))
	return &result, nil
}

func (ss *shopeeService) RefreshAccessToken(ctx context.Context, refreshToken, shopID string) (*dto.ShopeeRefreshTokenResponse, error) {
	path := "/api/v2/auth/access_token/get"
	payload := map[string]interface{}{
		"refresh_token": refreshToken,
		"shop_id":       shopID,
		"partner_id":    ss.config.PartnerID,
	}

	var result dto.ShopeeRefreshTokenResponse
	err := ss.doRequest(
		ctx,
		http.MethodPost,
		path,
		nil,
		payload,
		&result,
	)
	if err != nil {
		ss.logger.Error("failed refresh shopee access token", zap.Error(err), zap.String("shop_id", shopID))
		return nil, err
	}

	ss.logger.Info("success refresh shopee access token", zap.String("shop_id", shopID))
	return &result, nil
}

func (ss *shopeeService) SyncProducts(ctx context.Context, storePlatform *entity.StorePlatform) error {
	if storePlatform == nil {
		return fmt.Errorf("store platform is nil: %w", dto.ErrBadRequest)
	}
	if storePlatform.Credential == nil {
		return fmt.Errorf("store credential not found: %w", dto.ErrBadRequest)
	}

	accessToken := storePlatform.Credential.AccessToken
	refreshToken := storePlatform.Credential.RefreshToken
	shopID := storePlatform.ExternalShopID

	ss.logger.Info("start sync shopee products", zap.String("store_platform_id", storePlatform.ID.String()), zap.String("shop_id", shopID))
	// 1. CHECK TOKEN EXPIRATION
	if storePlatform.Credential.ExpiresAt != nil && time.Now().After(*storePlatform.Credential.ExpiresAt) {
		ss.logger.Info("access token expired, refreshing token", zap.String("shop_id", shopID))

		refreshResp, err := ss.RefreshAccessToken(
			ctx,
			refreshToken,
			shopID,
		)
		if err != nil {
			ss.logger.Error("failed refresh shopee token", zap.String("shop_id", shopID), zap.Error(err))
			return fmt.Errorf("failed refresh access token: %w", dto.ErrUnauthorized)
		}

		accessToken = refreshResp.AccessToken
		refreshToken = refreshResp.RefreshToken
		expiredAt := time.Now().Add(
			time.Duration(refreshResp.ExpireIn) * time.Second,
		)

		storePlatform.Credential.AccessToken = refreshResp.AccessToken
		storePlatform.Credential.RefreshToken = refreshResp.RefreshToken
		storePlatform.Credential.ExpiresAt = &expiredAt

		err = ss.storeCredentialRepo.UpdateStoreCredential(ctx, nil, storePlatform.Credential)
		if err != nil {
			ss.logger.Error("failed update refreshed credential", zap.Error(err))
			return fmt.Errorf("failed update credential: %w", dto.ErrInternal)
		}

		ss.logger.Info("success refresh shopee token", zap.String("shop_id", shopID))
	}

	// 2. GET ITEM LIST
	itemListResp, err := ss.GetItemList(ctx, accessToken, shopID)
	if err != nil {
		ss.logger.Error("failed get shopee item list", zap.String("shop_id", shopID), zap.Error(err))
		return fmt.Errorf("failed get item list: %w", dto.ErrInternal)
	}

	ss.logger.Info("success get shopee item list", zap.Int("total_items", len(itemListResp.Response.Item)))

	// 3. LOOP ITEM IDS
	for _, item := range itemListResp.Response.Item {
		itemID := item.ItemID

		// 4. GET ITEM DETAIL
		itemDetailResp, err := ss.GetItemBaseInfo(ctx, accessToken, shopID, []int64{itemID})
		if err != nil {
			ss.logger.Error("failed get item detail", zap.Int64("item_id", itemID), zap.Error(err))
			continue
		}
		if len(itemDetailResp.Response.ItemList) == 0 {
			continue
		}

		itemDetail := itemDetailResp.Response.ItemList[0]

		ss.logger.Info("syncing shopee product", zap.Int64("item_id", itemID), zap.String("item_name", itemDetail.ItemName))

		// 5. UPSERT PRODUCT
		externalID := strconv.FormatInt(itemID, 10)
		externalProduct, found, err := ss.externalProductRepo.GetExternalProductByExternalID(ctx, nil, &storePlatform.ID, externalID)
		if err != nil {
			ss.logger.Error("failed get external product", zap.String("external_id", externalID), zap.Error(err))
			return fmt.Errorf("failed get external product: %w", dto.ErrInternal)
		}

		// ambil harga marketplace
		var price int64
		if len(itemDetail.PriceInfo) > 0 {
			price = int64(itemDetail.PriceInfo[0].CurrentPrice)
		}

		var stock int
		if len(itemDetail.StockInfo) > 0 {
			stock = itemDetail.StockInfo[0].CurrentStock
		}

		// CREATE NEW PRODUCT
		if !found {
			ss.tx.Run(ctx, func(tx *gorm.DB) error {
				product := &entity.Product{
					ID:          uuid.New(),
					StoreID:     *storePlatform.StoreID,
					Name:        itemDetail.ItemName,
					Description: itemDetail.Description,
					SKU:         "",
					Stock:       stock,
				}
				_, err = ss.productRepo.CreateProduct(ctx, nil, product)
				if err != nil {
					ss.logger.Error("failed create internal product", zap.String("external_id", externalID), zap.Error(err))
					return fmt.Errorf("failed create product: %w", dto.ErrInternal)
				}

				rawPayload, _ := json.Marshal(itemDetail)
				newExternalProduct := &entity.ExternalProduct{
					ID:              uuid.New(),
					StorePlatformID: &storePlatform.ID,
					ProductID:       &product.ID,
					ExternalID:      externalID,
					ExternalName:    itemDetail.ItemName,
					Price:           price,
					RawPayload:      datatypes.JSON(rawPayload),
					LastSyncAt:      helper.PtrTime(time.Now()),
				}

				_, err = ss.externalProductRepo.CreateExternalProduct(ctx, nil, newExternalProduct)
				if err != nil {
					ss.logger.Error("failed create external product", zap.String("external_id", externalID), zap.Error(err))
					return fmt.Errorf("failed create external product: %w", dto.ErrInternal)
				}

				ss.logger.Info("success create new external product", zap.String("external_id", externalID), zap.String("product_id", product.ID.String()))
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			// UPDATE EXISTING PRODUCT
			ss.tx.Run(ctx, func(tx *gorm.DB) error {
				product := externalProduct.Product
				if product == nil {
					ss.logger.Error("product relation is nil", zap.String("external_id", externalID))
					return fmt.Errorf("failed get product: %w", dto.ErrInternal)
				}

				product.Name = itemDetail.ItemName
				product.Description = itemDetail.Description

				_, err = ss.productRepo.UpdateProductByStoreIDAndProductID(ctx, nil, product)
				if err != nil {
					ss.logger.Error("failed update product", zap.String("external_id", externalID), zap.Error(err))
					return fmt.Errorf("failed update product: %w", dto.ErrInternal)
				}

				rawPayload, _ := json.Marshal(itemDetail)
				externalProduct.ExternalName = itemDetail.ItemName
				externalProduct.Price = price
				externalProduct.RawPayload = datatypes.JSON(rawPayload)

				now := time.Now()
				externalProduct.LastSyncAt = helper.PtrTime(now)

				_, err = ss.externalProductRepo.UpdateExternalProductByStoreIDAndExprodID(ctx, nil, externalProduct)
				if err != nil {
					ss.logger.Error("failed update external product", zap.String("external_id", externalID), zap.Error(err))
					return fmt.Errorf("failed update external product: %w", dto.ErrInternal)
				}

				ss.logger.Info("success update external product", zap.String("external_id", externalID), zap.String("product_id", product.ID.String()))
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	// 6. UPDATE LAST SYNC
	now := time.Now()
	storePlatform.LastSyncAt = &now

	err = ss.storePlatformRepo.UpdateStorePlatform(ctx, nil, storePlatform)
	if err != nil {
		ss.logger.Error("failed update last sync", zap.Error(err))
		return fmt.Errorf("failed update sync time: %w", dto.ErrInternal)
	}

	ss.logger.Info("success sync shopee products", zap.String("store_platform_id", storePlatform.ID.String()))

	return nil
}

func (ss *shopeeService) GetItemList(ctx context.Context, accessToken string, shopID string) (*dto.ShopeeGetItemListResponse, error) {
	path := "/api/v2/product/get_item_list"
	timestamp := time.Now().Unix()

	sign := helper.GenerateShopeeShopSign(
		ss.config.PartnerID,
		path,
		timestamp,
		accessToken,
		shopID,
		ss.config.PartnerKey,
	)

	endpoint := fmt.Sprintf(
		"%s%s?partner_id=%s&timestamp=%d&access_token=%s&shop_id=%s&sign=%s&offset=0&page_size=100&item_status=NORMAL",
		ss.config.BaseURL,
		path,
		ss.config.PartnerID,
		timestamp,
		accessToken,
		shopID,
		sign,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		ss.logger.Error("failed create request shopee get item list", zap.Error(err))
		return nil, err
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		ss.logger.Error("failed request shopee get item list", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ss.logger.Error("failed read response body shopee get item list", zap.Error(err))
		return nil, err
	}

	var result dto.ShopeeGetItemListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		ss.logger.Error("failed unmarshal get item list", zap.Error(err), zap.ByteString("response", body))
		return nil, err
	}

	if result.Error != "" {
		ss.logger.Error("shopee error", zap.String("message", result.Message))
		return nil, fmt.Errorf("shopee error: %s", result.Message)
	}

	return &result, nil
}

func (ss *shopeeService) GetItemBaseInfo(ctx context.Context, accessToken string, shopID string, itemIDs []int64) (*dto.ShopeeGetItemBaseInfoResponse, error) {
	path := "/api/v2/product/get_item_base_info"
	timestamp := time.Now().Unix()

	sign := helper.GenerateShopeeShopSign(
		ss.config.PartnerID,
		path,
		timestamp,
		accessToken,
		shopID,
		ss.config.PartnerKey,
	)

	itemIDStrings := make([]string, 0)
	for _, id := range itemIDs {
		itemIDStrings = append(itemIDStrings, strconv.FormatInt(id, 10))
	}

	itemIDParam := strings.Join(itemIDStrings, ",")
	endpoint := fmt.Sprintf(
		"%s%s?partner_id=%s&timestamp=%d&access_token=%s&shop_id=%s&sign=%s&item_id_list=%s",
		ss.config.BaseURL,
		path,
		ss.config.PartnerID,
		timestamp,
		accessToken,
		shopID,
		sign,
		itemIDParam,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		ss.logger.Error("failed create request shopee get item base info", zap.Error(err))
		return nil, err
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		ss.logger.Error("failed request shopee get item base info", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ss.logger.Error("failed read response body shopee get item base info", zap.Error(err))
		return nil, err
	}

	var result dto.ShopeeGetItemBaseInfoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		ss.logger.Error("failed unmarshal get item base info", zap.Error(err), zap.ByteString("response", body))
		return nil, err
	}

	if result.Error != "" {
		ss.logger.Error("shopee error", zap.String("message", result.Message))
		return nil, fmt.Errorf("shopee error: %s", result.Message)
	}

	return &result, nil
}
