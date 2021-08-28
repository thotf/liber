package parser

import "fmt"
//golang支持

type gofuncer interface {
	run(p ...Value)              //运行函数
	getReturn()   	[]Value     //获取函数返回值
	paramlen()     int8       //参数值长度
}

var goFuncs map[string] gofuncer

//golang函数注册
func init() {
	goFuncs = map[string]gofuncer{
		"println":&goPrintln{nil,1},
		"len"    :&goLen{make([]Value,0),1},
	}
}

type goPrintln struct {
	ret      []Value
	paramNum int8
}

func (g *goPrintln) run (p ...Value) {
	//输入为list处理
	if p[0].T == VT_lIST {
		l := p[0].V.(listObject)

		if l.islist {
			for _,v := range l.list.buf {
				fmt.Print(" ",v.V)
			}
			fmt.Println()
		} else {
			for i,v := range l.cmap {
				fmt.Print(" ",i.V,v.V)
			}
			fmt.Println()
		}
	}


	if len(p) > 1 {
		for _,v := range p {
			fmt.Print(" ",v.V)
		}
		fmt.Println()
	} else if len(p) == 1 {
		fmt.Println(p[0].V)
	}else{
		fmt.Println()
	}

}

func (g *goPrintln) getReturn () []Value {
	return nil
}
func (g *goPrintln) paramlen () int8 {
	return g.paramNum
}

type goLen struct {
	ret 	[]Value
	paramNum	int8
}

func (g *goLen) run (p ...Value) {
	lenght := 0
	if p[0].T == VT_lIST {
		l := p[0].V.(listObject)
		if l.islist {
			lenght = l.list.positon
		} else {
			lenght = len(l.cmap)
		}
	}

	v := new(Value)
	v.T = VT_INT
	v.V = lenght
	g.ret = append(g.ret,*v)
}
func (g *goLen) getReturn() []Value {
	//重置返回值数据表
	r := g.ret
	g.ret = make([]Value,0)
	return  r
}
func (g *goLen) paramlen () int8 {
	return g.paramNum
}

//函数是一直保存在内存中的，每次调用后记得清除数据
func funcLoading (p *Parser) {
	for funcName,fn := range goFuncs {
		p.ObjModule.AddGloablVarFunc(funcName,def{nil,GODEF,fn})
	}
}