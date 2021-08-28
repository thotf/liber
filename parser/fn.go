package parser

//函数、子过程、编译单元都有一个闭包

const(
	//INIT_INSTRUCTION_SIZE = 128
	INIT_CONSTANTABLE_SIZE = 128
)

type fn struct {
	//  函数指令流
	instRuction		[]uint32
	//常量表
	constantTable	[]Value
	//所属模块
	module			*Module
	//本函数所覆盖的upvalue数量
	UpvalueQuantity  int
	//函数期望参数, 如果是一个函数闭包才有
	NumberParameters int
	//返回值数量
	ReturnValue      int8

	//所需要最大栈空间为 maxStackNum+1 从0开始计算
	maxStackNum		 uint
}

func NewFn() *fn {
	return &fn{
		instRuction:      make([]uint32,0),
		constantTable:    make([]Value,0,INIT_CONSTANTABLE_SIZE),
		module:           nil,
		UpvalueQuantity:  0,
		NumberParameters: 0,
		maxStackNum:0,
		ReturnValue:0,
	}
}
//存储一个常量、字面量
func (f *fn) StoreConst (v Value) uint32 {
	f.constantTable = append(f.constantTable,v)
	return uint32(len(f.constantTable)) - 1
}

//存储一条指令
func (f *fn) StoreOneInstruction (i uint32) int {
	f.instRuction = append(f.instRuction,i)
	return len(f.instRuction) - 1
}

func (f *fn) GetConst (index uint32) Value {
	return f.constantTable[index]
}

func (f *fn) writeBack (position uint32) {
	relative :=  len(f.instRuction) - int(position)
	f.instRuction[position] |= uint32(relative)
}

//stackv是指向栈中的value
//当作用域，vv将存储upvalue的值
type UpValue struct {
	stackV  *Value

	vv 		Value
	next    *UpValue
}

//闭包
type closure struct{
	f		*fn
	upValue []*UpValue
}

func NewClosure(f *fn) *closure {
	return &closure{
		f:       f,
		upValue: nil,
	}
}

type defType int8
const (
	LIBERDEF defType = iota
	GODEF
)

type def struct {
	c		  *closure
	t          defType
	gf         gofuncer
}


