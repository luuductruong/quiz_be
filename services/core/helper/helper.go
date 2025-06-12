package helper

import (
	"github.com/google/uuid"
	"math/rand"
	"sort"
	"time"
)

func MapList[S any, D any](source []S, mapper func(S) D) []D {
	rs := make([]D, len(source))
	for i := range source {
		rs[i] = mapper(source[i])
	}
	return rs
}

func RandString(n int) string {
	var source = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	max := len(source)
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = source[rand.Intn(max)]
	}

	return string(b)
}

func NewStringUUID() string {
	return uuid.New().String()
}

// Unique ...
func Unique[T comparable](arr []T) []T {
	u := make([]T, 0)
	m := make(map[T]bool)
	for _, val := range arr {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

// SelectMap copy T into T1, which each must have selectFunc true
func SelectMap[T any](src []T, selectFunc func(T) bool) []T {
	res := make([]T, 0)
	for _, v := range src {
		if selectFunc(v) {
			res = append(res, v)
		}
	}
	return res
}

// Sort clone T to T_after_sorting
func Sort[T any](src []T, less func(i, j int) bool) []T {
	res := make([]T, 0)
	sort.Slice(src, func(i, j int) bool {
		return less(i, j)
	})
	for _, v := range src {
		res = append(res, v)
	}
	return res
}

// AnyToPointer ...
func AnyToPointer[T any](src T) *T {
	return &src
}

// UniqBy ...
func UniqBy[T any, U comparable](collection []T, iteratee func(T) U) []T {
	result := make([]T, 0, len(collection))
	seen := make(map[U]struct{}, len(collection))

	for _, item := range collection {
		key := iteratee(item)

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, item)
	}

	return result
}
