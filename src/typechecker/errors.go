package typechecker

import "fmt"

var errorDivisionByZero error = fmt.Errorf("division by zero is forbidden")
var invalidOperationMsg string = "cannot evaluate numeric operation. unsupported operator %v"
