package validator

func IsValidData(featureID int, tagIDs []int) bool {
	if featureID < 0 {
		return false
	}
	for _, tagID := range tagIDs {
		if tagID < 0 {
			return false
		}
	}
	return true
}
