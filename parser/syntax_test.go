package parser

import (
	"fmt"
	"testing"
)

func TestExpression(t *testing.T) {
	vm := InitVm()
	vm.Run("abc")
	fmt.Println("指令",vm.curParser.CompileUnit.f.instRuction)
	fmt.Println("符号表",vm.curParser.CompileUnit.f.constantTable)
	fmt.Println("局部变量表",vm.curParser.CompileUnit.LocalVar)
	//fmt.Println("函数指令表",vm.curParser.ObjModule.V[0].V.(closure).f.instRuction)
	//fmt.Println("参数数量",vm.curParser.ObjModule.V[0].V.(closure).f.NumberParameters)
	//fmt.Println("栈大小",vm.curParser.ObjModule.V[0].V.(closure).f.maxStackNum)
	//fmt.Println("常量表",vm.curParser.ObjModule.V[0].V.(closure).f.constantTable)
	//fmt.Println("全局变量表",vm.curParser.ObjModule.name,vm.curParser.ObjModule.number)
	fmt.Println("栈",vm.curThread.Stack,vm.curThread.stackTop)
	fmt.Println(vm.curParser.ObjModule.V[1])
	//fmt.Println(vm.curParser.ObjModule.V[0].V.(listObject).list.buf)
	//for i:=1;i<5;i++ {
	//	fmt.Print("  ",vm.curParser.ObjModule.V[i].V.(int))
	//}
}
/*
将栈顶元素赋值给VarType，index 1 0
将栈顶元素赋值给VarType，index 1 1
2
指令 [33554432 33554433 83886080 33554434 83886080 33554435 83886080 33554436 83886080 33554437 335544
38 83886080]
符号表 [{0 1} {0 2} {0 3} {0 4} {0 5} {0 3} {0 1}]
局部变量表 <nil>
全局变量表 [a b
                                       ] 2
*/
