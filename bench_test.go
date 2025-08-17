package critbit

import (
	"math/rand/v2"
	"testing"
	"unsafe"
)

func Comparable(k Key) string {
	return unsafe.String(unsafe.SliceData(k.Data), len(k.Data))
}

func BenchmarkSet(b *testing.B) {
	N := 1024 * 256
	dataset := setupDataset(N)
	r := rand.New(rand.NewPCG(1, 1))
	perm := r.Perm(N)
	var m Tree[uint32]
	for _, p := range perm[:N/2] {
		data := dataset[p]
		m.Set(data.Key, data.Value)
	}
	perm = r.Perm(N)
	n := 0
	for b.Loop() {
		data := dataset[perm[n]]
		m.Set(data.Key, data.Value)
		n = (n + 1) % N
	}
}

func BenchmarkGet(b *testing.B) {
	N := 1024 * 256
	dataset := setupDataset(N)
	r := rand.New(rand.NewPCG(1, 1))
	perm := r.Perm(N)
	var m Tree[uint32]
	for _, p := range perm[:N/2] {
		data := dataset[p]
		m.Set(data.Key, data.Value)
	}
	perm = r.Perm(N)
	n := 0
	for b.Loop() {
		data := dataset[perm[n]]
		_, _ = m.Get(data.Key)
		n = (n + 1) % N
	}
}

func BenchmarkDelete(b *testing.B) {
	N := 1024 * 256
	dataset := setupDataset(N)
	r := rand.New(rand.NewPCG(1, 1))
	perm := r.Perm(N)
	var m Tree[uint32]
	for _, p := range perm[:N/2] {
		data := dataset[p]
		m.Set(data.Key, data.Value)
	}
	perm = r.Perm(N)
	n := 0
	for b.Loop() {
		data := dataset[perm[n]]
		m.Delete(data.Key)
		n = (n + 1) % N
	}
}

func BenchmarkStdMapSet(b *testing.B) {
	N := 1024 * 256
	dataset := setupDataset(N)
	r := rand.New(rand.NewPCG(1, 1))
	m := make(map[string]any)
	perm := r.Perm(N)
	for _, p := range perm[:N/2] {
		data := dataset[p]
		m[Comparable(data.Key)] = data.Value
	}
	perm = r.Perm(N)
	n := 0
	for b.Loop() {
		data := dataset[perm[n]]
		m[Comparable(data.Key)] = data.Value
		n = (n + 1) % N
	}
}

func BenchmarkStdMapGet(b *testing.B) {
	N := 1024 * 256
	dataset := setupDataset(N)
	r := rand.New(rand.NewPCG(1, 1))
	m := make(map[string]uint32)
	perm := r.Perm(N)
	for _, p := range perm[:N/2] {
		data := dataset[p]
		m[Comparable(data.Key)] = data.Value
	}
	perm = r.Perm(N)
	n := 0
	for b.Loop() {
		data := dataset[perm[n]]
		_ = m[Comparable(data.Key)]
		n = (n + 1) % N
	}
}

func BenchmarkStdMapDelete(b *testing.B) {
	N := 1024 * 256
	dataset := setupDataset(N)
	r := rand.New(rand.NewPCG(1, 1))
	m := make(map[string]uint32)
	perm := r.Perm(N)
	for _, p := range perm[:N/2] {
		data := dataset[p]
		m[Comparable(data.Key)] = data.Value
	}
	perm = r.Perm(N)
	n := 0
	for b.Loop() {
		data := dataset[perm[n]]
		delete(m, Comparable(data.Key))
		n = (n + 1) % N
	}
}

func BenchmarkOverhead(b *testing.B) {
	N := 1024 * 256
	dataset := setupDataset(N)
	r := rand.New(rand.NewPCG(1, 1))
	perm := r.Perm(N)
	n := 0
	for b.Loop() {
		_ = dataset[perm[n]]
		n = (n + 1) % N
	}
}
