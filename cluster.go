package kmeans

import (
	"math/rand"
	"time"
)

type clusterPool struct {
	distanceFunction Distance
	meanFunction     Mean
	clusterNum       int
	referencePoints  []*referencePoint
	memberPoints     []*memberPoint
}

func (pool *clusterPool) Result() [][]interface{} {
	for !pool.isStable() {
		// 根據現有分組，重新定義叢集的基準點
		pool.resetReferencePoints()

		for _, mp := range pool.memberPoints {
			// 重新計算每個點應該屬於哪個叢集
			pool.redirect(mp)
		}
	}

	// 結果
	return pool.currentClusters()
}

/*initialization related*/
func (pool *clusterPool) initRefPoints(data []interface{}) {
	// 用亂數選出基準點 (待改善)
	rand.Seed(time.Now().UnixNano())
	pool.referencePoints = make([]*referencePoint, 0, pool.clusterNum)

	for i := 0; i < pool.clusterNum; i++ {
		p := &referencePoint{
			rawData:   data[rand.Intn(len(data))],
			clusterID: i,
		}

		pool.referencePoints = append(pool.referencePoints, p)
	}
}

func (pool *clusterPool) initClusterPoints(data []interface{}) {
	pool.memberPoints = make([]*memberPoint, 0, len(data))

	for _, dp := range data {
		pool.memberPoints = append(pool.memberPoints, pool.transformToMemberPoint(dp))
	}
}

func (pool *clusterPool) transformToMemberPoint(dp interface{}) *memberPoint {
	clusterID := 0
	minDistance := pool.distanceFunction(dp, pool.referencePoints[0].rawData)

	for refPointID, refPoint := range pool.referencePoints {
		if minDistance > pool.distanceFunction(dp, refPoint.rawData) {
			clusterID = refPointID
			minDistance = pool.distanceFunction(dp, refPoint.rawData)
		}
	}

	pool.referencePoints[clusterID].totalMember++

	return &memberPoint{
		rawData:           dp,
		previousClusterID: -1,
		currentClusterID:  clusterID,
		isChanged:         true,
	}
}

/*grouping related*/
func (pool *clusterPool) resetReferencePoints() {
	clusters := pool.currentClusters()

	// 找到叢集的中間點
	for i, cluster := range clusters {
		newReferencePoint := pool.meanFunction(cluster...)
		pool.referencePoints[i].reset(newReferencePoint)
	}
}

func (pool *clusterPool) redirect(cp *memberPoint) *memberPoint {
	clusterID := 0
	minDistance := pool.distanceFunction(cp.rawData, pool.referencePoints[0].rawData)

	for refPointID, refPoint := range pool.referencePoints {
		if minDistance > pool.distanceFunction(cp.rawData, refPoint.rawData) {
			clusterID = refPointID
			minDistance = pool.distanceFunction(cp.rawData, refPoint.rawData)
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
func (pool *clusterPool) isStable() bool {
	hasChanged := false

	for _, cp := range pool.memberPoints {
		hasChanged = hasChanged || cp.isChanged
	}

	return !hasChanged
}

func (pool *clusterPool) currentClusters() [][]interface{} {
	// 初始化
	clusters := make([][]interface{}, pool.clusterNum)
	for _, refP := range pool.referencePoints {
		clusters[refP.clusterID] = make([]interface{}, 0, refP.totalMember)
	}

	// 將每個點分配到對應的叢集中
	for _, mp := range pool.memberPoints {
		clusterIDForMp := mp.currentClusterID
		clusters[clusterIDForMp] = append(clusters[clusterIDForMp], mp.rawData)
	}

	return clusters
}
