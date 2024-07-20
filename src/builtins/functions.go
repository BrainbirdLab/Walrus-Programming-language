package builtins

import (
	"fmt"
	"time"

	"walrus/typechecker"
)

func NativePrint(args ...typechecker.RuntimeValue) typechecker.RuntimeValue {

	//if no arguments
	if len(args) == 0 {
		fmt.Println()
		return typechecker.MakeVOID()
	}

	for _, arg := range args {
		val, err := typechecker.CastToStringValue(arg)

		if err != nil {
			continue
		}

		//colorize
		fmt.Print(val.Value)
	}
	fmt.Println()
	return typechecker.MakeVOID()
}

func NativeTime(args ...typechecker.RuntimeValue) typechecker.RuntimeValue {
	t := time.Now().Unix()
	return typechecker.MakeINT(t, 64, true)
}

func NativeLen(args ...typechecker.RuntimeValue) typechecker.RuntimeValue {
	switch a := args[0].(type) {
	case typechecker.ArrayValue:
		size := len(a.Values)		
		return typechecker.MakeINT(int64(size), 64, true)
	default:
		panic("invalid argument for len()")
	}
}