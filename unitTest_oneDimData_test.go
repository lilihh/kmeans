package kmeans

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

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

func TestOneDimensionalData(t *testing.T) {
	calculator := NewClustercalculator(distanceFunc, meanFunc)

	type args struct {
		clusterNum int
		data       []interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantResult [][]interface{}
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "1D, k=2",
			args: args{
				clusterNum: 2,
				data: []interface{}{
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
				},
			},
			wantResult: [][]interface{}{
				{
					&oneDimData{value: 10},
					&oneDimData{value: 15},
					&oneDimData{value: 11},
					&oneDimData{value: 12},
					&oneDimData{value: 17},
					&oneDimData{value: 8},
					&oneDimData{value: 13},
					&oneDimData{value: 9},
				},
				{
					&oneDimData{value: 51},
					&oneDimData{value: 50},
					&oneDimData{value: 60},
					&oneDimData{value: 40},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.Result(tt.args.clusterNum, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Result() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !isDeepEqual(result, tt.wantResult) {
				t.Errorf("result and wantResult is diff \n")
				fmt.Println("Calculator's result:")
				printCluster(result)
				fmt.Println("wantResult:")
				printCluster(tt.wantResult)
				return
			}
		})
	}
}

func isDeepEqual(c1, c2 [][]interface{}) bool {
	// TODO:待改善，叢集順序不同時，仍可視為一樣

	if len(c1) != len(c2) {
		return false
	}
	clusterNum := len(c1)

	for clusterIndex := 0; clusterIndex < clusterNum; clusterIndex++ {
		if len(c1[clusterIndex]) != len(c2[clusterIndex]) {
			return false
		}
		clusterLen := len(c1[clusterIndex])

		for memberIndex := 0; memberIndex < clusterLen; memberIndex++ {
			m1 := c1[clusterIndex][memberIndex]
			m2 := c1[clusterIndex][memberIndex]
			if !reflect.DeepEqual(m1, m2) {
				return false
			}
		}
	}

	return true
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
