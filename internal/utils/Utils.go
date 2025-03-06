package utils

import "slices"

// Checks if a string exists in an array of strings.
func IsArrayContainString(slice []string, str string) bool {
	return slices.Contains(slice, str)
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
