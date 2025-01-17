package data

type NotifyHeap []*Notify

func (n NotifyHeap) Peek() any {
	return n[0]
}

func (n NotifyHeap) Less(i, j int) bool {
	return n[i].Line < n[j].Line
}
func (n *NotifyHeap) Pop() any {
	res := (*n)[len(*n)-1]
	*n = (*n)[:len(*n)-1]
	return res
}
func (n *NotifyHeap) Push(v any) {
	*n = append(*n, v.(*Notify))
}
func (n NotifyHeap) Len() int {
	return len(n)
}
func (n NotifyHeap) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
