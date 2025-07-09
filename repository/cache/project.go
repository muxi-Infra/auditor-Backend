package cache

import (
	"context"
	"strconv"
	"time"
)

const TagsKey = "MuxiAuditor:tags:"
const TagExpiration = 10 * time.Minute

// ProjectCacheInterface 对真实的cache进行了一层封装方便拓展和更换
type ProjectCacheInterface interface {
	GetStringSlice(ctx context.Context, key string) ([]string, error)
	SetStringSlice(ctx context.Context, key string, val []string, expiration time.Duration) error
}
type ProjectCache struct {
	Ca ProjectCacheInterface
}

func NewProjectCache(ca ProjectCacheInterface) *ProjectCache {
	return &ProjectCache{Ca: ca}
}
func (p *ProjectCache) GetAllTags(ctx context.Context, pid uint) ([]string, error) {
	key := TagsKey + strconv.FormatUint(uint64(pid), 10)
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
