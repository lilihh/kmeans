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
			name: "one Dimension, 雙峰現象",
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
		{
			name: "one Dimension, 平均分佈",
			args: args{
				clusterNum: 2,
				data: []interface{}{
					&oneDimData{value: 10},
					&oneDimData{value: 1},
					&oneDimData{value: 2},
					&oneDimData{value: 7},
					&oneDimData{value: 8},
					&oneDimData{value: 4},
					&oneDimData{value: 3},
					&oneDimData{value: 9},
				},
			},
			wantResult: [][]interface{}{
				{
					&oneDimData{value: 1},
					&oneDimData{value: 2},
					&oneDimData{value: 4},
					&oneDimData{value: 3},
				},
				{
					&oneDimData{value: 10},
					&oneDimData{value: 7},
					&oneDimData{value: 8},
					&oneDimData{value: 9},
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
	// Remark:當叢集順序不同時，仍可視為一樣

	if len(c1) != len(c2) {
		return false
	}
	clusterNum := len(c1)

	// 把c2轉為Set
	c2Set := make(map[int][]interface{})
	for i := 0; i < clusterNum; i++ {
		c2Set[i] = c2[i]
	}

	// 針對c1的每個叢集，找遍c2Set看是否有一樣的
	for c1Index := 0; c1Index < clusterNum; c1Index++ {
		clusterOfC1 := c1[c1Index]
		isThereAMatch := false

		for c2Index, clusterOfC2 := range c2Set {
			// 一旦找到c2Set中有一樣的，就將他拿出c2Set
			if isSliceEqual(clusterOfC1, clusterOfC2) {
				delete(c2Set, c2Index)
				isThereAMatch = true
				break
			}
		}

		// 如果整個c2Set裡面都沒有，那就代表c1和c2一定不一樣
		if !isThereAMatch {
			return false
		}
	}

	return true
}

func isSliceEqual(s1, s2 []interface{}) bool {
	if len(s1) != len(s2) {
		return false
	}
	length := len(s1)

	for i := 0; i < length; i++ {
		m1 := s1[i]
		m2 := s2[i]
		if !reflect.DeepEqual(m1, m2) {
			return false
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
