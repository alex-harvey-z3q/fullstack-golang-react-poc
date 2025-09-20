package main

import (
	"bufio"
	"container/list"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

///////////////////////////////
// Strings & Arrays
///////////////////////////////

// ReverseStringUnicode reverses a string correctly for Unicode.
func ReverseStringUnicode(s string) string {
	r := []rune(s) // In Go, strings are slices of UTF-8 bytes.
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// IsPalindromeUnicode checks palindrome with Unicode runes.
func IsPalindromeUnicode(s string) bool {
	// Optional: normalize by removing spaces/punctuation; here we do a simple check.
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		if r[i] != r[j] {
			return false
		}
	}
	return true
}

// AreAnagrams checks if two strings are anagrams (Unicode-aware).
func AreAnagrams(a, b string) bool {
	ra, rb := []rune(a), []rune(b)
	if len(ra) != len(rb) {
		return false
	}
	m := make(map[rune]int, len(ra))
	for _, ch := range ra {
		m[ch]++
	}
	for _, ch := range rb {
		m[ch]--
		if m[ch] < 0 {
			return false
		}
	}
	return true
}

// KadaneMaxSubarray returns the maximum subarray sum.
func KadaneMaxSubarray(nums []int) int {
	best, cur := math.MinInt, 0
	for _, x := range nums {
		if cur > 0 {
			cur += x
		} else {
			cur = x
		}
		if cur > best {
			best = cur
		}
	}
	return best
}

// TwoSum returns indices of the two numbers adding to target (first found).
func TwoSum(nums []int, target int) (i, j int, ok bool) {
	m := make(map[int]int, len(nums))
	for idx, x := range nums {
		if k, found := m[target-x]; found {
			return k, idx, true
		}
		m[x] = idx
	}
	return 0, 0, false
}

///////////////////////////////
// Data Structures & Algorithms
///////////////////////////////

// Stack using slice
type Stack[T any] struct{ data []T }

func (s *Stack[T]) Push(v T) { s.data = append(s.data, v) }
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return v, true
}
func (s *Stack[T]) Len() int { return len(s.data) }

// Queue using circular slice
type Queue[T any] struct {
	data       []T
	head, size int
}

func NewQueue[T any](cap int) *Queue[T] { return &Queue[T]{data: make([]T, cap)} }
func (q *Queue[T]) grow() {
	newData := make([]T, max(2*len(q.data), 1))
	for i := 0; i < q.size; i++ {
		newData[i] = q.data[(q.head+i)%len(q.data)]
	}
	q.data, q.head = newData, 0
}
func (q *Queue[T]) Enqueue(v T) {
	if q.size == len(q.data) {
		q.grow()
	}
	q.data[(q.head+q.size)%len(q.data)] = v
	q.size++
}
func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}
	v := q.data[q.head]
	q.head = (q.head + 1) % len(q.data)
	q.size--
	return v, true
}
func (q *Queue[T]) Len() int { return q.size }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Singly linked list node
type ListNode struct {
	Val  int
	Next *ListNode
}

// ReverseLinkedList reverses in-place.
func ReverseLinkedList(head *ListNode) *ListNode {
	var prev *ListNode
	for head != nil {
		nxt := head.Next
		head.Next = prev
		prev = head
		head = nxt
	}
	return prev
}

// DetectCycle returns the node where a cycle begins (Floyd), or nil.
func DetectCycle(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			p := head
			for p != slow {
				p = p.Next
				slow = slow.Next
			}
			return p
		}
	}
	return nil
}

// MiddleNode returns the middle node (second middle for even length).
func MiddleNode(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}

// Binary search (iterative)
func BinarySearch(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		switch {
		case nums[mid] == target:
			return mid
		case nums[mid] < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return -1
}

// Binary search (recursive)
func BinarySearchRec(nums []int, target int) int {
	var f func(int, int) int
	f = func(lo, hi int) int {
		if lo > hi {
			return -1
		}
		mid := lo + (hi-lo)/2
		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			return f(mid+1, hi)
		}
		return f(lo, mid-1)
	}
	return f(0, len(nums)-1)
}

// Binary tree node
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// Traversals (recursive)
func Preorder(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	out := []int{root.Val}
	out = append(out, Preorder(root.Left)...)
	out = append(out, Preorder(root.Right)...)
	return out
}
func Inorder(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	out := Inorder(root.Left)
	out = append(out, root.Val)
	out = append(out, Inorder(root.Right)...)
	return out
}
func Postorder(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	out := Postorder(root.Left)
	out = append(out, Postorder(root.Right)...)
	out = append(out, root.Val)
	return out
}
func LevelOrder(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	var res []int
	q := []*TreeNode{root}
	for len(q) > 0 {
		n := q[0]
		q = q[1:]
		res = append(res, n.Val)
		if n.Left != nil {
			q = append(q, n.Left)
		}
		if n.Right != nil {
			q = append(q, n.Right)
		}
	}
	return res
}

// Graph DFS/BFS (adjacency list)
func DFS(start int, adj map[int][]int) []int {
	var res []int
	seen := map[int]bool{}
	var visit func(int)
	visit = func(u int) {
		seen[u] = true
		res = append(res, u)
		for _, v := range adj[u] {
			if !seen[v] {
				visit(v)
			}
		}
	}
	visit(start)
	return res
}
func BFS(start int, adj map[int][]int) []int {
	var res []int
	seen := map[int]bool{start: true}
	q := []int{start}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		res = append(res, u)
		for _, v := range adj[u] {
			if !seen[v] {
				seen[v] = true
				q = append(q, v)
			}
		}
	}
	return res
}

///////////////////////////////
// Concurrency & Goroutines
///////////////////////////////

// WorkerPool processes jobs concurrently and returns results. Safe to close jobs when done sending.
func WorkerPool[K any, V any](ctx context.Context, workers int, jobs <-chan K, fn func(context.Context, K) (V, error)) <-chan struct {
	Val V
	Err error
} {
	out := make(chan struct {
		Val V
		Err error
	})
	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case j, ok := <-jobs:
				if !ok {
					return
				}
				v, err := fn(ctx, j)
				select {
				case <-ctx.Done():
					return
				case out <- struct {
					Val V
					Err error
				}{v, err}:
				}
			}
		}
	}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// ProducerConsumer demonstrates using channels.
func ProducerConsumer(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			out <- i
		}
	}()
	return out
}

// SimpleRateLimiter token-bucket with capacity 1, interval d.
type SimpleRateLimiter struct {
	ticker *time.Ticker
	tokens chan struct{}
	stop   chan struct{}
}

func NewSimpleRateLimiter(d time.Duration) *SimpleRateLimiter {
	rl := &SimpleRateLimiter{
		ticker: time.NewTicker(d),
		tokens: make(chan struct{}, 1),
		stop:   make(chan struct{}),
	}
	// seed
	rl.tokens <- struct{}{}
	go func() {
		for {
			select {
			case <-rl.stop:
				return
			case <-rl.ticker.C:
				select { // non-blocking refill
				case rl.tokens <- struct{}{}:
				default:
				}
			}
		}
	}()
	return rl
}
func (r *SimpleRateLimiter) Allow() {
	<-r.tokens
}
func (r *SimpleRateLimiter) Close() {
	close(r.stop)
	r.ticker.Stop()
}

// DeadlockAvoidance via consistent lock ordering
type SafePair struct {
	mu1, mu2 sync.Mutex
	a, b     int
}

// IncrementBoth locks in a fixed order to avoid deadlocks.
func (p *SafePair) IncrementBoth() {
	// Always lock mu1 then mu2 everywhere in codebase.
	p.mu1.Lock()
	p.mu2.Lock()
	p.a++
	p.b++
	p.mu2.Unlock()
	p.mu1.Unlock()
}

///////////////////////////////
// Sorting & Searching
///////////////////////////////

// BubbleSort (educational)
func BubbleSort(nums []int) {
	n := len(nums)
	for swapped := true; swapped; {
		swapped = false
		for i := 1; i < n; i++ {
			if nums[i-1] > nums[i] {
				nums[i-1], nums[i] = nums[i], nums[i-1]
				swapped = true
			}
		}
		n--
	}
}

// Quicksort (in-place)
func QuickSort(nums []int) {
	var qs func(int, int)
	partition := func(lo, hi int) int {
		pivot := nums[hi]
		i := lo
		for j := lo; j < hi; j++ {
			if nums[j] < pivot {
				nums[i], nums[j] = nums[j], nums[i]
				i++
			}
		}
		nums[i], nums[hi] = nums[hi], nums[i]
		return i
	}
	qs = func(lo, hi int) {
		if lo >= hi {
			return
		}
		p := partition(lo, hi)
		qs(lo, p-1)
		qs(p+1, hi)
	}
	qs(0, len(nums)-1)
}

// MergeSort returns a new sorted slice.
func MergeSort(nums []int) []int {
	if len(nums) <= 1 {
		return append([]int(nil), nums...)
	}
	m := len(nums) / 2
	left := MergeSort(nums[:m])
	right := MergeSort(nums[m:])
	out := make([]int, 0, len(nums))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			out = append(out, left[i])
			i++
		} else {
			out = append(out, right[j])
			j++
		}
	}
	out = append(out, left[i:]...)
	out = append(out, right[j:]...)
	return out
}

// KthLargest using Quickselect (k=1 -> largest).
func KthLargest(nums []int, k int) (int, error) {
	if k < 1 || k > len(nums) {
		return 0, errors.New("k out of range")
	}
	target := len(nums) - k
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		p := partition(nums, lo, hi)
		if p == target {
			return nums[p], nil
		} else if p < target {
			lo = p + 1
		} else {
			hi = p - 1
		}
	}
	return 0, errors.New("unreachable")
}
func partition(a []int, lo, hi int) int {
	pivot := a[hi]
	i := lo
	for j := lo; j < hi; j++ {
		if a[j] < pivot {
			a[i], a[j] = a[j], a[i]
			i++
		}
	}
	a[i], a[hi] = a[hi], a[i]
	return i
}

// SearchRotated performs binary search in a rotated sorted array.
func SearchRotated(nums []int, target int) int {
	lo, hi := 0, len(nums)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if nums[mid] == target {
			return mid
		}
		if nums[lo] <= nums[mid] { // left sorted
			if nums[lo] <= target && target < nums[mid] {
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		} else { // right sorted
			if nums[mid] < target && target <= nums[hi] {
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
	}
	return -1
}

///////////////////////////////
// System Design-ish / Practical
///////////////////////////////

// LRU Cache (map + doubly linked list from container/list)
type LRUCache struct {
	capacity int
	ll       *list.List
	items    map[int]*list.Element
}
type kv struct{ k, v int }

func NewLRU(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		ll:       list.New(),
		items:    make(map[int]*list.Element, capacity),
	}
}
func (c *LRUCache) Get(key int) (int, bool) {
	if el, ok := c.items[key]; ok {
		c.ll.MoveToFront(el)
		return el.Value.(kv).v, true
	}
	return 0, false
}
func (c *LRUCache) Put(key, val int) {
	if el, ok := c.items[key]; ok {
		el.Value = kv{key, val}
		c.ll.MoveToFront(el)
		return
	}
	if c.ll.Len() == c.capacity {
		back := c.ll.Back()
		if back != nil {
			c.ll.Remove(back)
			delete(c.items, back.Value.(kv).k)
		}
	}
	el := c.ll.PushFront(kv{key, val})
	c.items[key] = el
}

// URL Shortener (toy) using auto-increment id + base62.
type Shortener struct {
	mu   sync.Mutex
	id   uint64
	urls map[string]string // code -> longURL
}
func NewShortener() *Shortener { return &Shortener{urls: map[string]string{}} }

const base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeBase62(n uint64) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append(b, base62[n%62])
		n /= 62
	}
	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func (s *Shortener) Shorten(longURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.id++
	code := encodeBase62(s.id)
	s.urls[code] = longURL
	return code
}
func (s *Shortener) Resolve(code string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.urls[code]
	return u, ok
}

// JSON parse/manipulate
type User struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Tags  []string `json:"tags,omitempty"`
	Email string   `json:"email"`
}

func ParseUserJSON(data []byte) (User, error) {
	var u User
	return u, json.Unmarshal(data, &u)
}
func DynamicJSONToMap(data []byte) (map[string]any, error) {
	var m map[string]any
	return m, json.Unmarshal(data, &m)
}

// Minimal REST API with net/http
func StartServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"echo": string(body)})
	})
	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		_ = srv.ListenAndServe()
	}()
	return srv
}

// File I/O: stream copy with buffering
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()
	bufIn := bufio.NewReader(in)
	bufOut := bufio.NewWriter(out)
	if _, err := io.Copy(bufOut, bufIn); err != nil {
		return err
	}
	return bufOut.Flush()
}

///////////////////////////////
// Math & Misc
///////////////////////////////

// Fibonacci iterative
func FibIter(n int) int {
	if n < 2 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// Fibonacci with memoization
func FibMemo(n int) int {
	memo := map[int]int{0: 0, 1: 1}
	var f func(int) int
	f = func(k int) int {
		if v, ok := memo[k]; ok {
			return v
		}
		memo[k] = f(k-1) + f(k-2)
		return memo[k]
	}
	return f(n)
}

// Primality test (deterministic for 32-bit range)
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	limit := int(math.Sqrt(float64(n)))
	for d := 3; d <= limit; d += 2 {
		if n%d == 0 {
			return false
		}
	}
	return true
}

// Sieve of Eratosthenes up to n.
func Sieve(n int) []int {
	if n < 2 {
		return nil
	}
	mark := make([]bool, n+1)
	var primes []int
	for p := 2; p <= n; p++ {
		if !mark[p] {
			primes = append(primes, p)
			for q := p * p; q <= n; q += p {
				mark[q] = true
			}
		}
	}
	return primes
}

// Concurrency-safe counter (mutex)
type SafeCounter struct {
	mu sync.Mutex
	n  int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

// Concurrency-safe counter (atomic)
type AtomicCounter struct{ n int64 }

func (c *AtomicCounter) Inc() { atomic.AddInt64(&c.n, 1) }
func (c *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&c.n)
}

// RotateMatrix90 rotates an N x N matrix clockwise in-place.
func RotateMatrix90(mat [][]int) {
	n := len(mat)
	// transpose
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			mat[i][j], mat[j][i] = mat[j][i], mat[i][j]
		}
	}
	// reverse each row
	for i := 0; i < n; i++ {
		for l, r := 0, n-1; l < r; l, r = l+1, r-1 {
			mat[i][l], mat[i][r] = mat[i][r], mat[i][l]
		}
	}
}

// NQueens returns solutions as slices of column indices (row -> col).
func NQueens(n int) [][]int {
	var res [][]int
	cols := make([]bool, n)
	d1 := make([]bool, 2*n) // r+c
	d2 := make([]bool, 2*n) // r-c+n
	sol := make([]int, n)
	var backtrack func(r int)
	backtrack = func(r int) {
		if r == n {
			cp := append([]int(nil), sol...)
			res = append(res, cp)
			return
		}
		for c := 0; c < n; c++ {
			if cols[c] || d1[r+c] || d2[r-c+n] {
				continue
			}
			cols[c], d1[r+c], d2[r-c+n] = true, true, true
			sol[r] = c
			backtrack(r + 1)
			cols[c], d1[r+c], d2[r-c+n] = false, false, false
		}
	}
	backtrack(0)
	return res
}

///////////////////////////////
// Tiny helpers for demos (optional)
///////////////////////////////

// Example demonstrating worker pool usage.
func ExampleWorkerPool() {
	jobs := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := WorkerPool(ctx, 4, jobs, func(ctx context.Context, x int) (int, error) {
		time.Sleep(10 * time.Millisecond)
		return x * x, nil
	})
	go func() {
		for i := 0; i < 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()
	for res := range out {
		if res.Err == nil {
			fmt.Println(res.Val)
		}
	}
}

// Simple text normalization (if you want stricter palindrome/anagram behavior)
func normalizeLetters(s string) string {
	var b strings.Builder
	for _, r := range s {
		if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') {
			b.WriteRune(r | 32) // lowercase ASCII
		}
	}
	return b.String()
}
