package saga

func fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	fibPrev := 0
	fibCurrent := 1

	for i := 2; i <= n; i++ {
		temp := fibCurrent
		fibCurrent += fibPrev
		fibPrev = temp
	}

	return fibCurrent
}
