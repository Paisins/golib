package algo

const (
	SoldierEps float32 = 25.0 // 5.0的平方，聚类中小兵间的最大距离
	SoldierPts int     = 2    // 聚类中小兵的最小数量
)

// Vec2 二维坐标
type Vec2 struct {

	// X id:1 x坐标
	X float32 `json:",omitempty" yaml:"x" tdr_field:"x" tdr_id:"1"`
	// Y id:2 y坐标
	Y float32 `json:",omitempty" yaml:"y" tdr_field:"y" tdr_id:"2"`
}

// 计算距离的平方和
func squaredDistance(a, b *Vec2) float32 {
	return (a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y)
}

// GetCentroidsFromDBSCAN 获取DBSCAN聚类算法得到的质心坐标
func GetCentroidsFromDBSCAN(data []*Vec2, eps float32, minPts int) []*Vec2 {
	labels := DBSCAN(data, eps, minPts)
	centroidsMap := CalculateCentroids(data, labels)
	centroids := make([]*Vec2, 0, len(centroidsMap))
	for _, v := range centroidsMap {
		centroids = append(centroids, v)
	}
	return centroids
}

// DBSCAN 聚类算法
func DBSCAN(data []*Vec2, eps float32, minPts int) []int {
	labels := make([]int, len(data))
	for i := range labels {
		labels[i] = -1 // 初始化为未分类
	}

	clusterID := 0
	for i := range data {
		if labels[i] != -1 {
			continue
		}

		neighbors := regionQuery(data, i, eps)
		if len(neighbors) < minPts {
			labels[i] = -2 // 标记为噪声
			continue
		}

		expandCluster(data, labels, i, neighbors, clusterID, eps, minPts)
		clusterID++
	}

	return labels
}

// CalculateCentroids 计算每个聚类的质心坐标
func CalculateCentroids(data []*Vec2, labels []int) map[int]*Vec2 {
	// 创建一个 map 来存储每个聚类的质心坐标
	centroids := make(map[int]*Vec2)
	// 创建一个 map 来存储每个聚类的点的数量
	clusterSizes := make(map[int]int)

	for i, label := range labels {
		if label < 0 {
			continue // 跳过噪声点
		}

		// 累加每个聚类的坐标
		if _, exists := centroids[label]; !exists {
			centroids[label] = &Vec2{X: 0, Y: 0}
		}
		centroids[label] = &Vec2{
			X: centroids[label].X + data[i].X,
			Y: centroids[label].Y + data[i].Y,
		}
		clusterSizes[label]++
	}

	// 计算每个聚类的质心坐标
	for label, centroid := range centroids {
		size := float32(clusterSizes[label])
		centroids[label] = &Vec2{
			X: centroid.X / size,
			Y: centroid.Y / size,
		}
	}

	return centroids
}

// 查找邻居
func regionQuery(data []*Vec2, idx int, eps float32) []int {
	var neighbors []int
	for i := range data {
		if squaredDistance(data[idx], data[i]) <= eps {
			neighbors = append(neighbors, i)
		}
	}
	return neighbors
}

// 扩展聚类
func expandCluster(data []*Vec2, labels []int, idx int, neighbors []int, clusterID int, eps float32, minPts int) {
	labels[idx] = clusterID
	for i := 0; i < len(neighbors); i++ {
		neighborIdx := neighbors[i]
		if labels[neighborIdx] == -2 {
			labels[neighborIdx] = clusterID
		}
		if labels[neighborIdx] != -1 {
			continue
		}

		labels[neighborIdx] = clusterID
		newNeighbors := regionQuery(data, neighborIdx, eps)
		if len(newNeighbors) >= minPts {
			neighbors = append(neighbors, newNeighbors...)
		}
	}
}
