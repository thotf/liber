package parser

import "fmt"

type LocalVarType int8

//常用语句解析，
//控制语句 if for do while
//数据处理语句 函数调用

const (
	NILL_VARIABLE   LocalVarType = iota
	GLOBAL_VARIABLE
	LOCAL_VARIABLE
	UPVALUE_VARIABLE
)

type LocalVar struct {
	LocalT 		LocalVarType
	Name        string       //局部变量名
	IsQuote		bool         //是否被低作用域所引用
	ScopeLevel  int          //变量所属作用域
	valueIndex  uint
	NextLocal  *LocalVar
}
/*
if expr {
1expr
} else {
2expr
}
写回：发现为if语句，expr的bool值入栈，写入一条空指令A，返回Aposition
解析1expr后，将1expr后面语句的指令的位置写入到Aposition。注意要先写入pop指令退出当前作用域函数后作用域后再写Aposition
写回执行：当expr为真则ip += 1 ，跳过Aposition的指令，执行1expr

expr为真，执行1expr，在1expr加一个JUMP命令跳转到2expr后面。
expr为假，跳转到2expr执行。

*/

func ifStatement (p *Parser) {
	p.GetToken()                //此时curToken为if后的表达式开头
	Expression(p,0)
	Assert(p,LBRACE)
	p.StoreOneInstruction(tableVarDisplace(IN_IF,0))   //写入这是if语句，接下来写入跳转到}后的位置
	writeBackPositon := uint32(p.CompileUnit.StoreOneInstruction(tableVarDisplace(IN_JUMP,0)))   //当前指令位置，用于写回技术
	analysisLRbrack(p,false)

	signOutElseBackPosition :=  uint32(p.CompileUnit.StoreOneInstruction(tableVarDisplace(IN_JUMP,0)))
	p.CompileUnit.f.writeBack(writeBackPositon)

	//解析else
	if p.curToken.T == ELSE  && p.preToken.line == p.curToken.line {
		p.GetToken()
		analysisLRbrack(p,false)
		p.CompileUnit.f.writeBack(signOutElseBackPosition)
		return
	}
	p.CompileUnit.f.writeBack(signOutElseBackPosition)
}

//解析{}内的表达式
func analysisLRbrack (p *Parser,isFunc bool) {
	Assert(p,LBRACE)
	p.CompileUnit.enterScope()   //进入作用域
	p.GetToken()             //读入“{”的下一个词素

	 //解析{}内的表达式
	 for p.curToken.T != RBRACE {
	 	p.analysisOneSentence()

	 	if p.curToken.T == EOF {
	 		panic(fmt.Sprint("缺少}",p.curToken.line))
		}
	 }

	Assert(p,RBRACE)
	p.GetToken()  //读入“}”的下一个词素
	if !isFunc {   //函数执行完直接销毁不用退出作用域
		p.CompileUnit.exitScope() //退出作用域
	}
}

//def abc(a,b,c) {}
//函数解析
//函数每个return必须具有相同的数量
func defStatement (p *Parser) {
	funcCompileUnit := InitCompileUnit(1,p.CompileUnit,p)
	p.CompileUnit = funcCompileUnit

	funcName := p.GetToken().Value.(string)   //拿到函数名
	p.GetToken()

	//解析函数参数，入栈
	Assert(p,LPAREN)
	for p.GetToken().T == ID {
		p.CompileUnit.AddParameter(p.curToken.Value.(string))
		p.CompileUnit.f.NumberParameters ++
		if p.GetToken().T != COMMA {
			break
		}
	}

	//
	Assert(p,RPAREN)

    //解析函数{} 中的内容
	p.GetToken()
	analysisLRbrack(p,true)              //函数中的指令加载

	p.StoreOneInstruction(insDisplace(IN_RETURN,0))  //这里又添加一个return，因为程序编写时可以不添加return
	newFunc := NewClosure(funcCompileUnit.f)
	newFunc.f.module = p.ObjModule
	p.CompileUnit = p.CompileUnit.upCompileUnit
	p.ObjModule.AddGloablVarFunc(funcName,def{newFunc,LIBERDEF,nil})   //将函数添加到全局变量中
}

func returnParser (p *Parser) {
	//p.StoreOneInstruction(varTypeLoadInstruct(VarType,uint8(index)))
	p.GetToken()

	//return后的返回值入栈
	p.CompileUnit.f.ReturnValue = 0    //重置return返回值数量，每个return返回数可能不同
	if p.curToken.T != RBRACE {
		Expression(p,0)
		p.CompileUnit.f.ReturnValue ++
		p.CompileUnit.f.maxStackNum ++
		for p.curToken.T == COMMA {
			p.CompileUnit.f.ReturnValue ++
			p.CompileUnit.f.maxStackNum ++
			p.GetToken()
			Expression(p,0)
		}
	}

	p.StoreOneInstruction(insDisplace(IN_RETURN,uint8(p.CompileUnit.f.ReturnValue)))
}


//返回返回值个数
func funcCall(p *Parser) int8 {
	fnname := p.curToken.Value.(string)
	fnValue,index,ok := p.ObjModule.findAndGetFunc(fnname)
	p.GetToken()
	p.GetToken()

	numberOfParameters := 0   //参数数量计算
	if ok {
		if p.curToken.T != RPAREN {
			numberOfParameters ++
			Expression(p, 0)
			p.CompileUnit.f.maxStackNum++
			for p.curToken.T == COMMA {
				p.CompileUnit.f.maxStackNum++
				numberOfParameters ++
				p.GetToken()
				Expression(p, 0)
			}
		}
		p.StoreOneInstruction(insDisplace(IN_CALL, index))
		p.CompileUnit.f.maxStackNum += uint(numberOfParameters)

	}else {
		fmt.Println("没找到名为",fnname,"的函数")
		panic("")
	}
	if fnValue.V.(def).t == LIBERDEF {
		return fnValue.V.(def).c.f.ReturnValue
	}

	return int8(len(fnValue.V.(def).gf.getReturn()))
}

//获取拆包变量
func getUnpackingVariables (p *Parser) []commaVar {
	commaVars := make([]commaVar,0)

	for   {
		var v commaVar
		VarType,index := returnVarindex(p)
		v.varType = VarType
		v.index   = index
		commaVars = append(commaVars,v)

		if p.curToken.T == ASSIGN {
			break
		}
		p.GetToken()
		if p.curToken.T == ID {
			p.GetToken()
		}
	}
	return commaVars
}

func unpacKingHandle(p *Parser) {
	//拆包
	//获取每个变量 varType，index 放在[]中
	commaVars := getUnpackingVariables(p)
	p.GetToken()
	// 加载函数，判断len（index）和函数返回值数量是否相同
	if int(funcCall(p)) == len(commaVars) {
		//函数所有返回值放在栈顶，
		//从[len-1] 开始，从栈顶拿到值进行赋值
		for i:=len(commaVars) - 1; i >= 0 ; i -- {
			assignmentInstruction(p,commaVars[i].varType,commaVars[i].index)
		}
		p.GetToken()
	fmt.Println(p.curToken)
	} else {
		fmt.Println("赋值变量数与返回值不匹配")
		panic("")
	}
}

//赋值函数处理
func assignmentStatementHandle (p *Parser) {
	//赋值语句处理
	//以局部变量、upValue、全局变量的顺序查找 “=”左侧的变量，不存在则创建局部变量
	VarType,index := returnVarindex(p)
	//解析表达式,getToken后为“=”右侧第一个值
	p.GetToken()
	Expression(p,0)
	//将栈顶元素赋值给标识符所代表的变量
	assignmentInstruction(p,VarType,index)
}

//A：while bool jump B {
//bool为真

//jump_B A  跳转到开头
//}
//B: bool为假退出
func whileStatement (p *Parser) {
	p.GetToken()
	writeBackPositonA := uint32(len(p.CompileUnit.f.instRuction)) //while内语句执行结束跳转到开头
	Expression(p,0)
	p.StoreOneInstruction(tableVarDisplace(IN_WHILE,0))
	writeBackPositonB := uint32(p.CompileUnit.StoreOneInstruction(tableVarDisplace(IN_JUMP,0))) //跳出while循环

	analysisLRbrack(p,false)
	p.CompileUnit.StoreOneInstruction(tableVarDisplace(IN_JUMP_B,writeBackPositonA))
	p.CompileUnit.f.writeBack(writeBackPositonB)
}

//切片赋值
func sectionAssignmentStatementHandle(p *Parser) {
	//VarType,index := returnVarindex(p)
	//p.GetToken()
	//Expression(p,0)    //切片索引在栈顶
	//Assert(p,RBRACK)
	//p.GetToken()
	//Assert(p,ASSIGN)       //a[123] =
	//
	//Expression(p,0)
	//assignmentInstruction(p,VarType,index)
	VarType,index,ok := findVarInLocalAndUpValueAndGlobal(p)
	if !ok {
		fmt.Println("不存在变量",returnIdName(p))
		panic("")
	}
	p.GetToken()
	Expression(p,0)
	Assert(p,RBRACK)
	p.GetToken()
	Assert(p,ASSIGN)
	p.GetToken()
	Expression(p,0)
	listAssignment(p,VarType,index)
}

func listAssignment (p *Parser,varType LocalVarType,index uint8) {
	switch varType {
	case GLOBAL_VARIABLE:
		p.StoreOneInstruction(tableVarDisplace(IN_LIST_OUT_OVER_LOAD, uint32(index)))
	case LOCAL_VARIABLE:
		p.StoreOneInstruction(tableVarDisplace(IN_LIST_OUT_MOVE, uint32(index)))
	}
}