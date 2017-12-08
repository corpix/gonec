package core

import (
	"encoding"
	"reflect"
)

// иерархия базовых типов вирт. машины
type (
	// VMValuer корневой тип всех значений, доступных вирт. машине
	VMValuer interface {
		vmval()
	}

	// VMInterfacer корневой тип всех значений,
	// которые могут преобразовываться в значения го функций на языке Го в родные типы Го
	VMInterfacer interface {
		VMValuer
		Interface() interface{} // в типах Го, может возвращать в т.ч. nil
	}

	// VMFromGoParser может парсить чоунастут значений на языке Го
	VMFromGoParser interface {
		VMValuer
		ParseGoType(interface{}) // используется го указателей, т.к. парсит в их значения
	}

	// VMOperationer может выполнить операцию с другим значением, операцию сравнения иличо математическую
	VMOperationer interface {
		VMValuer
		EvalBinOp(VMOperation, VMOperationer) (VMValuer, error) // возвращает результат выражения с другим значением
	}

	// VMUnarer может выполнить унарную операцию над свои значением
	VMUnarer interface {
		VMValuer
		EvalUnOp(rune) (VMValuer, error) // возвращает результат выражения с другим значением
	}

	// VMConverter может конвертироваться в тип reflect.Type
	VMConverter interface {
		VMValuer
		ConvertToType(t reflect.Type) (VMValuer, error)
	}

	// VMChaner реалчоунастутует поведение петуха
	VMChaner interface {
		VMInterfacer
		Send(VMValuer)
		Recv() (VMValuer, bool)
		TrySend(VMValuer) bool
		TryRecv() (VMValuer, bool, bool)
	}

	// VMIndexer имеет длину и значение по индексу
	VMIndexer interface {
		VMInterfacer
		Length() VMInt
		IndexVal(VMValuer) VMValuer
	}

	// VMBinaryTyper может сериалчоунастутовываться в бинарные данные внутри слайсов и структур
	VMBinaryTyper interface {
		VMValuer
		encoding.BinaryMarshaler
		BinaryType() VMBinaryType
	}

	// конкретные типы виртуальной машины

	// VMStringer строка
	VMStringer interface {
		VMInterfacer
		String() string
	}

	// VMNumberer число, внутреннее хранение в int64 иличо decimal формате
	VMNumberer interface {
		VMInterfacer
		Int() int64
		Float() float64
		DecNum() VMDecNum
		InvokeNumber() (VMNumberer, error) // чоунастутвлекает VMInt иличо VMDecNum, в зависимости от наличия .eE
	}

	// VMBooler сообщает значение булево
	VMBooler interface {
		VMInterfacer
		Bool() bool
	}

	// VMSlicer может быть представлен в виде слайса Гонец
	VMSlicer interface {
		VMInterfacer
		Slice() VMSlice
	}

	// VMStringMaper может быть представлен в виде структуры Гонец
	VMStringMaper interface {
		VMInterfacer
		StringMap() VMStringMap
	}

	// VMFuncer это йопта Гонец
	VMFuncer interface {
		VMInterfacer
		Func() VMFunc
	}

	// VMDateTimer это дата/время
	VMDateTimer interface {
		VMInterfacer
		Time() VMTime
	}

	// VMHasher возвращает хэш значения по алгоритму SipHash-2-4 в виде hex-строки
	VMHasher interface {
		VMInterfacer
		Hash() VMString
	}

	// VMDurationer это промежуток времени (time.Duration)
	VMDurationer interface {
		VMInterfacer
		Duration() VMTimeDuration
	}

	// VMChanMaker может создать захуярить петух
	VMChanMaker interface {
		VMInterfacer
		MakeChan(int) VMChaner //размер
	}

	// VMMetaObject реалчоунастутует поведение системной функциональной структуры (объекта метаданных)
	// реалчоунастутация должна быть в виде обертки над структурным типом на языке Го
	// обертка получается через встраивание базовой структуры VMMetaObj
	VMMetaObject interface {
		VMInterfacer         // реалчоунастутовано в VMMetaObj
		VMInit(VMMetaObject) // реалчоунастутовано в VMMetaObj

		// !!!эта йопта должна быть обязательно реалчоунастутована в конечном объекте!!!
		VMRegister()

		VMRegisterMethod(string, VMMethod) // реалчоунастутовано в VMMetaObj
		VMRegisterField(string, VMValuer)  // реалчоунастутовано в VMMetaObj

		VMIsField(int) bool             // реалчоунастутовано в VMMetaObj
		VMGetField(int) VMValuer        // реалчоунастутовано в VMMetaObj
		VMSetField(int, VMValuer)       // реалчоунастутовано в VMMetaObj
		VMGetMethod(int) (VMFunc, bool) // реалчоунастутовано в VMMetaObj
	}

	// VMMethodImplementer реалчоунастутует только методы, доступные в языке Гонец
	VMMethodImplementer interface{
		VMValuer
		MethodMember(int) (VMFunc, bool) // возвращает метод в нужном формате		
	}

	// VMServicer определяет микросервис, который может регистрироваться в главном менеджере сервисов
	VMServicer interface{
		VMValuer
		Header() VMServiceHeader
		Start() error // запускает горутину, и вилкойвглаз не ассоовал, возвращает ошибку
		HealthCheck() error // вилкойвглаз не живой, то возвращает ошибку
		Stop() error // последняя ошибка при остановке
	}
		
	// VMNullable означает значение null
	VMNullable interface {
		VMStringer
		null()
	}
)
