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

var _ = Describe("Linkedsecret controller Azure", func() {

	const (
		TIMEOUT  = time.Second * 60
		INTERVAL = time.Millisecond * 100
	)

	var azurePlain LinkedSecretTest
	var azureJSON LinkedSecretTest

	BeforeEach(func() {
		azureJSON = LinkedSecretTest{
			name:      "azure-example1",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:           "Azure",
				ProviderDataFormat: "JSON",
				ProviderOptions:    map[string]string{"secret": "opaque-secret-json", "keyvault": "linkedsecret"},
				SecretName:         "mysecret-azure-example1",
				Schedule:           "@every 1s",
				Suspended:          false,
			},
		}

		azurePlain = LinkedSecretTest{
			name:      "azure-example2",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:           "Azure",
				ProviderDataFormat: "PLAIN",
				ProviderOptions:    map[string]string{"secret": "opaque-secret-plain", "keyvault": "linkedsecret"},
				SecretName:         "mysecret-azure-example2",
				Schedule:           "@every 1s",
				Suspended:          false,
			},
		}
	})

	Describe("Creating Azure JSON Linkedsecret", func() {
		Context("Creating new Linkedsecret sinchronizing data with Azure", func() {
			It("Should be Azure Linkedsecret", func() {
				By("Creating new Azure linkedsecret")
				ctx := context.Background()
				linkedSecret := &securityv1.LinkedSecret{
					TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
					ObjectMeta: v1.ObjectMeta{Name: azureJSON.name, Namespace: azureJSON.namespace},
					Spec:       azureJSON.spec,
				}
				// Create new LinkeSecret
				Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

				linkedSecretLookupKey := types.NamespacedName{Namespace: azureJSON.namespace, Name: azureJSON.name}
				createdLinkedSecret := &securityv1.LinkedSecret{}

				// Get linkedSecret
				Eventually(func() bool {
					err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
					return err == nil
				}, TIMEOUT, INTERVAL).Should(BeTrue())

				// Check expected spec
				Expect(createdLinkedSecret.Spec.Provider).Should(Equal("Azure"))
				Expect(createdLinkedSecret.Spec.ProviderDataFormat).Should(Equal("JSON"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["secret"]).Should(Equal("opaque-secret-json"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["keyvault"]).Should(Equal("linkedsecret"))
				Expect(createdLinkedSecret.Spec.SecretName).Should(Equal("mysecret-azure-example1"))
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
				Expect(createdLinkedSecret.Status.CreatedSecret).Should(Equal("mysecret-azure-example1"))
				Expect(createdLinkedSecret.Status.CronJobID).Should(Equal(cron.EntryID(1)))
				Expect(createdLinkedSecret.Status.CronJobStatus).Should(Equal("Scheduled"))
				Expect(createdLinkedSecret.Status.CurrentProvider).Should(Equal("Azure"))
				Expect(createdLinkedSecret.Status.CurrentSchedule).Should(Equal("@every 1s"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["secret"]).Should(Equal("opaque-secret-json"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["keyvault"]).Should(Equal("linkedsecret"))

			})
		})
	})

	Describe("Creating Azure Plain Linkedsecret", func() {
		Context("Creating new Linkedsecret sinchronizing data with Azure", func() {
			It("Should be Azure Linkedsecret", func() {
				By("Creating new Azure linkedsecret")
				ctx := context.Background()
				linkedSecret := &securityv1.LinkedSecret{
					TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
					ObjectMeta: v1.ObjectMeta{Name: azurePlain.name, Namespace: azurePlain.namespace},
					Spec:       azurePlain.spec,
				}
				// Create new LinkeSecret
				Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

				linkedSecretLookupKey := types.NamespacedName{Namespace: azurePlain.namespace, Name: azurePlain.name}
				createdLinkedSecret := &securityv1.LinkedSecret{}

				// Get linkedSecret
				Eventually(func() bool {
					err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
					return err == nil
				}, TIMEOUT, INTERVAL).Should(BeTrue())

				// Check expected spec
				Expect(createdLinkedSecret.Spec.Provider).Should(Equal("Azure"))
				Expect(createdLinkedSecret.Spec.ProviderDataFormat).Should(Equal("PLAIN"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["secret"]).Should(Equal("opaque-secret-plain"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["keyvault"]).Should(Equal("linkedsecret"))
				Expect(createdLinkedSecret.Spec.SecretName).Should(Equal("mysecret-azure-example2"))
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
				Expect(createdLinkedSecret.Status.CreatedSecret).Should(Equal("mysecret-azure-example2"))
				Expect(createdLinkedSecret.Status.CronJobID).Should(Equal(cron.EntryID(1)))
				Expect(createdLinkedSecret.Status.CronJobStatus).Should(Equal("Scheduled"))
				Expect(createdLinkedSecret.Status.CurrentProvider).Should(Equal("Azure"))
				Expect(createdLinkedSecret.Status.CurrentSchedule).Should(Equal("@every 1s"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["secret"]).Should(Equal("opaque-secret-plain"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["keyvault"]).Should(Equal("linkedsecret"))

			})
		})
	})

})
