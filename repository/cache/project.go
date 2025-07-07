package cache

type ProjectCacheInterface interface {
	SetAllTags()
}
type ProjectCache struct {
	Ca Cache
}

func NewProjectCache(ca Cache) *ProjectCache {
	return &ProjectCache{Ca: ca}
}
func (p *ProjectCache) SetAllTags() {}
