package bindvar

import (
	"strconv"
	"unsafe"
)

type Int struct {
	n int
	b *BindVar

	events map[uintptr][]func(oldVal, newVal int) bool
}

func BindInt(n *int, fn func() int, ns ...*Int) {
	defaultBindVar.Int(n, fn, ns...)
}

func UnbindInt(p *Int) {
	defaultBindVar.UnbindInt(p)
}

func (a *Int) Ptr() uintptr {
	return uintptr(unsafe.Pointer(a))
}

func NewInt() *Int {
	return defaultBindVar.NewInt()
}

func (a *BindVar) NewInt() *Int {
	return &Int{b: a, events: map[uintptr][]func(int, int) bool{}}
}

func (a *Int) Set(n int) {
	a.n = n
	a.notify(n)
}

func (a *Int) Bind(fn func() int, ns ...*Int) {
	fn0 := func() {
		a.Set(fn())
	}
	a.b.bind(a, fn0, Ints(ns).ToBindableArray())
}

func (a *Int) AddOnChange(handler func(oldVal, newVal int) bool) {
	ptr := a.Ptr()
	arr, _ := a.events[ptr]
	a.events[ptr] = append(arr, handler)
}

func (a *Int) String() string     { return strconv.Itoa(a.n) }
func (a *Int) Value() interface{} { return a.n }

func (a *Int) notify(newVal int) {
	ptr := a.Ptr()
	for _, fn := range a.events[ptr] {
		if !fn(a.n, newVal) {
			return
		}
	}
	a.b.HaveUpdated(a)
}

type Ints []*Int

func (ns Ints) ToBindableArray() []IBindable {
	bs := make([]IBindable, len(ns))
	for i, vn := range ns {
		bs[i] = vn
	}
	return bs
}
