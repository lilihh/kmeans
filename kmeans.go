package kmeans

/*
To use kmeans, you have to define 2 functions on your own.
The first one is Distance function, which can calculate the distance between 2 data.
The second one is Mean function, which can give a mean result from a group of data.
After define these 2 functions, creating a kmeans calculator by using NewClusterCalculator().

Then call the method Result() of the calculator with the parameters it asks, and you will derive the clusters.
*/

type Distance func(d1, d2 interface{}) float64

type Mean func(...interface{}) interface{}

func NewClustercalculator(distanceFunc Distance, meanFunc Mean) *clusterCalculator {
	calculator := &clusterCalculator{
		distanceFunction: distanceFunc,
		meanFunction:     meanFunc,
	}

	return calculator
}
