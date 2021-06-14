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

var _ = Describe("Linkedsecret controller AWS", func() {

	const (
		TIMEOUT  = time.Second * 60
		INTERVAL = time.Millisecond * 100
	)

	var awsPlain LinkedSecretTest
	var awsJSON LinkedSecretTest

	BeforeEach(func() {
		awsJSON = LinkedSecretTest{
			name:      "aws-example1",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:           "AWS",
				ProviderDataFormat: "JSON",
				ProviderOptions:    map[string]string{"secret": "secret-json-tst", "region": "us-east-1", "version": "AWSCURRENT"},
				SecretName:         "mysecret-aws-example1",
				Schedule:           "@every 1s",
				Suspended:          false,
			},
		}

		awsPlain = LinkedSecretTest{
			name:      "aws-example2",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:           "AWS",
				ProviderDataFormat: "PLAIN",
				ProviderOptions:    map[string]string{"secret": "secret-plain-tst", "region": "us-east-1", "version": "AWSCURRENT"},
				SecretName:         "mysecret-aws-example2",
				Schedule:           "@every 1s",
				Suspended:          false,
			},
		}
	})

	Describe("Creating AWS JSON Linkedsecret", func() {
		Context("Creating new Linkedsecret sinchronizing data with AWS", func() {
			It("Should be AWS Linkedsecret", func() {
				By("Creating new AWS linkedsecret")
				ctx := context.Background()
				linkedSecret := &securityv1.LinkedSecret{
					TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
					ObjectMeta: v1.ObjectMeta{Name: awsJSON.name, Namespace: awsJSON.namespace},
					Spec:       awsJSON.spec,
				}
				// Create new LinkeSecret
				Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

				linkedSecretLookupKey := types.NamespacedName{Namespace: awsJSON.namespace, Name: awsJSON.name}
				createdLinkedSecret := &securityv1.LinkedSecret{}

				// Get linkedSecret
				Eventually(func() bool {
					err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
					return err == nil
				}, TIMEOUT, INTERVAL).Should(BeTrue())

				// Check expected spec
				Expect(createdLinkedSecret.Spec.Provider).Should(Equal("AWS"))
				Expect(createdLinkedSecret.Spec.ProviderDataFormat).Should(Equal("JSON"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["secret"]).Should(Equal("secret-json-tst"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["region"]).Should(Equal("us-east-1"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["version"]).Should(Equal("AWSCURRENT"))
				Expect(createdLinkedSecret.Spec.SecretName).Should(Equal("mysecret-aws-example1"))
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
				Expect(createdLinkedSecret.Status.CreatedSecret).Should(Equal("mysecret-aws-example1"))
				Expect(createdLinkedSecret.Status.CronJobID).Should(Equal(cron.EntryID(1)))
				Expect(createdLinkedSecret.Status.CronJobStatus).Should(Equal("Scheduled"))
				Expect(createdLinkedSecret.Status.CurrentProvider).Should(Equal("AWS"))
				Expect(createdLinkedSecret.Status.CurrentSchedule).Should(Equal("@every 1s"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["secret"]).Should(Equal("secret-json-tst"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["region"]).Should(Equal("us-east-1"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["version"]).Should(Equal("AWSCURRENT"))

			})
		})
	})

	Describe("Creating AWS Plain Linkedsecret", func() {
		Context("Creating new Linkedsecret sinchronizing data with AWS", func() {
			It("Should be AWS Linkedsecret", func() {
				By("Creating new AWS linkedsecret")
				ctx := context.Background()
				linkedSecret := &securityv1.LinkedSecret{
					TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
					ObjectMeta: v1.ObjectMeta{Name: awsPlain.name, Namespace: awsPlain.namespace},
					Spec:       awsPlain.spec,
				}
				// Create new LinkeSecret
				Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

				linkedSecretLookupKey := types.NamespacedName{Namespace: awsPlain.namespace, Name: awsPlain.name}
				createdLinkedSecret := &securityv1.LinkedSecret{}

				// Get linkedSecret
				Eventually(func() bool {
					err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
					return err == nil
				}, TIMEOUT, INTERVAL).Should(BeTrue())

				// Check expected spec
				Expect(createdLinkedSecret.Spec.Provider).Should(Equal("AWS"))
				Expect(createdLinkedSecret.Spec.ProviderDataFormat).Should(Equal("PLAIN"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["secret"]).Should(Equal("secret-plain-tst"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["region"]).Should(Equal("us-east-1"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["version"]).Should(Equal("AWSCURRENT"))
				Expect(createdLinkedSecret.Spec.SecretName).Should(Equal("mysecret-aws-example2"))
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
				Expect(createdLinkedSecret.Status.CreatedSecret).Should(Equal("mysecret-aws-example2"))
				Expect(createdLinkedSecret.Status.CronJobID).Should(Equal(cron.EntryID(1)))
				Expect(createdLinkedSecret.Status.CronJobStatus).Should(Equal("Scheduled"))
				Expect(createdLinkedSecret.Status.CurrentProvider).Should(Equal("AWS"))
				Expect(createdLinkedSecret.Status.CurrentSchedule).Should(Equal("@every 1s"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["secret"]).Should(Equal("secret-plain-tst"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["region"]).Should(Equal("us-east-1"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["version"]).Should(Equal("AWSCURRENT"))

			})
		})
	})

})
