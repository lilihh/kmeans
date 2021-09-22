package kmeans

type memberPoint struct {
	rawData           interface{}
	previousClusterID int
	currentClusterID  int
	isChanged         bool
}

type referencePoint struct {
	rawData     interface{}
	clusterID   int
	totalMember int
}

func (rp *referencePoint) reset(dp interface{}) {
	rp.rawData = dp
	rp.totalMember = 0
}
