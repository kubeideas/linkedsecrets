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

var _ = Describe("Linkedsecret controller IBM", func() {

	const (
		TIMEOUT                 = time.Second * 60
		DURATION                = time.Second * 10
		INTERVAL                = time.Millisecond * 250
		SECRETMANAGERINSTANCEID = "8d2350b3-7ce3-4852-8b4b-a5cc6fd5f146"
		JSONSECRETID            = "53a1db89-ce4e-0a39-a699-4a91ca9920a5"
		PLANTEXTSECRETID        = "5c5c4c05-31e8-7c5a-c7bd-e4d7e42d6547"
		REGION                  = "us-east"
	)

	var (
		ibmPlain LinkedSecretTest
		ibmJSON  LinkedSecretTest
	)

	BeforeEach(func() {
		ibmJSON = LinkedSecretTest{
			name:      "ibm-example1",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "IBM",
				ProviderSecretFormat: "JSON",
				ProviderOptions:      map[string]string{"secretManagerInstanceId": SECRETMANAGERINSTANCEID, "secretId": JSONSECRETID, "region": REGION},
				SecretName:           "mysecret-ibm-example1",
				Schedule:             "@every 1s",
				Suspended:            false,
			},
		}

		ibmPlain = LinkedSecretTest{
			name:      "ibm-example2",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "IBM",
				ProviderSecretFormat: "PLAIN",
				ProviderOptions:      map[string]string{"secretManagerInstanceId": SECRETMANAGERINSTANCEID, "secretId": PLANTEXTSECRETID, "region": REGION},
				SecretName:           "mysecret-ibm-example2",
				Schedule:             "@every 1s",
				Suspended:            false,
			},
		}
	})

	Context("When creating Linkedsecret ibm-example1", func() {

		It("Should create Linkedsecret ibm-example1", func() {

			By("Creating Linkedsecret ibm-example1")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: ibmJSON.name, Namespace: ibmJSON.namespace},
				Spec:       ibmJSON.spec,
			}

			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: ibmJSON.namespace, Name: ibmJSON.name}
			ibmExample1 := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, ibmExample1)
				if err != nil {
					return false
				}
				if ibmExample1.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check spec
			Expect(ibmExample1.Spec.Provider).Should(Equal("IBM"))
			Expect(ibmExample1.Spec.ProviderSecretFormat).Should(Equal("JSON"))
			Expect(ibmExample1.Spec.ProviderOptions["secretManagerInstanceId"]).Should(Equal(SECRETMANAGERINSTANCEID))
			Expect(ibmExample1.Spec.ProviderOptions["secretId"]).Should(Equal(JSONSECRETID))
			Expect(ibmExample1.Spec.ProviderOptions["region"]).Should(Equal(REGION))
			Expect(ibmExample1.Spec.SecretName).Should(Equal("mysecret-ibm-example1"))
			Expect(ibmExample1.Spec.Suspended).Should(Equal(false))
			Expect(ibmExample1.Spec.Schedule).Should(Equal("@every 1s"))

			// Check status
			Expect(ibmExample1.Status.CurrentSecret).Should(Equal("mysecret-ibm-example1"))
			Expect(ibmExample1.Status.CronJobID).Should(Equal(cron.EntryID(1)))
			Expect(ibmExample1.Status.CronJobStatus).Should(Equal("Scheduled"))
			Expect(ibmExample1.Status.CurrentSchedule).Should(Equal("@every 1s"))

		})
	})

	Context("When creating Linkedsecret ibm-example2", func() {

		It("Should create Linkedsecret ibm-example2", func() {

			By("Creating Linkedsecret ibm-example2")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: ibmPlain.name, Namespace: ibmPlain.namespace},
				Spec:       ibmPlain.spec,
			}
			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: ibmPlain.namespace, Name: ibmPlain.name}
			ibmExample2 := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, ibmExample2)
				if err != nil {
					return false
				}
				if ibmExample2.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check spec
			Expect(ibmExample2.Spec.Provider).Should(Equal("IBM"))
			Expect(ibmExample2.Spec.ProviderSecretFormat).Should(Equal("PLAIN"))
			Expect(ibmExample2.Spec.ProviderOptions["secretManagerInstanceId"]).Should(Equal(SECRETMANAGERINSTANCEID))
			Expect(ibmExample2.Spec.ProviderOptions["secretId"]).Should(Equal(PLANTEXTSECRETID))
			Expect(ibmExample2.Spec.ProviderOptions["region"]).Should(Equal(REGION))
			Expect(ibmExample2.Spec.SecretName).Should(Equal("mysecret-ibm-example2"))
			Expect(ibmExample2.Spec.Suspended).Should(Equal(false))
			Expect(ibmExample2.Spec.Schedule).Should(Equal("@every 1s"))

			// Check status
			Expect(ibmExample2.Status.CurrentSecret).Should(Equal("mysecret-ibm-example2"))
			Expect(ibmExample2.Status.CronJobID).Should(Equal(cron.EntryID(1)))
			Expect(ibmExample2.Status.CronJobStatus).Should(Equal("Scheduled"))
			Expect(ibmExample2.Status.CurrentSchedule).Should(Equal("@every 1s"))

		})
	})

})
