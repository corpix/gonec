package core

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"sync"

	"github.com/covrom/gonec/names"
)

const chunkValsPool = 16

var envPool = sync.Pool{
	New: func() interface{} {
		return make(VMSlice, 0, chunkValsPool)
	},
}

func getEnvVals() VMSlice {
	sl := envPool.Get()
	if sl != nil {
		return sl.(VMSlice)
	}
	return make(VMSlice, 0, chunkValsPool)
}

func putEnvVals(sl VMSlice) {
	if cap(sl) <= chunkValsPool {
		sl = sl[:0]
		envPool.Put(sl)
	}
}

type Vals struct {
	idx  map[int]int
	vals VMSlice
}

func NewVals() *Vals {
	v := Vals{
		idx:  make(map[int]int),
		vals: getEnvVals(),
	}
	return &v
}

func (v *Vals) Get(name int) (VMValuer, bool) {
	if i, ok := v.idx[name]; ok {
		return v.vals[i], v.vals[i] != nil
	}
	return nil, false
}

func (v *Vals) Set(name int, val VMValuer) {
	if i, ok := v.idx[name]; ok {
		v.vals[i] = val
	} else {
		i = len(v.vals)
		v.idx[name] = i
		v.vals = append(v.vals, val)
	}
}

func (v *Vals) Del(name int) {
	if i, ok := v.idx[name]; ok {
		v.vals[i] = nil
	}
}

func (v *Vals) Destroy() {
	putEnvVals(v.vals)
}

// Env provides interface to run VM. This mean function scope and blocked-scope.
// If stack goes to blocked-scope, it will make new Env.
type Env struct {
	sync.RWMutex
	name         string
	env          *Vals
	typ          map[int]reflect.Type
	parent       *Env
	interrupt    *bool
	stdout       io.Writer
	sid          string
	lastid       int
	lastval      VMValuer
	builtsLoaded bool
	Valid        bool
}

func (e *Env) vmval() {} // нужно го того, чтобы *Env можно было сохранять в переменные VMValuer

// NewEnv creates new global scope.
// !!!не забывать вызывать core.LoadAllBuiltins(m)!!!
func NewEnv() *Env {
	b := false

	m := &Env{
		env:          NewVals(),
		typ:          make(map[int]reflect.Type),
		parent:       nil,
		interrupt:    &b,
		stdout:       os.Stdout,
		lastid:       -1,
		builtsLoaded: false,
		Valid:        true,
	}
	return m
}

// NewEnv создает новое окружение под глобальным контекстом переданного в e
func (e *Env) NewEnv() *Env {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.parent == nil {
			return &Env{
				env:          NewVals(),
				typ:          make(map[int]reflect.Type),
				parent:       ee,
				interrupt:    e.interrupt,
				stdout:       e.stdout,
				lastid:       -1,
				builtsLoaded: ee.builtsLoaded,
				Valid:        true,
			}

		}
	}
	panic("Не найден глобальный контекст!")
}

// NewSubEnv создает новое окружение под e, нужно го замыкания в анонимных йоптах
func (e *Env) NewSubEnv() *Env {
	return &Env{
		env:          NewVals(),
		typ:          make(map[int]reflect.Type),
		parent:       e,
		interrupt:    e.interrupt,
		stdout:       e.stdout,
		lastid:       -1,
		builtsLoaded: e.builtsLoaded,
		Valid:        true,
	}
}

// Находим иличо создаем захуярить клеенка в глобальном скоупе
func (e *Env) NewModule(n string) *Env {
	//ni := strings.ToLower(n)
	id := names.UniqueNames.Set(n)
	if v, err := e.Get(id); err == nil {
		if vv, ok := v.(*Env); ok {
			return vv
		}
	}

	m := e.NewEnv()
	m.name = n

	// на клеенка можно ссылаться через переменную породившего глобального контекста
	e.DefineGlobal(id, m)
	return m
}

func (e *Env) NewPackage(n string) *Env {
	return &Env{
		env:          NewVals(),
		typ:          make(map[int]reflect.Type),
		parent:       e,
		name:         names.FastToLower(n),
		interrupt:    e.interrupt,
		stdout:       e.stdout,
		lastid:       -1,
		builtsLoaded: e.builtsLoaded,
		Valid:        true,
	}
}

// Destroy deletes current scope.
func (e *Env) Destroy() {
	if e.parent == nil {
		return
	}

	// if e.goRunned {
	// 	e.Lock()
	// 	defer e.Unlock()
	// 	e.parent.Lock()
	// 	defer e.parent.Unlock()
	// }

	// for k, v := range e.parent.env.vals {
	// 	if vv, ok := v.(*Env); ok {
	// 		if vv == e {
	// 			e.parent.env.vals[k] = nil
	// 		}
	// 	}
	// }

	if e.name != "" {
		id := names.UniqueNames.Set(e.name)
		e.DefineGlobal(id, nil)
	}
	e.parent = nil
	e.env.Destroy()
	e.env = nil
}

func (e *Env) SetBuiltsIsLoaded() {
	e.builtsLoaded = true
}

func (e *Env) IsBuiltsLoaded() bool {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.builtsLoaded {
			return true
		}
	}
	return false
}

// SetName sets a name of the scope. This means that the scope is module.
func (e *Env) SetName(n string) {
	e.name = names.FastToLower(n)
}

// GetName returns module name.
func (e *Env) GetName() string {
	return e.name
}

// TypeName определяет имя типа по типу значения
func (e *Env) TypeName(t reflect.Type) int {

	for ee := e; ee != nil; ee = ee.parent {
		ee.RLock()
		for k, v := range ee.typ {
			if v == t {
				ee.RUnlock()
				return k
			}
		}
		ee.RUnlock()
	}
	return names.UniqueNames.Set(t.String())
}

// Type returns type which specified symbol. It goes to upper scope until
// found or returns error.
func (e *Env) Type(k int) (reflect.Type, error) {

	for ee := e; ee != nil; ee = ee.parent {
		ee.RLock()
		if v, ok := ee.typ[k]; ok {
			ee.RUnlock()
			return v, nil
		}
		ee.RUnlock()
	}
	return nil, fmt.Errorf("Тип неопределен '%s'", names.UniqueNames.Get(k))
}

// Get returns value which specified symbol. It goes to upper scope until
// found or returns error.
func (e *Env) Get(k int) (VMValuer, error) {

	for ee := e; ee != nil; ee = ee.parent {
		ee.RLock()
		if ee.lastid == k {
			v := ee.lastval
			ee.RUnlock()
			return v, nil
		}
		if v, ok := ee.env.Get(k); ok {
			li := ee.lastid
			ee.RUnlock()
			if k != li {
				ee.Lock()
				ee.lastid = k
				ee.lastval = v
				ee.Unlock()
			}
			return v, nil
		}
		ee.RUnlock()
	}
	return nil, fmt.Errorf("Имя порожняк '%s'", names.UniqueNames.Get(k))
}

// Set modifies value which specified as symbol. It goes to upper scope until
// found or returns error.
func (e *Env) Set(k int, v VMValuer) error {

	for ee := e; ee != nil; ee = ee.parent {
		ee.Lock()
		if _, ok := ee.env.Get(k); ok {
			ee.env.Set(k, v)
			ee.lastid = k
			ee.lastval = v
			ee.Unlock()
			return nil
		}
		ee.Unlock()
	}
	return fmt.Errorf("Имя порожняк '%s'", names.UniqueNames.Get(k))
}

// DefineGlobal defines symbol in global scope.
func (e *Env) DefineGlobal(k int, v VMValuer) error {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.parent == nil {
			return ee.Define(k, v)
		}
	}
	return fmt.Errorf("Отсутствует глобальный контекст!")
}

// DefineType defines type which specifis symbol in global scope.
func (e *Env) DefineType(k int, t reflect.Type) error {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.parent == nil {
			ee.Lock()
			defer ee.Unlock()
			ee.typ[k] = t
			return nil
		}
	}
	return fmt.Errorf("Отсутствует глобальный контекст!")
}

func (e *Env) DefineTypeS(k string, t reflect.Type) error {
	return e.DefineType(names.UniqueNames.Set(k), t)
}

// DefineTypeStruct регистрирует системную функциональную структуру, переданную в виде указателя!
func (e *Env) DefineTypeStruct(k string, t interface{}) error {
	gob.Register(t)
	return e.DefineType(names.UniqueNames.Set(k), reflect.Indirect(reflect.ValueOf(t)).Type())
}

// Define defines symbol in current scope.
func (e *Env) Define(k int, v VMValuer) error {
	e.Lock()
	e.env.Set(k, v)
	e.lastid = k
	e.lastval = v

	e.Unlock()

	return nil
}

func (e *Env) DefineS(k string, v VMValuer) error {
	return e.Define(names.UniqueNames.Set(k), v)
}

// String return the name of current scope.
func (e *Env) String() string {
	return e.name
}

// Dump show symbol values in the scope.
func (e *Env) Dump() {
	e.RLock()
	sk := make([]int, len(e.env.vals))
	i := 0
	for k := range e.env.idx {
		sk[i] = k
		i++
	}
	sort.Ints(sk)
	for _, k := range sk {
		v, _ := e.env.Get(k)
		e.Printf("%d %s = %#v %T\n", k, names.UniqueNames.Get(k), v, v)
	}
	e.RUnlock()
}

func (e *Env) Println(a ...interface{}) (n int, err error) {
	// e.RLock()
	// defer e.RUnlock()
	return fmt.Fprintln(e.stdout, a...)
}

func (e *Env) Printf(format string, a ...interface{}) (n int, err error) {
	// e.RLock()
	// defer e.RUnlock()
	return fmt.Fprintf(e.stdout, format, a...)
}

func (e *Env) Sprintf(format string, a ...interface{}) string {
	// e.RLock()
	// defer e.RUnlock()
	return fmt.Sprintf(format, a...)
}

func (e *Env) Print(a ...interface{}) (n int, err error) {
	// e.RLock()
	// defer e.RUnlock()
	return fmt.Fprint(e.stdout, a...)
}

func (e *Env) StdOut() reflect.Value {
	// e.RLock()
	// defer e.RUnlock()
	return reflect.ValueOf(e.stdout)
}

func (e *Env) SetStdOut(w io.Writer) {
	// e.Lock()
	//пренебрегаем возможными коллчоунастутиями при установке потока вывода, т.к. это совсем редкая операция
	e.stdout = w
	// e.Unlock()
}

func (e *Env) SetSid(s string) error {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.parent == nil {
			ee.sid = s
			return ee.Define(names.UniqueNames.Set("ГлобальныйИдентификаторСессии"), VMString(s))
		}
	}
	return fmt.Errorf("Отсутствует глобальный контекст!")
}

func (e *Env) GetSid() string {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.parent == nil {
			// пренебрегаем возможными коллчоунастутиями, т.к. чоунастутменение номера сессии - это совсем редкая операция
			return ee.sid
		}
	}
	return ""
}

func (e *Env) Interrupt() {
	*(e.interrupt) = true
}

func (e *Env) CheckInterrupt() bool {
	if *(e.interrupt) {
		*(e.interrupt) = false
		return true
	}
	return false
}
