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
