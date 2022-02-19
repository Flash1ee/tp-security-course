package proxy

func concatByteRespond(args ...[]byte) []byte {
	res := make([]byte, 0, bufferSize)
	for _, arg := range args {
		res = append(res, arg...)
	}
	return res
}
