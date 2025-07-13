package pokecache

import (
	"testing"
	"time"
	//"github.com/ChrisR/internal/pokecache"
)

func TestCache(t *testing.T) {
	type data struct {
		interval time.Duration
		key      []string
		val      [][]byte
	}

	cases := []struct {
		input    data
		expected data
	}{
		{
			input:    data{interval: 5, key: []string{"test", "the", "cache"}, val: [][]byte{[]byte("testing"), []byte("this"), []byte("item")}},
			expected: data{interval: 1, key: []string{"test", "the", "cache"}, val: [][]byte{[]byte("testing"), []byte("this"), []byte("item")}},
		},
	}

	for _, val := range cases {
		cache := NewCache(val.input.interval)

		for i := range val.input.key {
			cache.Add(val.input.key[i], val.input.val[i])
		}

		shutdown := make(chan bool)
		ticker := time.NewTicker(val.input.interval)
		go func() {

			// Using for loop
			for {

				// Select statement
				select {

				// Case statement
				case <-shutdown:
					return

				// Case to print current time
				case <-ticker.C:
					for i, v := range val.expected.key {
						data, found := cache.Get(v)
						if !found {
							t.Errorf(" %s; expected %s", data, val.expected.val[i])
						}
					}
				}
			}
		}()

		time.Sleep(val.input.interval)
		shutdown <- true
		ticker.Stop()

		for _, v := range val.expected.key {
			data, found := cache.Get(v)
			if found {
				t.Errorf(" %s; expected %s", data, "nil")
			}
		}

		cache.ShutdownCacheCollection()
	}

}
