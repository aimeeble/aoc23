package lib

func GCD(a, b int) int {
	if b == 0 {
		return a
	}
	return GCD(b, a%b)
}

func LCM(vs ...int) int {
	res := 1
	for _, v := range vs {
		res = res * v / GCD(res, v)
	}
	return res
}
