package D

import (
	"fmt"
	"time"

	"github.com/allegro/bigcache"
	"github.com/kpango/fastime"

	"github.com/kokizzu/gotro/L"

	cmap "github.com/orcaman/concurrent-map"
)

type RamCache struct {
	evictionLogic *bigcache.BigCache
	store         cmap.ConcurrentMap
	expireSec     int64
}

func NewRamCache(dur time.Duration, sizeMB int) *RamCache {
	cfg := bigcache.DefaultConfig(dur)
	cfg.HardMaxCacheSize = sizeMB
	res := &RamCache{
		store:     cmap.New(),
		expireSec: int64(dur.Seconds()),
	}
	cfg.OnRemove = func(key string, entry []byte) {
		res.store.Remove(key)
	}
	expireLogic, err := bigcache.NewBigCache(cfg)
	L.PanicIf(err, `bigcache.NewBigCache failed`)
	res.evictionLogic = expireLogic
	return res
}

func (r *RamCache) Set(key string, value any) {
	suffix := r.secondSuffix()
	if r.evictionLogic != nil && r.evictionLogic.Set(key+suffix, []byte{1}) == nil {
		r.store.Set(key+suffix, value)
	}
}

func (r *RamCache) ClearAll() {
	r.store.Clear()
	L.IsError(r.evictionLogic.Reset(), `RamCache.ClearAll`)
}

func (r *RamCache) Get(key string) any {
	res, ok := r.store.Get(key + r.secondSuffix())
	if !ok {
		return nil
	}
	return res
}

func (r *RamCache) Delete(k string) {
	suffix := r.secondSuffix()
	r.store.Remove(k + suffix)
	L.IsError(r.evictionLogic.Delete(k+suffix), `RamCache.Delete`)
}

// force evict every n seconds
func (r *RamCache) secondSuffix() string {
	if r.expireSec < 2 {
		return ``
	}
	return `|` + fmt.Sprint(fastime.UnixNow()/r.expireSec)
}
