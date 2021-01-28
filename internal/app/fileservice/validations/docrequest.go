package validations

import "fmt"

// ValidateJSONDocumentRequest ...
func ValidateJSONDocumentRequest(properties map[string]interface{}) error {
	docClassProperty := "class"
	docTypeProperty := "type"
	docNumProperty := "number"

	_, isDocClassExists := properties[docClassProperty]
	_, isDocTypeExists := properties[docTypeProperty]
	_, isDocNumExists := properties[docNumProperty]

	if !isDocClassExists {
		return fmt.Errorf("Bad request. Field %s does not exists", docClassProperty)
	}
	if !isDocTypeExists {
		return fmt.Errorf("Bad request. Field %s does not exists", docTypeProperty)
	}
	if !isDocNumExists {
		return fmt.Errorf("Bad request. Field %s does not exists", docNumProperty)
	}
	return nil
}
