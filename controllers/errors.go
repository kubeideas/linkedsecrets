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
	return "AWS Provider option 'region' not informed or invalid."
}

// InvalidGoogleCloudProject error
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

// InvalidIBMSecretManagerID error
type InvalidIBMSecretManagerID struct {
}

func (i *InvalidIBMSecretManagerID) Error() string {
	return "IBM Provider option 'secretManagerInstanceId' not informed or invalid."
}

// InvalidIBMRegion error
type InvalidIBMRegion struct {
}

func (i *InvalidIBMRegion) Error() string {
	return "IBM Provider option 'region' not informed or invalid."
}

// InvalidIBMServiceApiKey error
type InvalidIBMServiceApiKey struct {
}

func (i *InvalidIBMServiceApiKey) Error() string {
	return "IBM environment variable [ IBM_SERVICE_API_KEY ] was not defined or invalid."
}

// InvalidUUIDFormat error
type InvalidUUIDFormat struct {
	fieldName string
}

func (i *InvalidUUIDFormat) Error() string {
	return fmt.Sprintf("[ %s ] has invalid UUID format", i.fieldName)
}

// DeploymentNotDefined error
type DeploymentNotDefined struct {
}

func (i *DeploymentNotDefined) Error() string {
	return "Deployment for rollout restart not defined."
}

// DeploymentNotFound error
type DeploymentNotFound struct {
	name string
}

func (i *DeploymentNotFound) Error() string {
	return fmt.Sprintf("Deployment [ %s ] not found", i.name)
}
