package main

import (
	"nanozkv/vm"
	"nanozkv/zk"
)

func main() {

	code := []vm.Instruction{
		{Op: vm.PUSH, Arg: 27},
		{Op: vm.PUSH, Arg: 3},
		{Op: vm.ADD},
		{Op: vm.PUSH, Arg: 81},
		{Op: vm.MUL},
		{Op: vm.HALT},
	}

	err := zk.RunTraceProof(code)

	if err != nil {
		panic(err)
	}
}
