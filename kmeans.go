package kmeans

type Distance func(d1, d2 interface{}) float64

type Mean func(...interface{}) interface{}

func NewClustercalculator(distanceFunc Distance, meanFunc Mean) *clusterCalculator {
	calculator := &clusterCalculator{
		distanceFunction: distanceFunc,
		meanFunction:     meanFunc,
	}

	return calculator
}
