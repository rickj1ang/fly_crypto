package main

import (
	"container/heap"

	"github.com/rickj1ang/fly_crypto/internal/data"
)

func main() {
	notifys := &data.NotifyHeap{}
	heap.Init(notifys)

}
