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

// InvalidAzureKeyvault error
type InvalidAzureKeyvault struct {
}

func (i *InvalidAzureKeyvault) Error() string {
	return "Provider option 'keyvault' not informed or invalid"
}

// InvalidAWSRegion error
type InvalidAWSRegion struct {
}

func (i *InvalidAWSRegion) Error() string {
	return "Provider option 'region' not informed or invalid."
}

// InvalidAWSRegion error
type InvalidGoogleCloudProject struct {
}

func (i *InvalidGoogleCloudProject) Error() string {
	return "Provider option 'project' not informed or invalid."
}

// InvalidSecretOption error
type InvalidSecretOption struct {
}

func (i *InvalidSecretOption) Error() string {
	return "Provider option 'secret' not informed or invalid."
}
