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

var _ = Describe("Linkedsecret controller Azure", func() {

	const (
		TIMEOUT  = time.Second * 60
		DURATION = time.Second * 10
		INTERVAL = time.Millisecond * 250
	)

	var (
		azurePlain LinkedSecretTest
		azureJSON  LinkedSecretTest
	)

	BeforeEach(func() {
		azureJSON = LinkedSecretTest{
			name:      "azure-example1",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "Azure",
				ProviderSecretFormat: "JSON",
				ProviderOptions:      map[string]string{"secret": "opaque-secret-json", "keyvault": "linkedsecret"},
				SecretName:           "mysecret-azure-example1",
				Schedule:             "@every 10s",
				Suspended:            false,
			},
		}

		azurePlain = LinkedSecretTest{
			name:      "azure-example2",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "Azure",
				ProviderSecretFormat: "PLAIN",
				ProviderOptions:      map[string]string{"secret": "opaque-secret-plain", "keyvault": "linkedsecret"},
				SecretName:           "mysecret-azure-example2",
				Schedule:             "@every 10s",
				Suspended:            false,
			},
		}
	})

	Context("When creating Linkedsecret azure-example1", func() {

		It("Should create Linkedsecret azure-example1", func() {

			By("Creating azure-example1")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: azureJSON.name, Namespace: azureJSON.namespace},
				Spec:       azureJSON.spec,
			}
			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: azureJSON.namespace, Name: azureJSON.name}
			azureExample1 := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, azureExample1)
				if err != nil {
					return false
				}
				if azureExample1.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check spec
			Expect(azureExample1.Spec.Provider).Should(Equal("Azure"))
			Expect(azureExample1.Spec.ProviderSecretFormat).Should(Equal("JSON"))
			Expect(azureExample1.Spec.ProviderOptions["secret"]).Should(Equal("opaque-secret-json"))
			Expect(azureExample1.Spec.ProviderOptions["keyvault"]).Should(Equal("linkedsecret"))
			Expect(azureExample1.Spec.SecretName).Should(Equal("mysecret-azure-example1"))
			Expect(azureExample1.Spec.Suspended).Should(Equal(false))
			Expect(azureExample1.Spec.Schedule).Should(Equal("@every 1s"))

			// Check status
			Expect(azureExample1.Status.CurrentSecret).Should(Equal("mysecret-azure-example1"))
			Expect(azureExample1.Status.CronJobID).Should(Equal(cron.EntryID(azureExample1.Status.CronJobID)))
			Expect(azureExample1.Status.CronJobStatus).Should(Equal("Scheduled"))

		})
	})

	Context("When creating Linkedsecret azure-example2", func() {

		It("Should create Linkedsecret azure-example2", func() {

			By("Creating azure-example1")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: azurePlain.name, Namespace: azurePlain.namespace},
				Spec:       azurePlain.spec,
			}
			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: azurePlain.namespace, Name: azurePlain.name}
			azureExample2 := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, azureExample2)
				if err != nil {
					return false
				}
				if azureExample2.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check expected spec
			Expect(azureExample2.Spec.Provider).Should(Equal("Azure"))
			Expect(azureExample2.Spec.ProviderSecretFormat).Should(Equal("PLAIN"))
			Expect(azureExample2.Spec.ProviderOptions["secret"]).Should(Equal("opaque-secret-plain"))
			Expect(azureExample2.Spec.ProviderOptions["keyvault"]).Should(Equal("linkedsecret"))
			Expect(azureExample2.Spec.SecretName).Should(Equal("mysecret-azure-example2"))
			Expect(azureExample2.Spec.Suspended).Should(Equal(false))
			Expect(azureExample2.Spec.Schedule).Should(Equal("@every 1s"))

			// Check expected status
			Expect(azureExample2.Status.CurrentSecret).Should(Equal("mysecret-azure-example2"))
			Expect(azureExample2.Status.CronJobID).Should(Equal(cron.EntryID(azureExample2.Status.CronJobID)))
			Expect(azureExample2.Status.CronJobStatus).Should(Equal("Scheduled"))
			Expect(azureExample2.Status.CurrentSchedule).Should(Equal("@every 1s"))

		})
	})

})
