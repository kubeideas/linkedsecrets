package controllers

import (
	"context"
	securityv1 "kubeideas/linkedsecrets/api/v1"
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
		DURATION = time.Second * 10
		INTERVAL = time.Millisecond * 250
	)

	var (
		gcpSecretPlain        LinkedSecretTest
		gcpSecretJSON         LinkedSecretTest
		gcpSecretDocker       LinkedSecretTest
		gcpScheduleParseError LinkedSecretTest
	)

	BeforeEach(func() {
		gcpSecretPlain = LinkedSecretTest{
			name:      "google-example1",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "Google",
				ProviderSecretFormat: "PLAIN",
				ProviderOptions:      map[string]string{"project": "linkedsecrets", "secret": "opaque-secret-plain", "version": "latest"},
				SecretName:           "mysecret-google-example1",
				Schedule:             "@every 10s",
				Suspended:            false,
				RolloutRestartDeploy: "myapp",
			},
		}

		gcpSecretJSON = LinkedSecretTest{
			name:      "google-example2",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "Google",
				ProviderSecretFormat: "JSON",
				ProviderOptions:      map[string]string{"project": "linkedsecrets", "secret": "opaque-secret-json", "version": "latest"},
				SecretName:           "mysecret-google-example2",
				Schedule:             "@every 10s",
				Suspended:            false,
			},
		}

		gcpSecretDocker = LinkedSecretTest{
			name:      "google-example3",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "Google",
				ProviderSecretFormat: "JSON",
				ProviderOptions:      map[string]string{"project": "linkedsecrets", "secret": "docker-secret-json", "version": "latest"},
				SecretName:           "mysecret-google-example3",
				Schedule:             "@every 10s",
				Suspended:            false,
			},
		}

		gcpScheduleParseError = LinkedSecretTest{
			name:      "google-example4",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "Google",
				ProviderSecretFormat: "JSON",
				ProviderOptions:      map[string]string{"project": "linkedsecrets", "secret": "opaque-secret-json", "version": "latest"},
				SecretName:           "mysecret-google-example4",
				Schedule:             "@every 10ss",
				Suspended:            false,
			},
		}

	})

	Context("When creating Linkedsecret google-example1", func() {
		It("Should create Linkedsecret google-example1", func() {

			By("Creating Linkedsecret google-example1")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: gcpSecretPlain.name, Namespace: gcpSecretPlain.namespace},
				Spec:       gcpSecretPlain.spec,
			}

			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())
			linkedSecretLookupKey := types.NamespacedName{Namespace: gcpSecretPlain.namespace, Name: gcpSecretPlain.name}
			googleExample1 := &securityv1.LinkedSecret{}

			// Get created linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, googleExample1)
				if err != nil {
					return false
				}
				if googleExample1.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check spec
			By("Checking google-example1 spec")
			Expect(googleExample1.Spec.Provider).Should(Equal("Google"))
			Expect(googleExample1.Spec.ProviderSecretFormat).Should(Equal("PLAIN"))
			Expect(googleExample1.Spec.ProviderOptions["project"]).Should(Equal("linkedsecrets"))
			Expect(googleExample1.Spec.ProviderOptions["secret"]).Should(Equal("opaque-secret-plain"))
			Expect(googleExample1.Spec.ProviderOptions["version"]).Should(Equal("latest"))
			Expect(googleExample1.Spec.SecretName).Should(Equal("mysecret-google-example1"))
			Expect(googleExample1.Spec.Suspended).Should(Equal(false))
			Expect(googleExample1.Spec.Schedule).Should(Equal("@every 10s"))

			// Check status
			By("Checking google-example1 status")
			Expect(googleExample1.Status.CurrentSecret).Should(Equal(googleExample1.Spec.SecretName))
			Expect(googleExample1.Status.CurrentSecretStatus).Should(Equal(STATUSSYNCHED))
			Expect(googleExample1.Status.CronJobStatus).Should(Equal(JOBSCHEDULED))
			Expect(googleExample1.Status.ObservedGeneration).Should(Equal(googleExample1.GetGeneration()))
			Expect(googleExample1.Status.CronJobID).Should(Equal(cron.EntryID(googleExample1.Status.CronJobID)))

		})
	})

	Context("When creating Linkedsecret google-example2 with JSON Data Secret", func() {
		It("Should create Linkedsecret google-example2", func() {

			By("Creating Linkedsecret google-example2")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: gcpSecretJSON.name, Namespace: gcpSecretJSON.namespace},
				Spec:       gcpSecretJSON.spec,
			}

			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())
			linkedSecretLookupKey := types.NamespacedName{Namespace: gcpSecretJSON.namespace, Name: gcpSecretJSON.name}
			googleExample2 := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, googleExample2)
				if err != nil {
					return false
				}
				if googleExample2.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check spec
			By("Checking google-example2 spec")
			Expect(googleExample2.Spec.Provider).Should(Equal("Google"))
			Expect(googleExample2.Spec.ProviderSecretFormat).Should(Equal("JSON"))
			Expect(googleExample2.Spec.ProviderOptions["project"]).Should(Equal("linkedsecrets"))
			Expect(googleExample2.Spec.ProviderOptions["secret"]).Should(Equal("opaque-secret-json"))
			Expect(googleExample2.Spec.ProviderOptions["version"]).Should(Equal("latest"))
			Expect(googleExample2.Spec.SecretName).Should(Equal("mysecret-google-example2"))
			Expect(googleExample2.Spec.Suspended).Should(Equal(false))
			Expect(googleExample2.Spec.Schedule).Should(Equal("@every 10s"))

			// Check status
			By("Checking google-example2 status")
			Expect(googleExample2.Status.CurrentSecret).Should(Equal(googleExample2.Spec.SecretName))
			Expect(googleExample2.Status.CurrentSecretStatus).Should(Equal(STATUSSYNCHED))
			Expect(googleExample2.Status.CronJobStatus).Should(Equal(JOBSCHEDULED))
			Expect(googleExample2.Status.ObservedGeneration).Should(Equal(googleExample2.GetGeneration()))
			Expect(googleExample2.Status.CronJobID).Should(Equal(cron.EntryID(googleExample2.Status.CronJobID)))

		})
	})

	Context("When updating Linkedsecret google-example2", func() {
		It("Should suspend Linkedsecret google-example2 schedule", func() {

			By("Getting google-example2 linkedsecret")
			ctx := context.Background()
			linkedSecretLookupKey := types.NamespacedName{Namespace: gcpSecretJSON.namespace, Name: gcpSecretJSON.name}
			googleExample2 := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, googleExample2)
				return err == nil
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			By("Changing spec field 'suspended' to true")
			googleExample2.Spec.Suspended = true

			Expect(k8sClient.Update(ctx, googleExample2)).Should(Succeed())
			updatedGoogleExample2 := &securityv1.LinkedSecret{}

			By("Checking updated google-example2 spec")
			// Get updated linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, updatedGoogleExample2)
				if err != nil {
					return false
				}
				if !updatedGoogleExample2.Spec.Suspended {
					return false
				}

				if updatedGoogleExample2.Status.CronJobStatus != JOBSUSPENDED {
					return false
				}

				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check spec
			By("Checking updated google-example2 spec")
			Expect(updatedGoogleExample2.Spec.Provider).Should(Equal("Google"))
			Expect(updatedGoogleExample2.Spec.ProviderSecretFormat).Should(Equal("JSON"))
			Expect(updatedGoogleExample2.Spec.ProviderOptions["project"]).Should(Equal("linkedsecrets"))
			Expect(updatedGoogleExample2.Spec.ProviderOptions["secret"]).Should(Equal("opaque-secret-json"))
			Expect(updatedGoogleExample2.Spec.ProviderOptions["version"]).Should(Equal("latest"))
			Expect(updatedGoogleExample2.Spec.SecretName).Should(Equal("mysecret-google-example2"))
			Expect(updatedGoogleExample2.Spec.Suspended).Should(Equal(true))
			Expect(updatedGoogleExample2.Spec.Schedule).Should(Equal("@every 10s"))

			// Check status
			By("Checking updated google-example2 status")
			Expect(updatedGoogleExample2.Status.CurrentSecret).Should(Equal(updatedGoogleExample2.Spec.SecretName))
			Expect(updatedGoogleExample2.Status.CurrentSecretStatus).Should(Equal(STATUSSYNCHED))
			Expect(updatedGoogleExample2.Status.CronJobStatus).Should(Equal(JOBSUSPENDED))
			Expect(updatedGoogleExample2.Status.ObservedGeneration).Should(Equal(updatedGoogleExample2.GetGeneration()))
			Expect(updatedGoogleExample2.Status.CronJobID).Should(Equal(cron.EntryID(updatedGoogleExample2.Status.CronJobID)))

		})
	})

	Context("When deleting Linkedsecret google-example2", func() {
		It("Should delete Linkedsecret google-example2 and secret mysecret-google-example2", func() {

			By("Getting google-example2 Linkedsecret")
			ctx := context.Background()

			googleExample2 := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: gcpSecretJSON.name, Namespace: gcpSecretJSON.namespace},
				Spec:       gcpSecretJSON.spec,
			}
			// Create new LinkeSecret
			Expect(k8sClient.Delete(ctx, googleExample2)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: gcpSecretJSON.namespace, Name: gcpSecretJSON.name}
			deleteLinkedSecret := &securityv1.LinkedSecret{}

			By("Checking deleted linkedsecret")
			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, deleteLinkedSecret)
				return err != nil
			}, TIMEOUT, INTERVAL).Should(BeFalse())

		})
	})

	Context("When creating Linkedsecret google-example3 with Docker config JSON", func() {
		It("Should create LInkedsecret Docker config JSON data", func() {

			By("Creating Linkedsecret google-example3 with Docker config JSON")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: gcpSecretDocker.name, Namespace: gcpSecretDocker.namespace},
				Spec:       gcpSecretDocker.spec,
			}
			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: gcpSecretDocker.namespace, Name: gcpSecretDocker.name}
			googleExample3 := &securityv1.LinkedSecret{}

			By("Getting google-example3 Linkedsecret")
			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, googleExample3)
				if err != nil {
					return false
				}
				if googleExample3.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check expected spec
			By("Checking google-example3 spec")
			Expect(googleExample3.Spec.Provider).Should(Equal("Google"))
			Expect(googleExample3.Spec.ProviderSecretFormat).Should(Equal("JSON"))
			Expect(googleExample3.Spec.ProviderOptions["project"]).Should(Equal("linkedsecrets"))
			Expect(googleExample3.Spec.ProviderOptions["secret"]).Should(Equal("docker-secret-json"))
			Expect(googleExample3.Spec.ProviderOptions["version"]).Should(Equal("latest"))
			Expect(googleExample3.Spec.SecretName).Should(Equal("mysecret-google-example3"))
			Expect(googleExample3.Spec.Suspended).Should(Equal(false))
			Expect(googleExample3.Spec.Schedule).Should(Equal("@every 10s"))

			// Check status
			By("Checking google-example3 status")
			Expect(googleExample3.Status.CurrentSecret).Should(Equal(googleExample3.Spec.SecretName))
			Expect(googleExample3.Status.CurrentSecretStatus).Should(Equal(STATUSSYNCHED))
			Expect(googleExample3.Status.CronJobStatus).Should(Equal(JOBSCHEDULED))
			Expect(googleExample3.Status.ObservedGeneration).Should(Equal(googleExample3.GetGeneration()))
			Expect(googleExample3.Status.CronJobID).Should(Equal(cron.EntryID(googleExample3.Status.CronJobID)))
		})
	})

	Context("When Creating Linkedsecret google-example4 with schedule parse error", func() {
		It("Should create Linkedsecret google-example4 with schedule parse error", func() {

			By("Creating Linkedsecret google-example4")
			ctx := context.Background()

			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: gcpScheduleParseError.name, Namespace: gcpScheduleParseError.namespace},
				Spec:       gcpScheduleParseError.spec,
			}

			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: gcpScheduleParseError.namespace, Name: gcpScheduleParseError.name}
			googleExample4 := &securityv1.LinkedSecret{}

			By("Checking Linkedsecret google-example4 with schedule parse error")
			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, googleExample4)
				if err != nil {
					return false
				}
				if googleExample4.Status.CronJobStatus != JOBFAILPARSESCHEDULE {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check Current status
			Expect(googleExample4.Status.CurrentSecret).Should(Equal("mysecret-google-example4"))
			Expect(googleExample4.Status.CronJobID).Should(Equal(cron.EntryID(googleExample4.Status.CronJobID)))
			Expect(googleExample4.Status.CronJobStatus).Should(Equal(JOBFAILPARSESCHEDULE))
			Expect(googleExample4.Status.CurrentSchedule).Should(Equal("@every 10ss"))

		})
	})

})
