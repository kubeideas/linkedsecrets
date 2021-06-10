/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"github.com/robfig/cron/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LinkedSecretSpec defines the desired state of LinkedSecret
type LinkedSecretSpec struct {
	// +kubebuilder:validation:Enum={"Google","AWS","Azure"}
	// +kubebuilder:validation:Required
	// Supported cloud secret manager. Valid options: Google,AWS,Azure.
	Provider string `json:"provider"`

	// +kubebuilder:validation:Enum={"PLAIN", "JSON"}
	// +kubebuilder:validation:Required
	// Supported formats: PLAIN and JSON
	// "PLAIN" format key/value must be delimited by character "=".
	// Empty lines, key without value and value without key will be skipped.
	// Leading and trailing whitespaces will be ignored. Ex: password=pass12@#=+$% or password = pass12@#=+$% (with whitespaces).
	// "JSON" format must be key/value format. Ex: {"pasword":"pass12@#=+$%","host":"myhost"}.
	ProviderDataFormat string `json:"providerDataFormat"`

	// +optional
	// Extra options necessary to fetch Cloud secret. If version is omitted, secret latest version will be used regardeless of Cloud provider.
	// Example GCP: project: <PROJECT-ID>, secret: <GCP-SECRET-NAME>, version: <latest|"1"|"2"|...>.
	// Example AWS: region: <AWS-REGION>, secret: <AWS-SECRET-NAME>, version: <AWSPREVIOUS|AWSCURRENT>.
	// Example Azure: keyvault: <KEYVAULT-NAME>, secret: <AZURE-SECRET-NAME>, version: <AZURE-SECRET-VERSION-ID>.
	ProviderOptions map[string]string `json:"providerOptions,omitempty"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*
	// Secret name expected to be created into kubernetes with data fetched from Cloud secret manager solution.
	SecretName string `json:"secretName"`

	// +kubebuilder:validation:Type=string
	// +optional
	// Schedule define interval to synchronize cloud secrets data and kubernetes secrets.
	// Examples of valid schedule: "@every 120s"(every 2 minutes), "@every 1m30s" (every 1 minute and 30 seconds),
	// "@every 10m" (every 10 minutes), "@every 1h" (every hour), "*/5 * * * * *" (every 5 minutes).
	// If empty schedule will be considered disabled and will be synchronized just on creation.
	// [IMPORTANT]: Please mind the interval you have chosen for data synchronization and
	// check Secret Manager pricing details in order to avoid unneeded cloud costs.
	Schedule string `json:"schedule,omitempty"`

	// +kubebuilder:validation:Type=boolean
	// +optional
	// Use this field to suspend cronjob temporarily. Valid values: {true, false}
	Suspended bool `json:"suspended,omitempty"`

	// +kubebuilder:validation:Type=boolean
	// +optional
	// Use this field keep secret after LinkedSecret deletion. Valid values: {true, false}
	KeepSecretOnDelete bool `json:"keepSecretOnDelete,omitempty"`

	// +kubebuilder:validation:Type=string
	// Deployment name which pods will be restarted if secret data key or value were changed.
	// Pods rollout update will happen 5 seconds after secret was updated.
	Deployment string `json:"deployment,omitempty"`
}

// LinkedSecretStatus defines the observed state of LinkedSecret
type LinkedSecretStatus struct {
	// if "CurrentSecretStatus = Synched" data between cloud provider and kubernetes secret were synchronized.
	// if "CurrentSecretStatus = NotSynched" may have occured an error during synchronization process.
	// Please check linkedsecret events for more details.
	CurrentSecretStatus string `json:"createdSecretStatus,omitempty"`

	// Secret name currently being used.
	CreatedSecret string `json:"createdSecret,omitempty"`

	// Secret namespace currently being used.
	CreatedSecretNamespace string `json:"createdSecretNamespace,omitempty"`

	// Last time secret was synchronized.
	LastScheduleExecution *metav1.Time `json:"lastScheduleExecution,omitempty"`
	//NextScheduleExecution  *metav1.Time `json:"nextScheduleExecution,omitempty"`

	// Provider name currently being used.
	CurrentProvider string `json:"currentProvider"`

	// Provider options currently being used.
	CurrentProviderOptions map[string]string `json:"currentProviderOptions,omitempty"`

	// Cronjob current status.
	// If "CronJobStatus = Scheduled" job schedule is normal.
	// If "CronJobStatus = NotScheduled" may have occured an error during schedule process, schedule is empty or schedule format is invalid.
	// Please check linkedsecret events for more details.
	CronJobStatus string `json:"cronJobStatus,omitempty"`

	// Cronjob current ID.
	// "If CronJobID > 0", job schedule is normal.
	// "If CronJobID = -1", may have occured an error during schedule process, schedule is empty or schedule format is invalid.
	// Please check linkedsecret events for more details.
	CronJobID cron.EntryID `json:"cronJobID,omitempty"`

	// Cronjob current schedule.
	CurrentSchedule string `json:"currentSchedule,omitempty"`

	// If value is "true" secret wont be deleted after LinkedSecret deletion.
	KeepSecretOnDelete bool `json:"keepSecretOnDelete,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.currentProvider",name="provider",type="string"
// +kubebuilder:printcolumn:JSONPath=".status.createdSecretStatus",name="status",type="string"
// +kubebuilder:printcolumn:JSONPath=".status.createdSecret",name="secret",type="string"
// +kubebuilder:printcolumn:JSONPath=".status.keepSecretOnDelete",name="keep-on-delete",type="boolean"
// +kubebuilder:printcolumn:JSONPath=".status.lastScheduleExecution",name="last-sync",type="string"
// +kubebuilder:printcolumn:JSONPath=".status.cronJobStatus",name="cron-job-status",type="string"
// +kubebuilder:printcolumn:JSONPath=".status.currentSchedule",name="current-schedule",type="string"
// +kubebuilder:resource:shortName=lns
// +kubebuilder:storageversion

// LinkedSecret is the Schema for the linkedsecrets API
type LinkedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LinkedSecretSpec   `json:"spec,omitempty"`
	Status LinkedSecretStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LinkedSecretList contains a list of LinkedSecret
type LinkedSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LinkedSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LinkedSecret{}, &LinkedSecretList{})
}
