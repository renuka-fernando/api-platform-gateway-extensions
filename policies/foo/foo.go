package foo

import "log/slog"

func Sum(a, b int) int {
	sum := a + b
	slog.Info("sum from foo v1", "sum", sum)
	return sum
}
