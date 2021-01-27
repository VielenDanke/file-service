package service

// PrepareMetadata ...
func PrepareMetadata(properties map[string]interface{}, fieldsToRemove []string) {
	for _, v := range fieldsToRemove {
		_, ok := properties[v]
		if ok {
			delete(properties, v)
		}
	}
}
