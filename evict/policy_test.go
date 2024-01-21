package evict

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLruPolicy_Push(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p := NewLruPolicy[string]()

		element1 := p.Push("key-1")
		p.Push("key-2")

		require.NotNil(t, element1.Prev())
		require.Equal(t, 2, p.Len())
	})
}

func TestLruPolicy_Evict(t *testing.T) {
	t.Run("success - push items", func(t *testing.T) {
		p := NewLruPolicy[string]()

		element1 := p.Push("key-1")
		element2 := p.Push("key-2")
		element3 := p.Push("key-3")

		p.Evict()

		require.Equal(t, 2, p.list.Len())
		require.Nil(t, element1.Next())
		require.Nil(t, element1.Prev())
		require.Equal(t, element2, element3.Next())
		require.Equal(t, element3, element2.Prev())
	})

	t.Run("success - increase score for last item", func(t *testing.T) {
		p := NewLruPolicy[string]()

		element1 := p.Push("key-1")
		element2 := p.Push("key-2")
		element3 := p.Push("key-3")

		p.IncScore(element1)

		p.Evict()

		require.Equal(t, 2, p.list.Len())
		require.Nil(t, element2.Next())
		require.Nil(t, element2.Prev())
		require.Equal(t, element3, element1.Next())
		require.Equal(t, element1, element3.Prev())
	})
}

func TestLruPolicy_Del(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p := NewLruPolicy[string]()

		element1 := p.Push("key-1")

		key := p.Del(element1)

		require.Equal(t, "key-1", key)
		require.Nil(t, element1.Next())
		require.Nil(t, element1.Prev())
	})
}
