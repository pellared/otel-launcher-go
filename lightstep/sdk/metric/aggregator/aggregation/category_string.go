// Code generated by "stringer -type=Category,Kind"; DO NOT EDIT.

package aggregation

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UndefinedCategory-0]
	_ = x[MonotonicSumCategory-1]
	_ = x[NonMonotonicSumCategory-2]
	_ = x[GaugeCategory-3]
	_ = x[HistogramCategory-4]
}

const _Category_name = "UndefinedCategoryMonotonicSumCategoryNonMonotonicSumCategoryGaugeCategoryHistogramCategory"

var _Category_index = [...]uint8{0, 17, 37, 60, 73, 90}

func (i Category) String() string {
	if i < 0 || i >= Category(len(_Category_index)-1) {
		return "Category(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Category_name[_Category_index[i]:_Category_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UndefinedKind-0]
	_ = x[DropKind-1]
	_ = x[AnySumKind-2]
	_ = x[MonotonicSumKind-3]
	_ = x[NonMonotonicSumKind-4]
	_ = x[LatestSumKind-5]
	_ = x[GaugeKind-6]
	_ = x[HistogramKind-7]
	_ = x[MinMaxSumCountKind-8]
}

const _Kind_name = "UndefinedKindDropKindAnySumKindMonotonicSumKindNonMonotonicSumKindLatestSumKindGaugeKindHistogramKindMinMaxSumCountKind"

var _Kind_index = [...]uint8{0, 13, 21, 31, 47, 66, 79, 88, 101, 119}

func (i Kind) String() string {
	if i < 0 || i >= Kind(len(_Kind_index)-1) {
		return "Kind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
