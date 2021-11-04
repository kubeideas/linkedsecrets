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

var _ = Describe("Linkedsecret controller IBM", func() {

	const (
		TIMEOUT                 = time.Second * 60
		INTERVAL                = time.Millisecond * 100
		SECRETMANAGERINSTANCEID = "8d2350b3-7ce3-4852-8b4b-a5cc6fd5f146"
		JSONSECRETID            = "53a1db89-ce4e-0a39-a699-4a91ca9920a5"
		PLANTEXTSECRETID        = "5c5c4c05-31e8-7c5a-c7bd-e4d7e42d6547"
		REGION                  = "us-east"
	)

	var ibmPlain LinkedSecretTest
	var ibmJSON LinkedSecretTest

	BeforeEach(func() {
		ibmJSON = LinkedSecretTest{
			name:      "ibm-example1",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:           "IBM",
				ProviderDataFormat: "JSON",
				ProviderOptions:    map[string]string{"secretManagerInstanceId": SECRETMANAGERINSTANCEID, "secretId": JSONSECRETID, "region": REGION},
				SecretName:         "mysecret-ibm-example1",
				Schedule:           "@every 1s",
				Suspended:          false,
			},
		}

		ibmPlain = LinkedSecretTest{
			name:      "ibm-example2",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:           "IBM",
				ProviderDataFormat: "PLAIN",
				ProviderOptions:    map[string]string{"secretManagerInstanceId": SECRETMANAGERINSTANCEID, "secretId": PLANTEXTSECRETID, "region": REGION},
				SecretName:         "mysecret-ibm-example2",
				Schedule:           "@every 1s",
				Suspended:          false,
			},
		}
	})

	Describe("Creating IBM JSON Linkedsecret", func() {
		Context("Creating new Linkedsecret sinchronizing data with IBM", func() {
			It("Should be IBM Linkedsecret", func() {
				By("Creating new IBM linkedsecret")
				ctx := context.Background()
				linkedSecret := &securityv1.LinkedSecret{
					TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
					ObjectMeta: v1.ObjectMeta{Name: ibmJSON.name, Namespace: ibmJSON.namespace},
					Spec:       ibmJSON.spec,
				}

				// Create new LinkeSecret
				Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

				linkedSecretLookupKey := types.NamespacedName{Namespace: ibmJSON.namespace, Name: ibmJSON.name}
				createdLinkedSecret := &securityv1.LinkedSecret{}

				// Get linkedSecret
				Eventually(func() bool {
					err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
					return err == nil
				}, TIMEOUT, INTERVAL).Should(BeTrue())

				// Check expected spec
				Expect(createdLinkedSecret.Spec.Provider).Should(Equal("IBM"))
				Expect(createdLinkedSecret.Spec.ProviderDataFormat).Should(Equal("JSON"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["secretManagerInstanceId"]).Should(Equal(SECRETMANAGERINSTANCEID))
				Expect(createdLinkedSecret.Spec.ProviderOptions["secretId"]).Should(Equal(JSONSECRETID))
				Expect(createdLinkedSecret.Spec.ProviderOptions["region"]).Should(Equal(REGION))
				Expect(createdLinkedSecret.Spec.SecretName).Should(Equal("mysecret-ibm-example1"))
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
				Expect(createdLinkedSecret.Status.CreatedSecret).Should(Equal("mysecret-ibm-example1"))
				Expect(createdLinkedSecret.Status.CronJobID).Should(Equal(cron.EntryID(1)))
				Expect(createdLinkedSecret.Status.CronJobStatus).Should(Equal("Scheduled"))
				Expect(createdLinkedSecret.Status.CurrentProvider).Should(Equal("IBM"))
				Expect(createdLinkedSecret.Status.CurrentSchedule).Should(Equal("@every 1s"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["secretManagerInstanceId"]).Should(Equal(SECRETMANAGERINSTANCEID))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["secretId"]).Should(Equal(JSONSECRETID))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["region"]).Should(Equal(REGION))

			})
		})
	})

	Describe("Creating IBM Plain Linkedsecret", func() {
		Context("Creating new Linkedsecret sinchronizing data with IBM", func() {
			It("Should be IBM Linkedsecret", func() {
				By("Creating new IBM linkedsecret")
				ctx := context.Background()
				linkedSecret := &securityv1.LinkedSecret{
					TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
					ObjectMeta: v1.ObjectMeta{Name: ibmPlain.name, Namespace: ibmPlain.namespace},
					Spec:       ibmPlain.spec,
				}
				// Create new LinkeSecret
				Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

				linkedSecretLookupKey := types.NamespacedName{Namespace: ibmPlain.namespace, Name: ibmPlain.name}
				createdLinkedSecret := &securityv1.LinkedSecret{}

				// Get linkedSecret
				Eventually(func() bool {
					err := k8sClient.Get(ctx, linkedSecretLookupKey, createdLinkedSecret)
					return err == nil
				}, TIMEOUT, INTERVAL).Should(BeTrue())

				// Check expected spec
				Expect(createdLinkedSecret.Spec.Provider).Should(Equal("IBM"))
				Expect(createdLinkedSecret.Spec.ProviderDataFormat).Should(Equal("PLAIN"))
				Expect(createdLinkedSecret.Spec.ProviderOptions["secretManagerInstanceId"]).Should(Equal(SECRETMANAGERINSTANCEID))
				Expect(createdLinkedSecret.Spec.ProviderOptions["secretId"]).Should(Equal(PLANTEXTSECRETID))
				Expect(createdLinkedSecret.Spec.ProviderOptions["region"]).Should(Equal(REGION))
				Expect(createdLinkedSecret.Spec.SecretName).Should(Equal("mysecret-ibm-example2"))
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
				Expect(createdLinkedSecret.Status.CreatedSecret).Should(Equal("mysecret-ibm-example2"))
				Expect(createdLinkedSecret.Status.CronJobID).Should(Equal(cron.EntryID(1)))
				Expect(createdLinkedSecret.Status.CronJobStatus).Should(Equal("Scheduled"))
				Expect(createdLinkedSecret.Status.CurrentProvider).Should(Equal("IBM"))
				Expect(createdLinkedSecret.Status.CurrentSchedule).Should(Equal("@every 1s"))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["secretManagerInstanceId"]).Should(Equal(SECRETMANAGERINSTANCEID))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["secretId"]).Should(Equal(PLANTEXTSECRETID))
				Expect(createdLinkedSecret.Status.CurrentProviderOptions["region"]).Should(Equal(REGION))

			})
		})
	})

})
