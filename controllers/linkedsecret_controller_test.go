/*


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

package controllers

import (
	"context"
	securityv1 "linkedsecrets/api/v1"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robfig/cron/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// encapsulate Linkedsecret name and spec
type linkedSecretTest struct {
	name      string
	namespace string
	spec      securityv1.LinkedSecretSpec
}

var _ = Describe("Linkedsecret controller", func() {

	const (
		TIMEOUT  = time.Second * 60
		INTERVAL = time.Millisecond * 100
	)

	var gcpPlain linkedSecretTest
	var gcpJSON linkedSecretTest

	BeforeEach(func() {
		gcpPlain = linkedSecretTest{
			name:      "google-example1",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:           "Google",
				ProviderDataFormat: "PLAIN",
				ProviderOptions:    map[string]string{"project": "project01-306719", "secret": "secret-plain-tst", "version": "latest"},
				SecretName:         "mysecret-google-example1",
				Schedule:           "@every 1s",
				Suspended:          false,
			},
		}

		gcpJSON = linkedSecretTest{
			name:      "google-example2",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:           "Google",
				ProviderDataFormat: "JSON",
				ProviderOptions:    map[string]string{"project": "project01-306719", "secret": "secret-json-tst", "version": "latest"},
				SecretName:         "mysecret-google-example2",
				Schedule:           "@every 1s",
				Suspended:          false,
			},
		}
	})

	Describe("Creating GCP Plain Linkedsecret", func() {
		Context("Creating new Linkedsecret sinchronizing data with GCP", func() {
			It("Should be GCP Linkedsecret", func() {
				By("Creating new GCP linkedsecret")
				ctx := context.Background()
				linkedSecret := &securityv1.LinkedSecret{
					TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
					ObjectMeta: v1.ObjectMeta{Name: gcpPlain.name, Namespace: gcpPlain.namespace},
					Spec:       gcpPlain.spec,
				}
				// Create new LinkeSecret
				Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

				linkedSecretLookupKey := types.NamespacedName{Namespace: gcpPlain.namespace, Name: gcpPlain.name}
				createdLinkedSecret := &securityv1.LinkedSecret{}

				// Get linkedSecret
				Eventually(func() bool {
					err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
					if err != nil {
						return false
					}
					return true
				}, TIMEOUT, INTERVAL).Should(BeTrue())

				// Check expected spec
				Expect(createdLinkedSecret.Spec.Provider).Should(Equal("Google"))
				Expect(createdLinkedSecret.Spec.ProviderDataFormat).Should(Equal("PLAIN"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["project"]).Should(Equal("project01-306719"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["secret"]).Should(Equal("secret-plain-tst"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["version"]).Should(Equal("latest"))
				Expect(createdLinkedSecret.Spec.SecretName).Should(Equal("mysecret-google-example1"))
				Expect(createdLinkedSecret.Spec.Suspended).Should(Equal(false))
				Expect(createdLinkedSecret.Spec.Schedule).Should(Equal("@every 1s"))

				// Expected status
				By("By checking linkedsecret status")
				// Get linkedSecret
				Eventually(func() bool {
					err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
					if err != nil {
						return false
					}
					if createdLinkedSecret.Status.CronJobID == cron.EntryID(0) {
						return false
					}
					return true
				}, TIMEOUT, INTERVAL).Should(BeTrue())

				// Check expected status
				Expect(createdLinkedSecret.Status.CreatedSecret).Should(Equal("mysecret-google-example1"))
				Expect(createdLinkedSecret.Status.CronJobID).Should(Equal(cron.EntryID(1)))
				Expect(createdLinkedSecret.Status.CronJobStatus).Should(Equal("Scheduled"))
				Expect(createdLinkedSecret.Status.CurrentProvider).Should(Equal("Google"))
				Expect(createdLinkedSecret.Status.CurrentSchedule).Should(Equal("@every 1s"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["project"]).Should(Equal("project01-306719"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["secret"]).Should(Equal("secret-plain-tst"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["version"]).Should(Equal("latest"))

			})
		})
	})

	Describe("Creating GCP JSON Linkedsecret", func() {
		Context("Creating new Linkedsecret sinchronizing data with GCP", func() {
			It("Should be GCP Linkedsecret", func() {
				By("Creating new GCP linkedsecret")
				ctx := context.Background()
				linkedSecret := &securityv1.LinkedSecret{
					TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
					ObjectMeta: v1.ObjectMeta{Name: gcpJSON.name, Namespace: gcpJSON.namespace},
					Spec:       gcpJSON.spec,
				}
				// Create new LinkeSecret
				Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

				linkedSecretLookupKey := types.NamespacedName{Namespace: gcpJSON.namespace, Name: gcpJSON.name}
				createdLinkedSecret := &securityv1.LinkedSecret{}

				// Get linkedSecret
				Eventually(func() bool {
					err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
					if err != nil {
						return false
					}
					return true
				}, TIMEOUT, INTERVAL).Should(BeTrue())

				// Check expected spec
				Expect(createdLinkedSecret.Spec.Provider).Should(Equal("Google"))
				Expect(createdLinkedSecret.Spec.ProviderDataFormat).Should(Equal("JSON"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["project"]).Should(Equal("project01-306719"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["secret"]).Should(Equal("secret-json-tst"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["version"]).Should(Equal("latest"))
				Expect(createdLinkedSecret.Spec.SecretName).Should(Equal("mysecret-google-example2"))
				Expect(createdLinkedSecret.Spec.Suspended).Should(Equal(false))
				Expect(createdLinkedSecret.Spec.Schedule).Should(Equal("@every 1s"))

				// Check expected status
				By("By checking linkedsecret status")
				// Get linkedSecret
				Eventually(func() bool {
					err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
					if err != nil {
						return false
					}
					if createdLinkedSecret.Status.CronJobID == cron.EntryID(0) {
						return false
					}
					return true
				}, TIMEOUT, INTERVAL).Should(BeTrue())

				// Check Current status
				Expect(createdLinkedSecret.Status.CreatedSecret).Should(Equal("mysecret-google-example2"))
				Expect(createdLinkedSecret.Status.CronJobID).Should(Equal(cron.EntryID(1)))
				Expect(createdLinkedSecret.Status.CronJobStatus).Should(Equal("Scheduled"))
				Expect(createdLinkedSecret.Status.CurrentProvider).Should(Equal("Google"))
				Expect(createdLinkedSecret.Status.CurrentSchedule).Should(Equal("@every 1s"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["project"]).Should(Equal("project01-306719"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["secret"]).Should(Equal("secret-json-tst"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["version"]).Should(Equal("latest"))

			})
		})
	})
})
