package cache

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCache_GetSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")
		c.Set("key-3", "val-3")

		act, ok := c.Get("key-2")
		require.True(t, ok)
		require.Equal(t, "val-2", act)
	})

	t.Run("miss value", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")
		c.Set("key-3", "val-3")

		act, ok := c.Get("unknown")
		require.False(t, ok)
		require.Empty(t, act)
	})

	t.Run("success - get some values with miss key (check locks)", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")
		c.Set("key-3", "val-3")

		act, ok := c.Get("key-1")
		require.True(t, ok)
		require.Equal(t, "val-1", act)

		act, ok = c.Get("unknown")
		require.False(t, ok)
		require.Equal(t, "", act)

		act, ok = c.Get("key-2")
		require.True(t, ok)
		require.Equal(t, "val-2", act)
	})
}

func TestCache_SetWithTTL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.SetWithTTL("key", "val", 10*time.Second)

		act, ok := c.Get("key")
		require.True(t, ok)
		require.Equal(t, "val", act)
	})

	t.Run("set value with empty ttl", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.SetWithTTL("key-1", "val-1", 0)

		c.now = func() int64 {
			return time.Now().Add(11 * time.Second).Unix()
		}

		act, ok := c.Get("key-1")
		require.True(t, ok)
		require.Equal(t, "val-1", act)
	})

	t.Run("miss value", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.SetWithTTL("key-1", "val-1", 10*time.Second)

		act, ok := c.Get("unknown")
		require.False(t, ok)
		require.Empty(t, act)
	})

	t.Run("expired value", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.SetWithTTL("key-1", "val-1", 10*time.Second)

		c.now = func() int64 {
			return time.Now().Add(11 * time.Second).Unix()
		}

		act, ok := c.Get("key-1")
		require.False(t, ok)
		require.Empty(t, act)
	})
}

func TestCache_Has(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")
		c.Set("key-3", "val-3")

		ok := c.Has("key-2")
		require.True(t, ok)
	})

	t.Run("miss value", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")
		c.Set("key-3", "val-3")

		ok := c.Has("unknown")
		require.False(t, ok)
	})
}

func TestCache_Del(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")

		c.Del("key-2")

		require.False(t, c.Has("key-2"))
		require.True(t, c.Has("key-1"))
	})

	t.Run("miss value", func(t *testing.T) {
		c := NewCache[string, string](10)

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")

		c.Del("unknown")

		require.True(t, c.Has("key-2"))
		require.True(t, c.Has("key-1"))
	})
}

func TestCache_LRU_Evict(t *testing.T) {
	t.Run("success - evict last recently used item by set function", func(t *testing.T) {
		c := NewCache[string, string](2)

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")

		c.Set("key-3", "val-3")

		require.False(t, c.Has("key-1"))
		require.True(t, c.Has("key-2"))
		require.True(t, c.Has("key-3"))
		require.Equal(t, 2, c.evictPolicy.Len())
		require.Equal(t, 2, len(c.items))
	})

	t.Run("success - loop evict", func(t *testing.T) {
		c := NewCache[string, string](2)

		for i := 0; i < 100; i++ {
			c.Set(fmt.Sprintf("key-%v", i), fmt.Sprintf("val-%v", i))
		}

		require.False(t, c.Has("key-97"))
		require.True(t, c.Has("key-98"))
		require.True(t, c.Has("key-99"))
		require.Equal(t, 2, c.evictPolicy.Len())
		require.Equal(t, 2, len(c.items))
	})

	t.Run("success - evict last recently used item by has function", func(t *testing.T) {
		c := NewCache[string, string](3)

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")
		c.Set("key-3", "val-3")

		c.Has("key-1")

		c.Set("key-4", "val-4")

		require.False(t, c.Has("key-2"))
		require.Equal(t, 3, c.evictPolicy.Len())
		require.Equal(t, 3, len(c.items))
	})
}
