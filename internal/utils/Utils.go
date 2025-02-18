package utils

// Checks if a string exists in an array of strings.
func IsArrayContainString(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func ConvertInterfaceArrayToStringArray(interfaceArray []interface{}) []string {
	var stringSlice []string
	for _, v := range interfaceArray {
		if str, ok := v.(string); ok {
			stringSlice = append(stringSlice, str)
		}
	}

	return stringSlice
}
