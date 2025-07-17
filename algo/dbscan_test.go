package algo

import (
	"fmt"
	"image/color"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// 绘制聚类结果和质心
func plotClustersWithCentroids(data []*Vec2, labels []int, centroids map[int]*Vec2, name string) {
	p := plot.New()

	p.Title.Text = "DBSCAN Clustering with Centroids"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// 设置坐标轴范围，增加边距
	p.X.Min = -50
	p.X.Max = 50
	p.Y.Min = -50
	p.Y.Max = 50

	// 使用 plotutil.Color 生成颜色
	for i, coord := range data {
		colorIndex := labels[i]
		if colorIndex < 0 {
			colorIndex = len(plotutil.DefaultColors) // 使用灰色表示噪声点
		}
		scatter, err := plotter.NewScatter(plotter.XYs{{X: float64(coord.X), Y: float64(coord.Y)}})
		if err != nil {
			panic(err)
		}
		scatter.GlyphStyle.Color = plotutil.Color(colorIndex)
		scatter.GlyphStyle.Radius = vg.Points(3)      // 调整点的大小
		scatter.GlyphStyle.Shape = draw.CircleGlyph{} // 使用实心圆形
		p.Add(scatter)
	}

	// 绘制质心
	for _, centroid := range centroids {
		scatter, err := plotter.NewScatter(plotter.XYs{{X: float64(centroid.X), Y: float64(centroid.Y)}})
		if err != nil {
			panic(err)
		}
		scatter.GlyphStyle.Color = color.Black       // 使用黑色
		scatter.GlyphStyle.Radius = vg.Points(5)     // 调整质心的大小
		scatter.GlyphStyle.Shape = draw.CrossGlyph{} // 使用十字形
		p.Add(scatter)
	}

	if err := p.Save(10*vg.Inch, 8*vg.Inch, fmt.Sprintf("clusters_with_centroids_%s.png", name)); err != nil {
		panic(err)
	}
}

func TestDBSCANOld(t *testing.T) {
	type args struct {
		data   []*Vec2
		eps    float32
		minPts int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "test_1",
			args: args{
				data: []*Vec2{
					{X: -45.478333, Y: -21.537476},
					{X: -44.873474, Y: -24.376282},
					{X: -46.16681, Y: -15.8029785},
					{X: -46.082886, Y: -18.698608},
					{X: -31.210938, Y: -29.40979},
					{X: -27.70221, Y: -26.615967},
					{X: -29.456665, Y: -28.013},
					{X: -25.947937, Y: -25.218933},
					{X: -21.464966, Y: -45.487427},
					{X: -18.626526, Y: -46.09546},
					{X: -15.731079, Y: -46.176575},
					{X: -24.302673, Y: -44.878967},
					{X: 24.497131, Y: 44.75885},
					{X: 15.900757, Y: 46.011353},
					{X: 18.792603, Y: 45.840576},
					{X: 21.644958, Y: 45.29962},
					{X: 25.042603, Y: 25.737854},
					{X: 26.437439, Y: 27.495605},
					{X: 27.832275, Y: 29.253235},
					{X: 29.22705, Y: 31.010742},
					{X: 45.998047, Y: 15.992798},
					{X: 45.824707, Y: 18.884827},
					{X: 45.285767, Y: 21.737488},
					{X: 44.747192, Y: 24.590149},
				},
				eps:    SoldierEps,
				minPts: SoldierPts,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := DBSCAN(tt.args.data, tt.args.eps, tt.args.minPts)
			// 输出聚类结果
			// fmt.Println("\n 聚类结果： ")
			// for i, label := range labels {
			// 	fmt.Printf("点 %d: (%f, %f) -> 簇 %d\n", i, tt.args.data[i].X, tt.args.data[i].Y, label)
			// }

			// 输出质心坐标
			centroids := CalculateCentroids(tt.args.data, labels)

			// 绘制聚类结果
			plotClustersWithCentroids(tt.args.data, labels, centroids, tt.name)
		})
	}
}
