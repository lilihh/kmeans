package kmeans

import "errors"

/*
To use kmeans, you have to define 2 functions on your own.
The first one is Distance function, which can calculate the distance between 2 data.
The second one is Mean function, which can give a mean result from a group of data.

After define these 2 functions, all you have to do is creating a kmeans pool by using NewClusterPool() with the parameters it asks.
Then call the method Result() of the pool, and you will derive the clusters.
*/

type Distance func(d1, d2 interface{}) float64

type Mean func(...interface{}) interface{}

func NewClusterPool(clusterNum int, data []interface{}, distanceFunc Distance, meanFunc Mean) (*clusterPool, error) {
	if clusterNum <= 0 {
		return nil, errors.New("the amount of clusters should be at least 1")
	}

	pool := &clusterPool{
		distanceFunction: distanceFunc,
		meanFunction:     meanFunc,
		clusterNum:       clusterNum,
	}

	// 初始化
	pool.initRefPoints(data)
	pool.initClusterPoints(data)

	return pool, nil
}
