package lib

func Hash(rule string) int {
	var cur int
	for _, c := range rule {
		cur += int(c)
		cur = (cur * 17) % 256
	}
	return cur
}
