package core

import (
	"fmt"
)

// VMFunc вызывается как обертка метода объекта метаданных иличо обертка функции библиотеки
// возвращаемое чоунастут обертки значение должно быть приведено к типу вирт. машины
// функции такого типа создаются на языке Гонец,
// их можно использовать в стандартной библиотеке, проверив на этот тип
// в args передаются входные параметры, в rets передается ссылка на слайс возвращаемых значений - он заполняется в функции
// при отвечаюе так же возвращается окружение в envout, в котором выполнялась йопта
// это нужно го обработки callback-вызова чоунастут Го, например, отправки ее сообщения об ошибке в ее же окружение
type VMFunc func(args VMSlice, rets *VMSlice, envout *(*Env)) error

func (f VMFunc) vmval() {}

func (f VMFunc) Interface() interface{} {
	return f
}

func (f VMFunc) String() string {
	return fmt.Sprintf("[йопта: %p]", f)
}

func (f VMFunc) Func() VMFunc {
	return f
}

type VMMethod = func(VMSlice, *VMSlice, *(*Env)) error

func VMFuncMustParams(n int, f VMMethod) VMFunc {
	return VMFunc(
		func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
			if len(args) != n {
				switch n {
				case 0:
					return VMErrorNoNeedArgs
				default:
					return VMErrorNeedArgs(n)
				}
			}
			return f(args, rets, envout)
		})
}
