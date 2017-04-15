package bindvar

import (
	"fmt"
	"unsafe"
)

type String struct {
	s string
	b *BindVar

	events map[uintptr][]func(oldVal, newVal string) bool
}

func BindStrFmt(p *string, format string, os ...interface{}) {
	defaultBindVar.StrFmt(p, format, os...)
}

func BindStr(p *string, os ...interface{}) {
	defaultBindVar.Str(p, os...)
}

func UnbindStr(p *String) {
	defaultBindVar.UnbindStr(p)
}

func (a *String) Ptr() uintptr {
	return uintptr(unsafe.Pointer(a))
}

func NewString() *String {
	return defaultBindVar.NewString()
}

func (a *BindVar) NewString() *String {
	return &String{b: a, events: map[uintptr][]func(string, string) bool{}}
}

func (a *String) Set(s string) {
	a.s = s
	a.notify(s)
}

func (a *String) Bind(os ...interface{}) {
	os2 := copyStrParams(os)
	fn := func() {
		updateBindableInParams(os, os2)
		a.Set(fmt.Sprint(os2...))
	}

	a.b.bind(a, fn, getbindables(os))
}

func (a *String) BindFmt(format string, os ...interface{}) {
	os2 := copyStrParams(os)
	fn := func() {
		updateBindableInParams(os, os2)
		a.Set(fmt.Sprintf(format, os2...))
	}

	a.b.bind(a, fn, getbindables(os))
}

func (a *String) AddOnChange(handler func(oldVal, newVal string) bool) {
	ptr := a.Ptr()
	arr, _ := a.events[ptr]
	a.events[ptr] = append(arr, handler)
}

func (a *String) String() string     { return a.s }
func (a *String) Value() interface{} { return a.s }

func (a *String) notify(newVal string) {
	ptr := a.Ptr()
	for _, fn := range a.events[ptr] {
		if !fn(a.s, newVal) {
			return
		}
	}
	a.b.HaveUpdated(a)
}

func updateBindableInParams(base, os2 []interface{}) {
	for i, o := range base {
		if v, ok := o.(IBindable); ok {
			os2[i] = v.Value()
		}
	}
}

func copyStrParams(os []interface{}) []interface{} {
	os2 := make([]interface{}, len(os))
	copy(os2, os)
	return os2
}
