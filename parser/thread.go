package parser

//读取当前fn中的指令，执行每条指令

type Thread struct {
	openValue	 *UpValue   //包含所有打开的Upvalue.

	Stack        []Value    //总栈

	curFrame     *Frame    //当前活动记录

	ip           uint

	stackTop	 int     //当前活动栈顶，永远指向最新活动的栈顶

	vm           *VM
}

//func (thread *Thread)
//执行指令
func (th *Thread) Run () {
	for  {
		if th.explainAnInstruction() == -1 {
			break
		}
	}
}

//解释一条指令
func (th *Thread) explainAnInstruction () int {
	i := th.getAnInstruction(th.ip)
	switch i.OpCode() {
	case IN_STORE:
		th.StoreConst(th.popStackTop())
		
	case IN_MOVE:
		th.StoreTop(th.Stack[i.R1()])

	case IN_OUT_MOVE:
		th.Stack[i.R1()] = th.popStackTop()

	case IN_LOAD:
		th.StoreTop(th.GetConst(i.R1()))

	case IN_OVER_LOAD:
		th.StoreTop(th.GetModuleVariable(i.R1()))

	case IN_OUT_OVER_LOAD:
		th.SetModuleVariable(i.R1(), th.popStackTop())

	case IN_PUSH:
		th.StoreTop(Value{})

	case IN_POP:
		th.popStackTop()

	case IN_ADD:
		fallthrough
	case IN_SUB:
		fallthrough
	case IN_MUL:
		fallthrough
	case IN_QUO:
		fallthrough
	case IN_REM:
		fallthrough
	case IN_SHL:
		fallthrough
	case IN_SHR:
		A2 := th.popStackTop()
		A1 := th.popStackTop()
		th.StoreTop(Calculation(A1,A2,i.OpCode()))

	case IN_EQL:
		fallthrough
	case IN_LSS:
		fallthrough
	case IN_GTR:
		fallthrough
	case IN_LAND:
		fallthrough
	case IN_LOR:
		A2 := th.popStackTop()
		A1 := th.popStackTop()
		th.StoreTop(BoolCalculation(A1,A2,i.OpCode()))
	case IN_NOT:
		th.StoreTop(ReverseBool(th.popStackTop()))

	case IN_NEGATE:
		th.StoreTop(TakeNegative(th.popStackTop()))

	case IN_JUMP:
		th.ip += uint(i.R1())
		return 0

	case IN_IF:
		top := th.popStackTop()
		if arbitrarilyToBool(top) {
			th.ip += 2
			return 0
		}

	case IN_RETURN:
		returnNum := int(i.R1())
		if returnNum > 0 {
			//有返回值,返回值已经在栈顶了
			returnNumberValue := make([]Value,returnNum)   //拿到返回值
			for num := 0 ; num < returnNum ; num ++ {
				returnNumberValue[num] = th.popStackTop()
			}
			th.curFrame = th.curFrame.preFrame
			th.resetThreadInFrame()  //返回parent栈
			//返回值Push进parent栈
			for num := returnNum - 1 ; num >=0 ; num -- {
				th.StoreTop(returnNumberValue[num])
			}
		} else {
			th.curFrame = th.curFrame.preFrame
			th.resetThreadInFrame()
		}

	case IN_CALL:
		callFun := th.GetFunc(i.R1()).(def)
		th.callFuns(callFun)
		return 0

	case IN_STACK_CALL:
		param := i.R1()
		a := make([]Value,0)
		for i:=uint32(0) ; i < param ; i ++ {
			a = append(a,th.popStackTop())
		}
		fn := th.popStackTop()

		for i:=int32(param) - 1; i >= 0 ; i -- {
			th.StoreTop(a[i])
		}

		th.callFuns(fn.V.(def))
		return 0

	case IN_WHILE:
		b := th.popStackTop()
		if	arbitrarilyToBool(b) {
			th.ip += 2
			return 0
		}

	case IN_JUMP_B:
		position := i.R1()
		th.ip = uint(position)
		return 0

	case IN_NEW_LIST:
		th.StoreTop(NewListValue())

	case IN_LIST_APPEND:
		index := i.R1()
		v := th.popStackTop()
		list := th.popStackTop().V.(listObject)
		list.add(Value{
			T: VT_INT,
			V: index,
		},v)
		th.StoreTop(Value{VT_lIST,list})

	case IN_LIST_ADD:
		v := th.popStackTop()
		key := th.popStackTop()
		list := th.popStackTop().V.(listObject)
		list.add(key,v)
		th.StoreTop(Value{VT_lIST,list})

	case IN_LIST_GET:
		index := th.popStackTop()
		list  := th.popStackTop()
		if list.T == VT_lIST {
			l := list.V.(listObject)
			th.StoreTop(l.Get(index))
		} else {
			panic("需要一个list")
		}

	case IN_LIST_OUT_OVER_LOAD:
		if th.GetModuleVariable(i.R1()).T == VT_lIST {
			v := th.popStackTop()
			index := th.popStackTop()
			list := th.GetModuleVariable(i.R1()).V.(listObject)
			list.Set(index, v)
			th.SetModuleVariable(i.R1(), Value{T: VT_lIST, V: list})
		}

	case IN_LIST_OUT_MOVE:
		if th.Stack[i.R1()].T == VT_lIST {
			v := th.popStackTop()
			index := th.popStackTop()
			list := th.Stack[i.R1()].V.(listObject)
			list.Set(index, v)
			th.Stack[i.R1()] = Value{T: VT_lIST, V: list}
		}
	case IN_EXIT:
		return -1
	}
	th.ip ++
	return 0
}

func (th *Thread) SetModuleVariable (index uint32,v Value)  {
	th.vm.curParser.ObjModule.SetVariableValue(index,&v)
}

func (th *Thread) GetModuleVariable (index uint32) Value {
	return  th.vm.curParser.ObjModule.GetModuleVariable(index)
}
//将数据存储到栈顶
func (th *Thread) StoreTop (v Value) {
	th.stackTop ++
	th.Stack[th.stackTop] = v
}

func (th *Thread) GetFunc (index uint32) interface{} {
	return th.vm.curParser.ObjModule.GetModuleVariable(index).V
}

func (th *Thread) StoreConst (v Value) {
	th.curFrame.curFn.f.StoreConst(v)
}

func (th *Thread) GetConst (index uint32) Value {
	return th.curFrame.curFn.f.GetConst(index)
}

func (th *Thread) getAnInstruction(ip uint) order {
	return order(th.curFrame.GetAnInstruction(ip))
}

//返回并弹出栈顶
func (th *Thread) popStackTop() Value {
	v := th.Stack[th.stackTop]
	th.Stack[th.stackTop] = Value{}
	th.stackTop --
	return v
}

func (th *Thread) contextSaving (fn closure) {
	th.curFrame.StackTop = th.stackTop
	th.curFrame.stack = th.Stack
	th.curFrame.ip = th.ip
	th.curFrame = InitNewFrame(&fn,0,-1,th.curFrame)
}

//根据活动框架重置线程
func (th *Thread) resetThreadInFrame() {
	th.ip = th.curFrame.ip
	th.stackTop = th.curFrame.StackTop
	if th.curFrame.stack == nil {        //调用新的函数时，需要初始化新函数的栈空间
		th.Stack = make([]Value,ThreadStackNumCount(th.curFrame.curFn.f))
	} else {
		th.Stack = th.curFrame.stack
	}
}

func InitNewThread (f *fn,vm *VM) *Thread {
	return &Thread{
		openValue: nil,
		Stack:     make([]Value,ThreadStackNumCount(f)),
		curFrame:  InitNewFrame(NewClosure(f),0,-1,nil),
		ip:        0,
		stackTop:  -1,
		vm:vm,
	}
}

//栈的大小为：所有变量数+2,2个空间是用来计算的
func ThreadStackNumCount (f *fn) uint {
	return f.maxStackNum + 4
}

func (th *Thread) callFuns (callFun def) {
	switch callFun.t{
	case LIBERDEF:
		parameterNumber := make([]Value, callFun.c.f.NumberParameters)

		//保存实际参数
		for i := 0; i < len(parameterNumber); i++ {
			parameterNumber[i] = th.popStackTop()
		}
		th.contextSaving(*callFun.c)
		th.resetThreadInFrame()
		//实际参数入栈
		for i := len(parameterNumber) - 1; i >= 0; i-- {
			th.StoreTop(parameterNumber[i])
		}

	case GODEF:
		//是否需要参数
		if callFun.gf.paramlen() > 0 {
			parameterNumber := make([]Value, callFun.gf.paramlen())
			for i := len(parameterNumber) - 1; i>=0 ; i-- {
				parameterNumber[i] = th.popStackTop()
			}
			callFun.gf.run(parameterNumber...)
		} else {
			callFun.gf.run()
		}

		//有返回值
		vs := callFun.gf.getReturn()
		if vs != nil {
			for _,v := range vs {
				th.StoreTop(v)
			}
		}

		th.ip ++
	}
}
