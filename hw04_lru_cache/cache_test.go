package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	c := NewCache(5)

	wasInCache := c.Set("aaa", 100)
	require.False(t, wasInCache)

	wasInCache = c.Set("bbb", 200)
	require.False(t, wasInCache)

	val, ok := c.Get("aaa")
	require.True(t, ok)
	require.Equal(t, 100, val)

	val, ok = c.Get("bbb")
	require.True(t, ok)
	require.Equal(t, 200, val)

	wasInCache = c.Set("aaa", 300)
	require.True(t, wasInCache)

	val, ok = c.Get("aaa")
	require.True(t, ok)
	require.Equal(t, 300, val)

	val, ok = c.Get("ccc")
	require.False(t, ok)
	require.Nil(t, val)
}

func TestCacheClear(t *testing.T) {
	c := NewCache(10)

	c.Set("aaa", 100)
	c.Set("bbb", 200)
	c.Clear()

	_, ok := c.Get("aaa")
	require.False(t, ok)
}

func TestCachePure(t *testing.T) {
	c := NewCache(3)

	c.Set("aaa", 100)
	c.Set("bbb", 200)
	c.Set("ccc", 400)
	c.Set("ddd", 500)

	_, ok := c.Get("aaa")
	require.False(t, ok)

	c.Set("bbb", 250)
	c.Set("eee", 600)

	_, ok = c.Get("ccc")
	require.False(t, ok)
}

func TestCacheEmpty(t *testing.T) {
	c := NewCache(10)

	_, ok := c.Get("aaa")
	require.False(t, ok)

	_, ok = c.Get("bbb")
	require.False(t, ok)
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
