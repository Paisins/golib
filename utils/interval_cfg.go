package utils

import (
	"strconv"
	"strings"
)

type IntervalItem[T, V interface{ float64 | int32 }] struct {
	Interval [2]*T `yaml:"interval"`
	Value    V     `yaml:"value"`
}

// FormatIntervalItem 区间格式化
func FormatIntervalItem[T, V interface{ float64 | int32 }](intervals []IntervalItem[T, V]) string {
	var strBuilder strings.Builder
	for i, v := range intervals {
		if i != 0 {
			strBuilder.WriteString(",")
		}
		strBuilder.WriteString("[(")
		if v.Interval[0] == nil {
			strBuilder.WriteString("~")
		} else {
			strBuilder.WriteString(strconv.FormatFloat(float64(*v.Interval[0]), 'f', -1, 64))
		}
		strBuilder.WriteString(",")
		if v.Interval[0] == nil {
			strBuilder.WriteString("~")
		} else {
			strBuilder.WriteString(strconv.FormatFloat(float64(*v.Interval[0]), 'f', -1, 64))
		}
		strBuilder.WriteString("),")
		strBuilder.WriteString(strconv.FormatFloat(float64(v.Value), 'f', -1, 64))
		strBuilder.WriteString("]")
	}
	return strBuilder.String()
}

// GetIntervalValue 泛型区间值获取函数，左闭右开
func GetIntervalValue[T, V interface{ float64 | int32 }](intervals []IntervalItem[T, V], metric T) (V, bool) {
	for _, v := range intervals {
		if v.Interval[0] != nil && metric < *v.Interval[0] {
			continue
		}
		if v.Interval[1] != nil && metric >= *v.Interval[1] {
			continue
		}
		return v.Value, true
	}
	return 0, false
}
