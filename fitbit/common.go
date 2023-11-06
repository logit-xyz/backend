package fitbit

import (
    "strconv"
)

func ConvertFloat(f float64, n int) string {
    return strconv.FormatFloat(f, 'f', n, 64)
}
