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

var _ = Describe("Linkedsecret controller GCP", func() {

	const (
		TIMEOUT  = time.Second * 60
		INTERVAL = time.Millisecond * 100
	)

	var gcpPlain LinkedSecretTest
	var gcpJSON LinkedSecretTest
	var gcpScheduleParseError LinkedSecretTest

	BeforeEach(func() {
		gcpPlain = LinkedSecretTest{
			name:      "google-example1",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:           "Google",
				ProviderDataFormat: "PLAIN",
				ProviderOptions:    map[string]string{"project": "project01-306719", "secret": "secret-plain-tst", "version": "latest"},
				SecretName:         "mysecret-google-example1",
				Schedule:           "@every 1s",
				Suspended:          false,
				Deployment:         "myapp",
			},
		}

		gcpJSON = LinkedSecretTest{
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

		gcpScheduleParseError = LinkedSecretTest{
			name:      "google-example3-schedule-parse-error",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:           "Google",
				ProviderDataFormat: "JSON",
				ProviderOptions:    map[string]string{"project": "project01-306719", "secret": "secret-json-tst", "version": "latest"},
				SecretName:         "mysecret-google-example3-schedule-parse-error",
				Schedule:           "@every 1ss",
				Suspended:          false,
			},
		}
	})

	Context("When creating new GCP PLAIN Linkedsecret", func() {
		It("Should create GCP PLAIN Linkedsecret", func() {

			By("By creating new GCP linkedsecret with PLAIN secret")
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
				return err == nil
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check expected spec
			By("By checking created linkedsecret spec")
			Expect(createdLinkedSecret.Spec.Provider).Should(Equal("Google"))
			Expect(createdLinkedSecret.Spec.ProviderDataFormat).Should(Equal("PLAIN"))
			Expect(createdLinkedSecret.Spec.ProviderOptions["project"]).Should(Equal("project01-306719"))
			Expect(createdLinkedSecret.Spec.ProviderOptions["secret"]).Should(Equal("secret-plain-tst"))
			Expect(createdLinkedSecret.Spec.ProviderOptions["version"]).Should(Equal("latest"))
			Expect(createdLinkedSecret.Spec.SecretName).Should(Equal("mysecret-google-example1"))
			Expect(createdLinkedSecret.Spec.Suspended).Should(Equal(false))
			Expect(createdLinkedSecret.Spec.Schedule).Should(Equal("@every 1s"))

			// Expected status
			By("By checking created linkedsecret status")
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

	Context("When updating GCP PLAIN Linkedsecret ", func() {
		It("Should suspend GCP PLAIN Linkedsecret", func() {

			By("Getting GCP PLAIN linkedsecret with GCP PLAIN Secret")
			ctx := context.Background()
			linkedSecretLookupKey := types.NamespacedName{Namespace: gcpPlain.namespace, Name: gcpPlain.name}
			linkedSecret := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, linkedSecret)
				return err == nil
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			By("Changing spec field 'suspended' to true")
			linkedSecret.Spec.Suspended = true

			// Create new LinkeSecret
			Expect(k8sClient.Update(ctx, linkedSecret)).Should(Succeed())

			updatedLinkedSecret := &securityv1.LinkedSecret{}

			By("By checking updated linkedsecret spec")
			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, updatedLinkedSecret)
				if err != nil {
					return false
				}
				if !updatedLinkedSecret.Spec.Suspended {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check expected spec
			Expect(updatedLinkedSecret.Spec.Provider).Should(Equal("Google"))
			Expect(updatedLinkedSecret.Spec.ProviderDataFormat).Should(Equal("PLAIN"))
			Expect(updatedLinkedSecret.Spec.ProviderOptions["project"]).Should(Equal("project01-306719"))
			Expect(updatedLinkedSecret.Spec.ProviderOptions["secret"]).Should(Equal("secret-plain-tst"))
			Expect(updatedLinkedSecret.Spec.ProviderOptions["version"]).Should(Equal("latest"))
			Expect(updatedLinkedSecret.Spec.SecretName).Should(Equal("mysecret-google-example1"))
			Expect(updatedLinkedSecret.Spec.Suspended).Should(Equal(true))
			Expect(updatedLinkedSecret.Spec.Schedule).Should(Equal("@every 1s"))

			By("By checking updated linkedsecret status")
			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, updatedLinkedSecret)
				if err != nil {
					return false
				}
				if updatedLinkedSecret.Status.CronJobStatus != JOBSUSPENDED {
					return false
				}
				if updatedLinkedSecret.Status.CronJobID != cron.EntryID(-1) {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())
			// Check expected status
			Expect(updatedLinkedSecret.Status.CreatedSecret).Should(Equal("mysecret-google-example1"))
			Expect(updatedLinkedSecret.Status.CronJobID).Should(Equal(cron.EntryID(-1)))
			Expect(updatedLinkedSecret.Status.CronJobStatus).Should(Equal(JOBSUSPENDED))
			Expect(updatedLinkedSecret.Status.CurrentProvider).Should(Equal("Google"))
			Expect(updatedLinkedSecret.Status.CurrentSchedule).Should(Equal("@every 1s"))
			Expect(updatedLinkedSecret.Status.CurrentProviderOptions["project"]).Should(Equal("project01-306719"))
			Expect(updatedLinkedSecret.Status.CurrentProviderOptions["secret"]).Should(Equal("secret-plain-tst"))
			Expect(updatedLinkedSecret.Status.CurrentProviderOptions["version"]).Should(Equal("latest"))

		})

	})

	Context("When creating new GCP JSON Linkedsecret", func() {
		It("Should create GCP JSON Linkedsecret", func() {

			By("By creating new GCP linkedsecret")
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

			By("By checking created linkedsecret spec")
			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
				return err == nil
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

	Context("When deleting GCP JSON Linkedsecret", func() {
		It("Should delete GCP JSON Linkedsecret", func() {

			By("By delete GCP linkedsecret")
			ctx := context.Background()

			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: gcpJSON.name, Namespace: gcpJSON.namespace},
				Spec:       gcpJSON.spec,
			}
			// Create new LinkeSecret
			Expect(k8sClient.Delete(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: gcpJSON.namespace, Name: gcpJSON.name}
			deleteLinkedSecret := &securityv1.LinkedSecret{}

			By("By checking deleted linkedsecret")
			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, deleteLinkedSecret)
				return err != nil
			}, TIMEOUT, INTERVAL).Should(BeFalse())

		})
	})

	Context("When Creating GCP JSON Linkedsecret with schedule parse error", func() {
		It("Should create a GCP Linkedsecret with schedule parse error", func() {

			By("By Create GCP linkedsecret with schedule parse error")
			ctx := context.Background()

			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: gcpScheduleParseError.name, Namespace: gcpScheduleParseError.namespace},
				Spec:       gcpScheduleParseError.spec,
			}

			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: gcpScheduleParseError.namespace, Name: gcpScheduleParseError.name}
			createdLinkedSecretWithScheduleParseError := &securityv1.LinkedSecret{}

			By("By checking Created linkedsecret with schedule parse error")
			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecretWithScheduleParseError)
				if err != nil {
					return false
				}
				if createdLinkedSecretWithScheduleParseError.Status.CronJobStatus != JOBFAILPARSESCHEDULE {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check Current status
			Expect(createdLinkedSecretWithScheduleParseError.Status.CreatedSecret).Should(Equal("mysecret-google-example3-schedule-parse-error"))
			Expect(createdLinkedSecretWithScheduleParseError.Status.CronJobID).Should(Equal(cron.EntryID(-1)))
			Expect(createdLinkedSecretWithScheduleParseError.Status.CronJobStatus).Should(Equal(JOBFAILPARSESCHEDULE))
			Expect(createdLinkedSecretWithScheduleParseError.Status.CurrentProvider).Should(Equal("Google"))
			Expect(createdLinkedSecretWithScheduleParseError.Status.CurrentSchedule).Should(Equal("@every 1ss"))
			Expect(createdLinkedSecretWithScheduleParseError.Status.CurrentProviderOptions["project"]).Should(Equal("project01-306719"))
			Expect(createdLinkedSecretWithScheduleParseError.Status.CurrentProviderOptions["secret"]).Should(Equal("secret-json-tst"))
			Expect(createdLinkedSecretWithScheduleParseError.Status.CurrentProviderOptions["version"]).Should(Equal("latest"))

		})
	})

})
