package parser

import (
	"bufio"
	"os"
)

type VM struct {
	curThread	*Thread
	curParser	*Parser
}

func InitVm(s string) *VM {
	f, err := os.Open(s)
	if err != nil {
		panic(err)
	}
	aFile := bufio.NewReader(f)

	vm := new(VM)
	vm.curParser =  InitNewParser(aFile, f.Name())

	return vm
}

func (v *VM) Run () {
	v.curParser.ReadyToCompile()

	v.curThread = InitNewThread(v.curParser.CompileUnit.ReturnFn(),v)
	v.curThread.Run()
}