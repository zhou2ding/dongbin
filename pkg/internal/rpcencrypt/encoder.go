package rpcencrypt

type encoder struct {
	codeStr string
	verify  func(byte) bool
}

func (e *encoder) getEncoder() string {
	return e.codeStr
}

func (e *encoder) contains(b byte) bool {
	return e.verify(b)
}

func (e *encoder) get(index int) byte {
	if index >= len(e.codeStr) {
		panic("out of range")
	}
	return e.codeStr[index]
}
