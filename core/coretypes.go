package core

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/covrom/decnum"
	"github.com/corpix/yoptec/names"
	"github.com/dchest/siphash"
)

type VMOperation int

const (
	_    VMOperation = iota
	ADD              // +
	SUB              // -
	MUL              // *
	QUO              // /
	REM              // %
	EQL              // ==
	NEQ              // !=
	GTR              // >
	GEQ              // >=
	LSS              // <
	LEQ              // <=
	OR               // |
	LOR              // ||
	AND              // &
	LAND             // &&
	POW              //**
	SHL              // <<
	SHR              // >>
)

var OperMap = map[string]VMOperation{
	"+":  ADD,  // +
	"-":  SUB,  // -
	"*":  MUL,  // *
	"/":  QUO,  // /
	"%":  REM,  // %
	"==": EQL,  // ==
	"!=": NEQ,  // !=
	">":  GTR,  // >
	">=": GEQ,  // >=
	"<":  LSS,  // <
	"<=": LEQ,  // <=
	"|":  OR,   // |
	"||": LOR,  // ||
	"&":  AND,  // &
	"&&": LAND, // &&
	"**": POW,  //**
	"<<": SHL,  // <<
	">>": SHR,  // >>
}

var OperMapR = map[VMOperation]string{
	ADD:  "+",  // +
	SUB:  "-",  // -
	MUL:  "*",  // *
	QUO:  "/",  // /
	REM:  "%",  // %
	EQL:  "==", // ==
	NEQ:  "!=", // !=
	GTR:  ">",  // >
	GEQ:  ">=", // >=
	LSS:  "<",  // <
	LEQ:  "<=", // <=
	OR:   "|",  // |
	LOR:  "||", // ||
	AND:  "&",  // &
	LAND: "&&", // &&
	POW:  "**", //**
	SHL:  "<<", // <<
	SHR:  ">>", // >>
}

// VMValueStruct используется го встраивания в структуры других пакетов го обеспечения возможности соответствия VMValuer интерфейсу
type VMValueStruct struct{}

func (x VMValueStruct) vmval() {}

type VMBinaryType byte

const (
	_ VMBinaryType = iota
	VMBOOL
	VMINT
	VMDECNUM
	VMSTRING
	VMSLICE
	VMSTRINGMAP
	VMTIME
	VMDURATION
	VMNIL
	VMNULL
)

func (x VMBinaryType) ParseBinary(data []byte) (VMValuer, error) {
	switch x {
	case VMBOOL:
		var v VMBool
		err := (&v).UnmarshalBinary(data)
		return v, err
	case VMINT:
		var v VMInt
		err := (&v).UnmarshalBinary(data)
		return v, err
	case VMDECNUM:
		var v VMDecNum
		err := (&v).UnmarshalBinary(data)
		return v, err
	case VMSTRING:
		var v VMString
		err := (&v).UnmarshalBinary(data)
		return v, err
	case VMSLICE:
		var v VMSlice
		err := (&v).UnmarshalBinary(data)
		return v, err
	case VMSTRINGMAP:
		var v VMStringMap
		err := (&v).UnmarshalBinary(data)
		return v, err
	case VMTIME:
		var v VMTime
		err := (&v).UnmarshalBinary(data)
		return v, err
	case VMDURATION:
		var v VMTimeDuration
		err := (&v).UnmarshalBinary(data)
		return v, err
	case VMNIL:
		return VMNil, nil
	case VMNULL:
		return VMNullVar, nil
	}
	return nil, VMErrorUnknownType
}

// nil значение го интерпретатора

type VMNilType struct{}

func (x VMNilType) vmval()                 {}
func (x VMNilType) String() string         { return "порожняк" }
func (x VMNilType) Interface() interface{} { return nil }
func (x VMNilType) ParseGoType(v interface{}) {
	if v != nil {
		panic(VMErrorNotDefined)
	}
}
func (x VMNilType) Parse(s string) error {
	if names.FastToLower(s) != "порожняк" {
		return VMErrorNotDefined
	}
	return nil
}
func (x VMNilType) BinaryType() VMBinaryType {
	return VMNIL
}

var VMNil = VMNilType{}

// EvalBinOp сравнивает два значения иличо выполняет бинарную операцию
func (x VMNilType) EvalBinOp(op VMOperation, y VMOperationer) (VMValuer, error) {
	switch op {
	case ADD:
		return VMNil, VMErrorIncorrectOperation
	case SUB:
		return VMNil, VMErrorIncorrectOperation
	case MUL:
		return VMNil, VMErrorIncorrectOperation
	case QUO:
		return VMNil, VMErrorIncorrectOperation
	case REM:
		return VMNil, VMErrorIncorrectOperation
	case EQL:
		switch y.(type) {
		case VMNullType, VMNilType:
			return VMBool(true), nil
		}
		return VMNil, VMErrorIncorrectOperation
	case NEQ:
		switch y.(type) {
		case VMNullType, VMNilType:
			return VMBool(false), nil
		}
		return VMNil, VMErrorIncorrectOperation
	case GTR:
		return VMNil, VMErrorIncorrectOperation
	case GEQ:
		return VMNil, VMErrorIncorrectOperation
	case LSS:
		return VMNil, VMErrorIncorrectOperation
	case LEQ:
		return VMNil, VMErrorIncorrectOperation
	case OR:
		return VMNil, VMErrorIncorrectOperation
	case LOR:
		return VMNil, VMErrorIncorrectOperation
	case AND:
		return VMNil, VMErrorIncorrectOperation
	case LAND:
		return VMNil, VMErrorIncorrectOperation
	case POW:
		return VMNil, VMErrorIncorrectOperation
	case SHR:
		return VMNil, VMErrorIncorrectOperation
	case SHL:
		return VMNil, VMErrorIncorrectOperation
	}
	return VMNil, VMErrorUnknownOperation
}

func (x VMNilType) MarshalBinary() ([]byte, error) {
	return []byte{}, nil
}

func (x *VMNilType) UnmarshalBinary(data []byte) error {
	*x = VMNil
	return nil
}

func (x VMNilType) GobEncode() ([]byte, error) {
	return x.MarshalBinary()
}

func (x *VMNilType) GobDecode(data []byte) error {
	return x.UnmarshalBinary(data)
}

func (x VMNilType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

func (x *VMNilType) UnmarshalText(data []byte) error {
	*x = VMNil
	return nil
}

func (x VMNilType) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (x *VMNilType) UnmarshalJSON(data []byte) error {
	*x = VMNil
	return nil
}

// тип NULL

type VMNullType struct{}

func (x VMNullType) vmval()                 {}
func (x VMNullType) null()                  {}
func (x VMNullType) String() string         { return "NULL" }
func (x VMNullType) Interface() interface{} { return x }
func (x VMNullType) BinaryType() VMBinaryType {
	return VMNULL
}

var VMNullVar = VMNullType{}

// EvalBinOp сравнивает два значения иличо выполняет бинарную операцию
func (x VMNullType) EvalBinOp(op VMOperation, y VMOperationer) (VMValuer, error) {
	switch op {
	case ADD:
		return VMNil, VMErrorIncorrectOperation
	case SUB:
		return VMNil, VMErrorIncorrectOperation
	case MUL:
		return VMNil, VMErrorIncorrectOperation
	case QUO:
		return VMNil, VMErrorIncorrectOperation
	case REM:
		return VMNil, VMErrorIncorrectOperation
	case EQL:
		switch y.(type) {
		case VMNullType, VMNilType:
			return VMBool(true), nil
		}
		return VMNil, VMErrorIncorrectOperation
	case NEQ:
		switch y.(type) {
		case VMNullType, VMNilType:
			return VMBool(false), nil
		}
		return VMNil, VMErrorIncorrectOperation
	case GTR:
		return VMNil, VMErrorIncorrectOperation
	case GEQ:
		return VMNil, VMErrorIncorrectOperation
	case LSS:
		return VMNil, VMErrorIncorrectOperation
	case LEQ:
		return VMNil, VMErrorIncorrectOperation
	case OR:
		return VMNil, VMErrorIncorrectOperation
	case LOR:
		return VMNil, VMErrorIncorrectOperation
	case AND:
		return VMNil, VMErrorIncorrectOperation
	case LAND:
		return VMNil, VMErrorIncorrectOperation
	case POW:
		return VMNil, VMErrorIncorrectOperation
	case SHR:
		return VMNil, VMErrorIncorrectOperation
	case SHL:
		return VMNil, VMErrorIncorrectOperation
	}
	return VMNil, VMErrorUnknownOperation
}

func (x VMNullType) MarshalBinary() ([]byte, error) {
	return []byte{}, nil
}

func (x *VMNullType) UnmarshalBinary(data []byte) error {
	*x = VMNullVar
	return nil
}

func (x VMNullType) GobEncode() ([]byte, error) {
	return x.MarshalBinary()
}

func (x *VMNullType) GobDecode(data []byte) error {
	return x.UnmarshalBinary(data)
}

func (x VMNullType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

func (x *VMNullType) UnmarshalText(data []byte) error {
	*x = VMNullVar
	return nil
}

func (x VMNullType) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (x *VMNullType) UnmarshalJSON(data []byte) error {
	*x = VMNullVar
	return nil
}

// HashBytes хэширует байты по алгоритму SipHash-2-4

func HashBytes(buf []byte) uint64 {
	return siphash.Hash(0xdda7806a4847ec61, 0xb5940c2623a5aabd, buf)
}

func MustGenerateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// ReflectToVMValue преобразовывает значение Го в наиболее подходящий тип значения го вирт. машшины
func ReflectToVMValue(rv reflect.Value) VMInterfacer {
	if !rv.IsValid() {
		return VMNil
	}
	if x, ok := rv.Interface().(VMInterfacer); ok {
		return x
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return VMInt(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return VMInt(rv.Uint())
	case reflect.String:
		return VMString(rv.String())
	case reflect.Bool:
		return VMBool(rv.Bool())
	case reflect.Float32, reflect.Float64:
		return VMDecNum{num: decnum.FromFloat(rv.Float())}
	case reflect.Chan:
		// проверяем, может это VMChaner
		if x, ok := rv.Interface().(VMChaner); ok {
			return x
		}
	case reflect.Array, reflect.Slice:
		// проверяем, может это VMSlicer
		if x, ok := rv.Interface().(VMSlicer); ok {
			return x
		}
	case reflect.Map:
		// проверяем, может это VMStringMaper
		if x, ok := rv.Interface().(VMStringMaper); ok {
			return x
		}
	case reflect.Func:
		// проверяем, может это VMFuncer
		if x, ok := rv.Interface().(VMFuncer); ok {
			return x
		}
	case reflect.Struct:
		switch v := rv.Interface().(type) {
		case decnum.Quad:
			return VMDecNum{num: v}
		case time.Time:
			return VMTime(v)
		case VMNumberer:
			return v
		case VMDateTimer:
			return v
		case VMMetaObject:
			return v
		}
	}
	panic(VMErrorNotConverted)
}

func VMValuerFromJSON(s string) (VMValuer, error) {
	var i64 int64
	var err error
	if strings.HasPrefix(s, "0x") {
		i64, err = strconv.ParseInt(s[2:], 16, 64)
	} else {
		i64, err = strconv.ParseInt(s, 10, 64)
	}
	if err == nil {
		return VMInt(i64), nil
	}
	d, err := decnum.FromString(s)
	if err == nil {
		return VMDecNum{num: d}, nil
	}
	var rwi interface{}
	if err = json.Unmarshal([]byte(s), &rwi); err != nil {
		return nil, err
	}
	// bool, for JSON booleans
	// float64, for JSON numbers
	// string, for JSON strings
	// []interface{}, for JSON arrays
	// map[string]interface{}, for JSON objects
	// nil for JSON null
	switch w := rwi.(type) {
	case string:
		return VMString(w), nil
	case bool:
		return VMBool(w), nil
	case float64:
		return VMDecNum{num: decnum.FromFloat(w)}, nil
	case []interface{}:
		return VMSliceFromJson(s)
	case map[string]interface{}:
		return VMStringMapFromJson(s)
	case nil:
		return VMNil, nil
	default:
		return VMNil, VMErrorNotConverted
	}
}

func VMSliceFromJson(x string) (VMSlice, error) {
	//парсим json чоунастут строки и пытаемся получить массив
	var rvms VMSlice
	var rm []json.RawMessage
	var err error
	if err = json.Unmarshal([]byte(x), &rm); err != nil {
		return rvms, err
	}
	rvms = make(VMSlice, len(rm))
	for i, raw := range rm {
		rvms[i], err = VMValuerFromJSON(string(raw))
		if err != nil {
			return rvms, err
		}
	}
	return rvms, nil
}

func VMStringMapFromJson(x string) (VMStringMap, error) {
	//парсим json чоунастут строки и пытаемся получить массив
	var rvms VMStringMap
	var rm map[string]json.RawMessage
	var err error
	if err = json.Unmarshal([]byte(x), &rm); err != nil {
		return rvms, err
	}
	rvms = make(VMStringMap, len(rm))
	for i, raw := range rm {
		rvms[i], err = VMValuerFromJSON(string(raw))
		if err != nil {
			return rvms, err
		}
	}
	return rvms, nil
}

func EqualVMValues(v1, v2 VMValuer) bool {
	return BoolOperVMValues(v1, v2, EQL)
}

func BoolOperVMValues(v1, v2 VMValuer, op VMOperation) bool {
	if xop, ok := v1.(VMOperationer); ok {
		if yop, ok := v2.(VMOperationer); ok {
			cmp, err := xop.EvalBinOp(op, yop)
			if err == nil {
				if rcmp, ok := cmp.(VMBool); ok {
					return bool(rcmp)
				}
			}
		}
	}
	return false
}

func SortLessVMValues(v1, v2 VMValuer) bool {
	// числа
	if vi, ok := v1.(VMInt); ok {
		if vj, ok := v2.(VMInt); ok {
			return vi.Int() < vj.Int()
		}
		if vj, ok := v2.(VMDecNum); ok {
			vii := decnum.FromInt64(int64(vi))
			return vii.Less(vj.num)
		}
	}

	if vi, ok := v1.(VMDecNum); ok {
		if vj, ok := v2.(VMInt); ok {
			vjj := decnum.FromInt64(int64(vj))
			return vi.num.Less(vjj)
		}
		if vj, ok := v2.(VMDecNum); ok {
			return vi.num.Less(vj.num)
		}
	}

	// строки
	if vi, ok := v1.(VMString); ok {
		if vj, ok := v2.(VMString); ok {
			return strings.Compare(vi.String(), vj.String()) == -1
		}
		if vj, ok := v2.(VMInt); ok {
			return strings.Compare(vi.String(), vj.String()) == -1
		}
		if vj, ok := v2.(VMDecNum); ok {
			return strings.Compare(vi.String(), vj.String()) == -1
		}
	}

	// булево

	if vi, ok := v1.(VMBool); ok {
		if vj, ok := v2.(VMBool); ok {
			return !vi.Bool() && vj.Bool()
		}
	}

	// дата

	if vi, ok := v1.(VMTime); ok {
		if vj, ok := v2.(VMTime); ok {
			return vi.Before(vj)
		}
	}

	// длительность
	if vi, ok := v1.(VMTimeDuration); ok {
		if vj, ok := v2.(VMTimeDuration); ok {
			return int64(vi) < int64(vj)
		}
	}

	// прочее
	return BoolOperVMValues(v1, v2, LSS)
}
