package model

// FeatureMap holds information about Oracle database feature
type FeatureMap map[string]interface{}

// Name getter
func (feature *FeatureMap) Name() string {
	name := (*feature)["Name"]
	return name.(string)
}

// Status getter
func (feature *FeatureMap) Status() bool {
	status := (*feature)["Status"]
	return status.(bool)
}

// DiffFeatureMap return a map that contains the difference of status between the oldFeature and newFeature
func DiffFeatureMap(oldFeatures []FeatureMap, newFeatures []FeatureMap) map[string]int {
	result := make(map[string]int)

	//Add the features to the result assuming that the all new features are inactive
	for _, feature := range oldFeatures {
		if feature.Status() {
			result[feature.Name()] = DiffFeatureDeactivated
		} else {
			result[feature.Name()] = DiffFeatureInactive
		}
	}

	//Activate/deactivate missing feature
	for _, feature := range newFeatures {
		if (result[feature.Name()] == DiffFeatureInactive || result[feature.Name()] == DiffFeatureMissing) && !feature.Status() {
			result[feature.Name()] = DiffFeatureInactive
		} else if (result[feature.Name()] == DiffFeatureDeactivated) && !feature.Status() {
			result[feature.Name()] = DiffFeatureDeactivated
		} else if (result[feature.Name()] == DiffFeatureInactive || result[feature.Name()] == DiffFeatureMissing) && feature.Status() {
			result[feature.Name()] = DiffFeatureActivated
		} else if (result[feature.Name()] == DiffFeatureDeactivated) && feature.Status() {
			result[feature.Name()] = DiffFeatureActive
		}
	}

	return result
}
