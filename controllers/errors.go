package controllers

import "fmt"

// EmptySecret error
type EmptySecret struct {
	name string
	err  string
}

func (e *EmptySecret) Error() string {
	return fmt.Sprintf("Cloud secret %s: %s.", e.name, e.err)
}

// InvalidCloudSecret error
type InvalidCloudSecret struct {
}

func (i *InvalidCloudSecret) Error() string {
	return "Invalid Cloud Secret data format."
}
