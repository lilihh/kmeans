package kmeans

import (
	"errors"
	"math/rand"
	"time"
)

type clusterCalculator struct {
	distanceFunction Distance
	meanFunction     Mean
	clusterNum       int
	referencePoints  []*referencePoint
	memberPoints     []*memberPoint
}

func (calculator *clusterCalculator) Result(clusterNum int, data []interface{}) ([][]interface{}, error) {
	// 初始化
	if clusterNum <= 0 {
		return nil, errors.New("the amount of clusters should be at least 1")
	} else {
		calculator.clusterNum = clusterNum
	}

	calculator.initRefPoints(data)
	calculator.initClusterPoints(data)

	// 計算
	for !calculator.isStable() {
		// 根據現有分組，重新定義叢集的基準點
		calculator.resetReferencePoints()

		for _, mp := range calculator.memberPoints {
			// 重新計算每個點應該屬於哪個叢集
			calculator.redirect(mp)
		}
	}

	// 結果
	return calculator.currentClusters(), nil
}

/*initialization related*/
func (calculator *clusterCalculator) initRefPoints(data []interface{}) {
	// TODO:待改善，當k值接近len(data)時，容易選出一樣的基準點
	// 用亂數選出基準點
	rand.Seed(time.Now().UnixNano())
	calculator.referencePoints = make([]*referencePoint, 0, calculator.clusterNum)

	for i := 0; i < calculator.clusterNum; i++ {
		p := &referencePoint{
			rawData:   data[rand.Intn(len(data))],
			clusterID: i,
		}

		calculator.referencePoints = append(calculator.referencePoints, p)
	}
}

func (calculator *clusterCalculator) initClusterPoints(data []interface{}) {
	calculator.memberPoints = make([]*memberPoint, 0, len(data))

	for _, dp := range data {
		calculator.memberPoints = append(calculator.memberPoints, calculator.transformToMemberPoint(dp))
	}
}

func (calculator *clusterCalculator) transformToMemberPoint(dp interface{}) *memberPoint {
	clusterID := 0
	minDistance := calculator.distanceFunction(dp, calculator.referencePoints[0].rawData)

	for refPointID, refPoint := range calculator.referencePoints {
		if minDistance > calculator.distanceFunction(dp, refPoint.rawData) {
			clusterID = refPointID
			minDistance = calculator.distanceFunction(dp, refPoint.rawData)
		}
	}

	calculator.referencePoints[clusterID].totalMember++

	return &memberPoint{
		rawData:           dp,
		previousClusterID: -1,
		currentClusterID:  clusterID,
		isChanged:         true,
	}
}

/*grouping related*/
func (calculator *clusterCalculator) resetReferencePoints() {
	clusters := calculator.currentClusters()

	// 找到叢集的中間點
	for i, cluster := range clusters {
		newReferencePoint := calculator.meanFunction(cluster...)
		calculator.referencePoints[i].reset(newReferencePoint)
	}
}

func (calculator *clusterCalculator) redirect(cp *memberPoint) *memberPoint {
	clusterID := 0
	minDistance := calculator.distanceFunction(cp.rawData, calculator.referencePoints[0].rawData)

	for refPointID, refPoint := range calculator.referencePoints {
		if minDistance > calculator.distanceFunction(cp.rawData, refPoint.rawData) {
			clusterID = refPointID
			minDistance = calculator.distanceFunction(cp.rawData, refPoint.rawData)
		}
	}

	cp.previousClusterID = cp.currentClusterID
	cp.currentClusterID = clusterID

	if cp.previousClusterID == cp.currentClusterID {
		cp.isChanged = false
	} else {
		cp.isChanged = true
	}

	return cp
}

/*current state related*/
func (calculator *clusterCalculator) isStable() bool {
	hasChanged := false

	for _, cp := range calculator.memberPoints {
		hasChanged = hasChanged || cp.isChanged
	}

	return !hasChanged
}

func (calculator *clusterCalculator) currentClusters() [][]interface{} {
	// 初始化
	clusters := make([][]interface{}, calculator.clusterNum)
	for _, refP := range calculator.referencePoints {
		clusters[refP.clusterID] = make([]interface{}, 0, refP.totalMember)
	}

	// 將每個點分配到對應的叢集中
	for _, mp := range calculator.memberPoints {
		clusterIDForMp := mp.currentClusterID
		clusters[clusterIDForMp] = append(clusters[clusterIDForMp], mp.rawData)
	}

	return clusters
}
