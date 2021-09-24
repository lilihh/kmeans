# kmeans

## Introduction
This repo implement [k-means algorithm](https://en.wikipedia.org/wiki/K-means_clustering) by using Go.

## Installation
```
$ go get github.com/lilihh/kmeans
```

## How to use
To use kmeans, you have to 
1. Define 2 functions on your own.

    * `Distance` function, which can calculate the distance between 2 data.

    * `Mean` function, which can give a mean result from a group of data.

2. Create a kmeans calculator by using `NewClusterCalculator(distanceFunc Distance, meanFunc Mean)`.
3. Call the method `Result(clusterNum int, data []interface{})` of the calculator with the parameters it asks, and you will derive the clusters.

## Example
```go
package main

import (
    "fmt"
    "math"

    "github.com/lilihh/kmeans"
)

func main() {
    // Step1: Define 2 functions

    type oneDimData struct {
    	value float64
    }

    var distanceFunc = func(d1, d2 interface{}) float64 {
    	return math.Abs(d1.(*oneDimData).value - d2.(*oneDimData).value)
    }

    var meanFunc = func(dps ...interface{}) interface{} {
    	meanValue := float64(0)
    	for _, dp := range dps {
    		meanValue += dp.(*oneDimData).value
    	}
    	meanValue = meanValue / float64(len(dps))

    	result := &oneDimData{
    		value: meanValue,
    	}

    	return result
    }

    // Step2: Create a kmeans calculator
    calculator := kmeans.NewClustercalculator(distanceFunc, meanFunc)

    // Step3: Input data and derive the clusters
    clusterNum := 2
    data := []interface{}{
    	&oneDimData{value: 10},
    	&oneDimData{value: 15},
    	&oneDimData{value: 51},
    	&oneDimData{value: 11},
    	&oneDimData{value: 50},
    	&oneDimData{value: 12},
    	&oneDimData{value: 17},
    	&oneDimData{value: 60},
    	&oneDimData{value: 8},
    	&oneDimData{value: 40},
    	&oneDimData{value: 13},
    	&oneDimData{value: 9},
    }

    clusters := calculator.Result(clusterNum, data)
    printCluster(clusters)
}

func printCluster(result [][]interface{}) {
	clusterNum := len(result)

	for clusterIndex := 0; clusterIndex < clusterNum; clusterIndex++ {
		fmt.Printf("cluster #%d\n", clusterIndex)
		clusterLen := len(result[clusterIndex])

		for memberIndex := 0; memberIndex < clusterLen; memberIndex++ {
			member := result[clusterIndex][memberIndex]
			fmt.Printf("%v ", member)
		}
		fmt.Println()
	}
}

```