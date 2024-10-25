package server

import (
	"context"
	"encoding/json"
	"github.com/go-oauth2/oauth2/v4"
	"store/pkg/config"
	"store/pkg/redis"
	"store/pkg/tools"
	"time"
)

type TokeStore struct {
	rdb  *redis.Client
	conf *config.GlobalConfig
}

func NewTokeStore(rdb *redis.Client, conf *config.GlobalConfig) *TokeStore {
	return &TokeStore{
		rdb:  rdb,
		conf: conf,
	}
}

func (t *TokeStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	code := tools.CreateID()
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	if err := t.rdb.Set(ctx, code, data, time.Hour*time.Duration(24*t.conf.JWT.RefreshExpiry)); err != nil {
		return err
	}

	if err := t.rdb.Set(ctx, info.GetAccess(), code, time.Hour*time.Duration(24*t.conf.JWT.AccessExpiry)); err != nil {
		return err
	}

	if err := t.rdb.Set(ctx, info.GetRefresh(), code, time.Hour*time.Duration(24*t.conf.JWT.RefreshExpiry)); err != nil {
		return err
	}

	return nil
}

// delete the authorization code
func (t *TokeStore) RemoveByCode(ctx context.Context, code string) error {
	data, err := t.rdb.Get(ctx, code)
	if err != nil {
		return err
	}

	var info oauth2.TokenInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return err
	}

	return t.rdb.Del(ctx, code, info.GetAccess(), info.GetRefresh())
}

// use the access token to delete the token information
func (t *TokeStore) RemoveByAccess(ctx context.Context, access string) error {
	code, err := t.rdb.Get(ctx, access)
	if err != nil {
		return err
	}

	data, err := t.rdb.Get(ctx, code)
	if err != nil {
		return err
	}

	var info oauth2.TokenInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return err
	}

	return t.rdb.Del(ctx, access, code, info.GetRefresh())
}

// use the refresh token to delete the token information
func (t *TokeStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	code, err := t.rdb.Get(ctx, refresh)
	if err != nil {
		return err
	}

	data, err := t.rdb.Get(ctx, code)
	if err != nil {
		return err
	}

	var info oauth2.TokenInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return err
	}

	return t.rdb.Del(ctx, refresh, code, info.GetAccess())
}

// use the authorization code for token information data
func (t *TokeStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	data, err := t.rdb.Get(ctx, code)
	if err != nil {
		return nil, err
	}

	var info oauth2.TokenInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, err
	}
	return info, nil
}

// use the access token for token information data
func (t *TokeStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	code, err := t.rdb.Get(ctx, access)
	if err != nil {
		return nil, err
	}

	data, err := t.rdb.Get(ctx, code)
	if err != nil {
		return nil, err
	}

	var info oauth2.TokenInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, err
	}
	return info, nil
}

// use the refresh token for token information data
func (t *TokeStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	code, err := t.rdb.Get(ctx, refresh)
	if err != nil {
		return nil, err
	}

	data, err := t.rdb.Get(ctx, code)
	if err != nil {
		return nil, err
	}

	var info oauth2.TokenInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, err
	}

	return info, nil
}
