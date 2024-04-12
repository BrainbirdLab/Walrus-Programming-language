package utils

func IF(conddition bool, a, b any) any {
	if conddition {
		return a
	}
	return b
}