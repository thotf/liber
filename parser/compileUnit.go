package parser

//每个作用域都有一个编译单元，
//编译单元在编译期间保存数据，指向数据，记录各种信息
//通过编译单元的各种信息得到生成指令

type CompileUnit struct {
	//所编译函数
	f  				*fn

	//本层函数锁使用的upvalue个数
	UpvalueSum		int

	//本层处于的作用域
	ScopeLevel    	int

	//当前循环层
	Loop			*Loop

	//直接外层
	upCompileUnit	*CompileUnit

	//当前parser
	p 				*Parser

	//局部变量链
	LocalVar        *LocalVar
}

func (c *CompileUnit) ReturnFn() *fn {
	return c.f
}
//这里常量指的是字面量
func (c *CompileUnit) storeConst (v Value) uint32 {
	return c.f.StoreConst(v)
}

//存储一条指令,返回当前写入指令的位置
func (c *CompileUnit) StoreOneInstruction (i uint32) int {
	return c.f.StoreOneInstruction(i)
}

//添加一个新的局部变量
func (c *CompileUnit)AddLocalVarS (name string) *LocalVar {
	var newLocal  LocalVar
	newLocal.LocalT     = LOCAL_VARIABLE
	newLocal.ScopeLevel = c.ScopeLevel
	newLocal.Name		= name
	newLocal.IsQuote	= false
	c.f.maxStackNum ++
	return &newLocal
}

//添加新的局部变量到局部变量链表,返回局部变量位置
func (c *CompileUnit) addLocalToLinkList (local *LocalVar) uint8 {

	localHead := c.LocalVar
	index     := uint8(0)

	if localHead == nil {
		c.LocalVar = local
		return 0
	}
	index ++
	for localHead.NextLocal != nil {
		localHead = localHead.NextLocal
		index ++
	}
	localHead.NextLocal = local
	return index
}

//添加一个局部变量
func (c *CompileUnit) AddLocalVar (name string) uint8 {
	c.StoreOneInstruction(push())
	return c.addLocalToLinkList(c.AddLocalVarS(name))
}

//添加一个参数，不用push入栈，在后端会有特别处理实际参数入栈不用指令
func (c *CompileUnit) AddParameter(name string) uint8 {
	return c.addLocalToLinkList(c.AddLocalVarS(name))
}

//查找局部变量并返回对应下标，没找到返回false
func (c *CompileUnit) findLocalVar(name string) (uint8,bool) {
	index := uint8(0)
	aLocal := c.LocalVar

	for aLocal != nil {
		if aLocal.Name == name {
			return index,true
		}
		index ++
		aLocal = aLocal.NextLocal
	}

	return 0,false
}

//查找局部变量，不存在则创建新的局部变量并返回index
func (c *CompileUnit) FindLocalAndCreateNewLocal(name string) uint8 {
	if index,ok := c.findLocalVar(name) ; ok == false  {
		return c.AddLocalVar(name)
	}else {
		return uint8(index)
	}
}

func (c *CompileUnit) enterScope () {
	c.ScopeLevel ++
}

//弹出局部作用域变量
func (c *CompileUnit) popUpLocalScopeVariables (local *LocalVar) *LocalVar {
	if local == nil {
		return nil
	}

	local.NextLocal = c.popUpLocalScopeVariables(local.NextLocal)
	if local.ScopeLevel >= c.ScopeLevel {
		c.StoreOneInstruction(pop())
		return nil
	} else {
		return local
	}
}

func (c *CompileUnit) exitScope () {
	//弹出所有当前作用域的变量
	c.LocalVar = c.popUpLocalScopeVariables(c.LocalVar)

	c.ScopeLevel --
}

type Loop struct {}

func InitCompileUnit(scopeLevel int,upCompileUnit *CompileUnit,p *Parser) *CompileUnit {

	return &CompileUnit{
		f:             &fn{},
		UpvalueSum:    0,
		ScopeLevel:    scopeLevel,
		Loop:          nil,
		upCompileUnit: upCompileUnit,
		p:             p,
		LocalVar:nil,
	}
}
