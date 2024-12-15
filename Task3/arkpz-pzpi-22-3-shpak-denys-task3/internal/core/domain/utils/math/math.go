package math

import "math"

func Transpose(matrix [][]float64) [][]float64 {
	rows := len(matrix)
	cols := len(matrix[0])
	result := make([][]float64, cols)
	for i := range result {
		result[i] = make([]float64, rows)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			result[j][i] = matrix[i][j]
		}
	}
	return result
}

func MultiplyMatrices(a, b [][]float64) [][]float64 {
	rowsA := len(a)
	colsA := len(a[0])
	colsB := len(b[0])
	result := make([][]float64, rowsA)
	for i := range result {
		result[i] = make([]float64, colsB)
	}
	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsB; j++ {
			for k := 0; k < colsA; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

func Inverse(matrix [][]float64) [][]float64 {

	n := len(matrix)
	augmented := make([][]float64, n)
	for i := range augmented {
		augmented[i] = make([]float64, 2*n)
		copy(augmented[i], matrix[i])
		augmented[i][n+i] = 1
	}

	for i := 0; i < n; i++ {

		maxRow := i
		for j := i + 1; j < n; j++ {
			if math.Abs(augmented[j][i]) > math.Abs(augmented[maxRow][i]) {
				maxRow = j
			}
		}

		augmented[i], augmented[maxRow] = augmented[maxRow], augmented[i]

		pivot := augmented[i][i]
		for j := 0; j < 2*n; j++ {
			augmented[i][j] /= pivot
		}

		for j := i + 1; j < n; j++ {
			factor := augmented[j][i]
			for k := 0; k < 2*n; k++ {
				augmented[j][k] -= factor * augmented[i][k]
			}
		}
	}

	for i := n - 1; i >= 0; i-- {
		for j := i - 1; j >= 0; j-- {
			factor := augmented[j][i]
			for k := 0; k < 2*n; k++ {
				augmented[j][k] -= factor * augmented[i][k]
			}
		}
	}

	inverse := make([][]float64, n)
	for i := range inverse {
		inverse[i] = make([]float64, n)
		copy(inverse[i], augmented[i][n:])
	}

	return inverse
}

func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func Mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

func StdDev(data []float64) float64 {
	mean := Mean(data)
	sumSquares := 0.0
	for _, value := range data {
		diff := value - mean
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(data))
	return math.Sqrt(variance)
}

func Square(data []float64) []float64 {
	result := make([]float64, len(data))
	for i, value := range data {
		result[i] = value * value
	}
	return result
}

func Multiply(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range result {
		result[i] = a[i] * b[i]
	}
	return result
}

func Sum(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum
}
