package utils

func MapDereference(data map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			data[key] = MapDereference(v)
		case *interface{}:
			data[key] = *v
		}
	}
	return data
}

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
