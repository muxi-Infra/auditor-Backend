package cache

import (
	"context"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/stringx"
	"strconv"
	"time"
)

const (
	TagsKey             = "MuxiAuditor:tags:"
	TagExpiration       = 10 * time.Minute
	AuditRoleKey        = "MuxiAuditor:auditRole:"
	AuditRoleExpiration = 1 * time.Hour
)

type CacheInterface interface {
	GetStringSlice(ctx context.Context, key string) ([]string, error)
	SetStringSlice(ctx context.Context, key string, val []string, expiration time.Duration) error
	SetString(ctx context.Context, key string, val string, expiration time.Duration) error
	GetString(ctx context.Context, key string) (string, error)
}

type ProjectCacheInterface interface {
	GetAllTags(ctx context.Context, pid uint) ([]string, error)
	SetAllTags(ctx context.Context, pid uint, tags []string) error
	GetAuditRole(ctx context.Context, pid uint) (string, error)
	SetAuditRole(ctx context.Context, pid uint, role string) error
}

type ProjectCache struct {
	Ca CacheInterface
}

func NewProjectCache(ca CacheInterface) *ProjectCache {
	return &ProjectCache{Ca: ca}
}

func (p *ProjectCache) GetAllTags(ctx context.Context, pid uint) ([]string, error) {
	key := stringx.Build(TagsKey, strconv.FormatUint(uint64(pid), 10))
	v, err := p.Ca.GetStringSlice(ctx, key)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (p *ProjectCache) SetAllTags(ctx context.Context, pid uint, tags []string) error {
	key := TagsKey + strconv.FormatUint(uint64(pid), 10)
	return p.Ca.SetStringSlice(ctx, key, tags, TagExpiration)
}

func (p *ProjectCache) SetAuditRole(ctx context.Context, pid uint, role string) error {
	key := stringx.Build(AuditRoleKey, strconv.FormatUint(uint64(pid), 10))
	return p.Ca.SetString(ctx, key, role, AuditRoleExpiration)
}

func (p *ProjectCache) GetAuditRole(ctx context.Context, pid uint) (string, error) {
	key := stringx.Build(AuditRoleKey, strconv.FormatUint(uint64(pid), 10))
	v, err := p.Ca.GetString(ctx, key)
	if err != nil {
		return "", err
	}
	return v, nil
}
