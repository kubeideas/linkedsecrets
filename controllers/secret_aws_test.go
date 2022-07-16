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

var _ = Describe("Linkedsecret controller AWS", func() {

	const (
		TIMEOUT  = time.Second * 60
		DURATION = time.Second * 10
		INTERVAL = time.Millisecond * 250
	)

	var (
		awsPlain          LinkedSecretTest
		awsJSON           LinkedSecretTest
		awsSecretNotFound LinkedSecretTest
	)

	BeforeEach(func() {

		awsJSON = LinkedSecretTest{
			name:      "aws-example1",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "AWS",
				ProviderSecretFormat: "JSON",
				ProviderOptions:      map[string]string{"secret": "opaque-secret-json", "region": "us-east-1", "version": "AWSCURRENT"},
				SecretName:           "mysecret-aws-example1",
				Schedule:             "@every 10s",
				Suspended:            false,
			},
		}

		awsPlain = LinkedSecretTest{
			name:      "aws-example2",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "AWS",
				ProviderSecretFormat: "PLAIN",
				ProviderOptions:      map[string]string{"secret": "opaque-secret-plain", "region": "us-east-1", "version": "AWSCURRENT"},
				SecretName:           "mysecret-aws-example2",
				Schedule:             "@every 10s",
				Suspended:            false,
			},
		}

		awsSecretNotFound = LinkedSecretTest{
			name:      "aws-example3",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "AWS",
				ProviderSecretFormat: "PLAIN",
				ProviderOptions:      map[string]string{"secret": "secret-not-found", "region": "us-east-1", "version": "AWSCURRENT"},
				SecretName:           "mysecret-aws-example3",
				Schedule:             "@every 10s",
				Suspended:            false,
			},
		}
	})

	Context("When creating Linkedsecret aws-example1", func() {

		It("Should create Linkedsecret aws-example1", func() {

			By("Creating aws-example1")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: awsJSON.name, Namespace: awsJSON.namespace},
				Spec:       awsJSON.spec,
			}

			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: awsJSON.namespace, Name: awsJSON.name}
			awsExample1 := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, awsExample1)
				if err != nil {
					return false
				}
				if awsExample1.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check spec
			Expect(awsExample1.Spec.Provider).Should(Equal("AWS"))
			Expect(awsExample1.Spec.ProviderSecretFormat).Should(Equal("JSON"))
			Expect(awsExample1.Spec.ProviderOptions["secret"]).Should(Equal("opaque-secret-json"))
			Expect(awsExample1.Spec.ProviderOptions["region"]).Should(Equal("us-east-1"))
			Expect(awsExample1.Spec.ProviderOptions["version"]).Should(Equal("AWSCURRENT"))
			Expect(awsExample1.Spec.SecretName).Should(Equal("mysecret-aws-example1"))
			Expect(awsExample1.Spec.Suspended).Should(Equal(false))
			Expect(awsExample1.Spec.Schedule).Should(Equal("@every 10s"))

			// Check status
			Expect(awsExample1.Status.CurrentSecret).Should(Equal("mysecret-aws-example1"))
			Expect(awsExample1.Status.CronJobID).Should(Equal(cron.EntryID(awsExample1.Status.CronJobID)))
			Expect(awsExample1.Status.CronJobStatus).Should(Equal("Scheduled"))
			Expect(awsExample1.Status.CurrentSchedule).Should(Equal("@every 10s"))

		})
	})

	Context("When Creating Linkedsecret aws-example2", func() {
		It("Should create Linkedsecret aws-example2", func() {

			By("Creating aws-example2")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: awsPlain.name, Namespace: awsPlain.namespace},
				Spec:       awsPlain.spec,
			}
			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: awsPlain.namespace, Name: awsPlain.name}
			awsExample2 := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, awsExample2)
				if err != nil {
					return false
				}
				if awsExample2.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check spec
			Expect(awsExample2.Spec.Provider).Should(Equal("AWS"))
			Expect(awsExample2.Spec.ProviderSecretFormat).Should(Equal("PLAIN"))
			Expect(awsExample2.Spec.ProviderOptions["secret"]).Should(Equal("opaque-secret-plain"))
			Expect(awsExample2.Spec.ProviderOptions["region"]).Should(Equal("us-east-1"))
			Expect(awsExample2.Spec.ProviderOptions["version"]).Should(Equal("AWSCURRENT"))
			Expect(awsExample2.Spec.SecretName).Should(Equal("mysecret-aws-example2"))
			Expect(awsExample2.Spec.Suspended).Should(Equal(false))
			Expect(awsExample2.Spec.Schedule).Should(Equal("@every 10s"))

			// Check status
			Expect(awsExample2.Status.CurrentSecret).Should(Equal("mysecret-aws-example2"))
			Expect(awsExample2.Status.CronJobID).Should(Equal(cron.EntryID(awsExample2.Status.CronJobID)))
			Expect(awsExample2.Status.CronJobStatus).Should(Equal("Scheduled"))
			Expect(awsExample2.Status.CurrentSchedule).Should(Equal("@every 10s"))

		})
	})

	Context("When Creating Linkedsecret aws-example3", func() {
		It("Should create Linkedsecret aws-example3", func() {

			By("Creating aws-example3")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: awsSecretNotFound.name, Namespace: awsSecretNotFound.namespace},
				Spec:       awsSecretNotFound.spec,
			}
			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: awsSecretNotFound.namespace, Name: awsSecretNotFound.name}
			awsExample3 := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, awsExample3)
				if err != nil {
					return false
				}
				if awsExample3.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check spec
			Expect(awsExample3.Spec.Provider).Should(Equal("AWS"))
			Expect(awsExample3.Spec.ProviderSecretFormat).Should(Equal("PLAIN"))
			Expect(awsExample3.Spec.ProviderOptions["secret"]).Should(Equal("secret-not-found"))
			Expect(awsExample3.Spec.ProviderOptions["region"]).Should(Equal("us-east-1"))
			Expect(awsExample3.Spec.ProviderOptions["version"]).Should(Equal("AWSCURRENT"))
			Expect(awsExample3.Spec.SecretName).Should(Equal("mysecret-aws-example3"))
			Expect(awsExample3.Spec.Suspended).Should(Equal(false))
			Expect(awsExample3.Spec.Schedule).Should(Equal("@every 10s"))

			// Check status
			Expect(awsExample3.Status.CurrentSecret).Should(Equal("mysecret-aws-example3"))
			Expect(awsExample3.Status.CronJobStatus).Should(Equal("NotScheduled"))
			Expect(awsExample3.Status.CurrentSchedule).Should(Equal("@every 10s"))

		})
	})

})
