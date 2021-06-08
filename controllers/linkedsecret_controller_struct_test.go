package controllers

import (
	securityv1 "linkedsecrets/api/v1"
)

// encapsulate Linkedsecret name and spec
type LinkedSecretTest struct {
	name      string
	namespace string
	spec      securityv1.LinkedSecretSpec
}
