package metronome

import "errors"

var ChunkSize int = 1024
var NotEnoughSamples error = errors.New("not enough samples")

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Energy(data [][2]float64) (float64, error) {
	if len(data) < ChunkSize {
		return 0, NotEnoughSamples
	}
	e := 0.0
	for i := 0; i < ChunkSize; i++ {
		l := data[i][0]
		r := data[i][1]
		e += l*l + r*r
	}
	return e, nil
}

func mean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func variance(data []float64, mean float64) float64 {
	sum := 0.0
	for _, v := range data {
		u := (v - mean)
		sum += u * u
	}
	return sum / float64(len(data))
}

func getC(variance float64) float64 {
	return -0.0000015*variance + 1.5142857
}
