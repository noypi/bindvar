package bindvar

import (
	"fmt"
)

type IBindable interface {
	fmt.Stringer
	Ptr() uintptr
	Value() interface{}
}

var defaultBindVar BindVar

func init() {
	defaultBindVar.Init()
}

func Unbind(p interface{}) {
	defaultBindVar.Unbind(p)
}

type _bindInfo struct {
	bs []IBindable
	fn func()
}
type BindVar struct {
	fns map[interface{}]*_bindInfo
	// bindable by vars
	m map[uintptr]map[interface{}]struct{}
}

func (a *BindVar) Init() {
	a.fns = map[interface{}]*_bindInfo{}
	a.m = map[uintptr]map[interface{}]struct{}{}
}

func (a *BindVar) StrFmt(p *string, format string, os ...interface{}) {
	os2 := copyStrParams(os)
	fn := func() {
		updateBindableInParams(os, os2)
		*p = fmt.Sprintf(format, os2...)
	}

	a.bind(p, fn, getbindables(os))
}

func (a *BindVar) Str(p *string, os ...interface{}) {
	os2 := copyStrParams(os)
	fn := func() {
		updateBindableInParams(os, os2)
		*p = fmt.Sprint(os2...)
	}

	a.bind(p, fn, getbindables(os))
}

func (a *BindVar) Int(n *int, fn func() int, ns ...*Int) {
	fn0 := func() {
		*n = fn()
	}
	a.bind(n, fn0, Ints(ns).ToBindableArray())
}

func (a *BindVar) UnbindStr(p *String) {
	a.Unbind(p)
}

func (a *BindVar) UnbindInt(p *Int) {
	a.Unbind(p)
}

func (a *BindVar) Unbind(p interface{}) {
	info := a.fns[p]
	if nil == info {
		return
	}

	for _, b := range info.bs {
		m := a.m[b.Ptr()]
		delete(m, p)
		if 0 == len(m) {
			delete(a.m, b.Ptr())
		}
	}
	delete(a.fns, p)
}

func (a *BindVar) HaveUpdated(o IBindable) {
	ptr := o.Ptr()
	if nil == a.m {
		return
	}

	m := a.m[ptr]
	if nil == m {
		return
	}

	for vptr, _ := range m {
		info := a.fns[vptr]
		if nil != info {
			info.fn()
		}
	}
}

func (a *BindVar) bind(vptr interface{}, fn func(), bs []IBindable) {
	a.bindFn(vptr, fn, bs)
	a.bindParams(vptr, bs)
	// call once
	fn()
}

func (a *BindVar) bindFn(vptr interface{}, fn func(), bs []IBindable) {
	a.fns[vptr] = &_bindInfo{
		bs: bs,
		fn: fn,
	}
}

func (a *BindVar) bindParams(vptr interface{}, bs []IBindable) {
	for _, b := range bs {
		ptr := b.Ptr()
		m := a.m[ptr]
		if nil == m {
			m = map[interface{}]struct{}{}
			a.m[ptr] = m
		}
		m[vptr] = struct{}{}
	}

}

func getbindables(os []interface{}) (bs []IBindable) {
	for _, o := range os {
		switch v := o.(type) {
		case IBindable:
			bs = append(bs, v)
		}
	}
	return
}
