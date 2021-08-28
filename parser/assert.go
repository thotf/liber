package parser

import "fmt"

func Assert (parser *Parser,t TokenT) {
	if parser.curToken.T != t {
		fmt.Println("此时为",parser.curToken)
		panic(fmt.Sprint("缺少",t))
	}
}