package cache

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCache_GetSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewCache[string, string]()

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")
		c.Set("key-3", "val-3")

		act, ok := c.Get("key-2")
		require.True(t, ok)
		require.Equal(t, "val-2", act)
	})

	t.Run("miss value", func(t *testing.T) {
		c := NewCache[string, string]()

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")
		c.Set("key-3", "val-3")

		act, ok := c.Get("unknown")
		require.False(t, ok)
		require.Empty(t, act)
	})
}

func TestCache_Has(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewCache[string, string]()

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")
		c.Set("key-3", "val-3")

		ok := c.Has("key-2")
		require.True(t, ok)
	})

	t.Run("miss value", func(t *testing.T) {
		c := NewCache[string, string]()

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")
		c.Set("key-3", "val-3")

		ok := c.Has("unknown")
		require.False(t, ok)
	})
}

func TestCache_Del(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := NewCache[string, string]()

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")

		c.Del("key-2")

		require.False(t, c.Has("key-2"))
		require.True(t, c.Has("key-1"))
	})

	t.Run("miss value", func(t *testing.T) {
		c := NewCache[string, string]()

		c.Set("key-1", "val-1")
		c.Set("key-2", "val-2")

		c.Del("unknown")

		require.True(t, c.Has("key-2"))
		require.True(t, c.Has("key-1"))
	})
}
