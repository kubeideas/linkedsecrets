package controllers

// check if a slice contain string
func containString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// remove string from a slice
func removeString(slice []string, s string) []string {
	for i, item := range slice {
		if item == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
