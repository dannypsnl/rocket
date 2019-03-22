package filepath

func Join(paths ...string) string {
	pathLen := len(paths[0])
	for _, p := range paths[1:] {
		pathLen += len(p) + 1
	}
	path := make([]byte, pathLen)
	lastOne := len(paths) - 1
	index := 0
	for _, v := range paths[:lastOne] {
		copy(path[index:], []byte(v))
		index += len(v)
		path[index] = byte('/')
		index++
	}
	copy(path[index:], []byte(paths[lastOne]))
	return string(path)
}
