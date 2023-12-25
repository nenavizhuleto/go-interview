package company

type feature int

const (
	FeatureBranchEmployees feature = iota
	FeatureBranchNetwork

	FeatureCount // Total count of plugin features
)

var features []bool = make([]bool, FeatureCount)

func EnableFeature(f feature) {
	if f == FeatureCount {
		return
	}

	features[f] = true
}

func FeatureEnabled(f feature) bool {
	if f == FeatureCount {
		return false
	}
	return features[f]
}

func DisableFeature(f feature) {
	if f == FeatureCount {
		return
	}

	features[f] = false
}
