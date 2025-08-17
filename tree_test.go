package critbit

import (
	"math/rand/v2"
	"testing"
)

type TestData struct {
	Key   Key
	Value uint32
}

func setupDataset(n int) []TestData {
	dataset := make([]TestData, n)
	for i := range n {
		dataset[i] = TestData{
			Key:   Uint32Key(uint32(i)),
			Value: uint32(i),
		}
	}
	return dataset
}

func TestTreeUint32(t *testing.T) {
	N := 256
	r := rand.New(rand.NewPCG(1, 1))
	var m Tree[uint32]
	dataset := setupDataset(N)
	t.Run("Set", func(t *testing.T) {
		for _, p := range r.Perm(N) {
			data := dataset[p]
			m.Set(data.Key, data.Value)
		}
	})
	t.Run("Replace", func(t *testing.T) {
		for _, p := range r.Perm(N) {
			data := dataset[p]
			m.Set(data.Key, data.Value)
		}
	})
	t.Run("Len", func(t *testing.T) {
		if m.Len() != len(dataset) {
			t.Errorf("want %v; but got %v", len(dataset), m.Len())
		}
	})
	t.Run("Get", func(t *testing.T) {
		for _, data := range dataset {
			val, found := m.Get(data.Key)
			if !found {
				t.Fatalf("%x not found", data.Key)
			}
			if val != data.Value {
				t.Errorf("want %v; but got %v", data.Value, val)
			}
		}
	})
	t.Run("Get unkown key", func(t *testing.T) {
		key := Key{}
		_, found := m.Get(key)
		if found {
			t.Fatalf("%x exists", key)
		}
	})
	t.Run("All", func(t *testing.T) {
		i := 0
		for key, val := range m.All() {
			data := dataset[i]
			if !key.Equal(data.Key) {
				t.Errorf("want %v; but got %v", data.Key, key)
			}
			if val != data.Value {
				t.Errorf("want %v; but got %v", data.Value, val)
			}
			i++
		}
	})
	t.Run("All break", func(t *testing.T) {
		i := 0
		for range m.All() {
			if i >= len(dataset)/2 {
				break
			}
			i++
		}
	})
	t.Run("Keys", func(t *testing.T) {
		i := 0
		for key := range m.Keys() {
			data := dataset[i]
			if !key.Equal(data.Key) {
				t.Errorf("want %v; but got %v", data.Key, key)
			}
			i++
		}
	})
	t.Run("Keys break", func(t *testing.T) {
		i := 0
		for range m.Keys() {
			if i >= len(dataset)/2 {
				break
			}
			i++
		}
	})
	t.Run("Values", func(t *testing.T) {
		i := 0
		for val := range m.Values() {
			data := dataset[i]
			if val != data.Value {
				t.Errorf("want %v; but got %v", data.Value, val)
			}
			i++
		}
	})
	t.Run("Values break", func(t *testing.T) {
		i := 0
		for range m.Values() {
			if i >= len(dataset)/2 {
				break
			}
			i++
		}
	})
	t.Run("Reverse", func(t *testing.T) {
		i := len(dataset) - 1
		s := NewScanner(m.root, true)
		for {
			leaf := s.Scan()
			if leaf == nil {
				break
			}
			data := dataset[i]
			if !leaf.Key.Equal(data.Key) {
				t.Errorf("want %v; but got %v", data.Key, leaf.Key)
			}
			if leaf.Value != data.Value {
				t.Errorf("want %v; but got %v", data.Value, leaf.Value)
			}
			i--
		}
	})
	t.Run("Delete", func(t *testing.T) {
		for _, p := range r.Perm(N) {
			data := dataset[p]
			m.Delete(data.Key)
			// already deleted
			m.Delete(data.Key)
		}
	})
	t.Run("empty Delete", func(t *testing.T) {
		for _, p := range r.Perm(N) {
			data := dataset[p]
			m.Delete(data.Key)
		}
	})
	t.Run("Len", func(t *testing.T) {
		if m.Len() != 0 {
			t.Errorf("want %v; but got %v", 0, m.Len())
		}
	})
	t.Run("Get", func(t *testing.T) {
		for _, data := range dataset {
			_, found := m.Get(data.Key)
			if found {
				t.Fatalf("%x exists", data.Key)
			}
		}
	})
}
