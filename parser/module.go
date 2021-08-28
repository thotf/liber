package parser

import "fmt"

//每个文件都是一个模块
//模块中有模块级变量也就是全局变量，包括函数
const(
	MAX_MODULE_VARIABLE = 128
)

type Module struct{
	//模块变量值
	V      			[]*Value

	//模块变量名
	name			[]string

	//变量数
	number 			uint8

	p               *Parser
}

func (m *Module) InsertModuleVariable (name string) uint8 {
	if m.number == (MAX_MODULE_VARIABLE - 1) {
		panic("变量添加失败，变量数量超过限制")
	}
		m.name[m.number] = name
		m.V[m.number] = &Value{}
		m.number ++

	return m.number - 1

}

func (m *Module) GetModuleVariable (index uint32) Value {
	return *m.V[index]
}

 //查看变量是否存在
func (m *Module) IsExist (name string) (bool,uint8) {
		for index,vName := range m.name {
			if vName == name {
				return true,uint8(index)
			}
		}
	return false,0
}

//返回index为 变量位置,不存在返回-1
func (m *Module) FindVariable (name string) (uint8,bool) {
	if ok,index := m.IsExist(name) ; ok {
		return index,true
	}
	return 0,false
	//return m.InsertModuleVariable(name)
}

func InitModule(name string) *Module {
	return &Module{
		V:    make([]*Value,MAX_MODULE_VARIABLE),
		name: make([]string,MAX_MODULE_VARIABLE),
	}
}

func (m *Module) SetVariableValue (index uint32,v *Value) {
	m.V[index] = v
}

//将函数作为全局变量添加到模块中
func (m *Module) AddGloablVarFunc(name string,c def) {
	if _,ok:= m.FindVariable(name) ; ok != false {
		panic(fmt.Sprint("重复定义",name))
	}
	v := new(Value)
	v.T = VT_FUN
	v.V = c
	m.SetVariableValue(uint32(m.InsertModuleVariable(name)),v)
}

func (m *Module) findAndGetFunc (name string) (Value,uint8,bool) {
	if ok,index := m.IsExist(name) ; ok {
		if V := m.GetModuleVariable(uint32(index)); V.T == VT_FUN {
			return V,index,ok
		}
		return Value{},index,false
	}
	return Value{},0,false
}