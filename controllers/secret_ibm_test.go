package controllers

import (
	"context"
	securityv1 "kubeideas/linkedsecrets/api/v1"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robfig/cron/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Linkedsecret controller IBM", func() {

	const (
		TIMEOUT  = time.Second * 60
		DURATION = time.Second * 10
		INTERVAL = time.Millisecond * 250
		//secretManagerInstance = "f0a9ef5b-de69-484b-ab84-4181390eec1e"
		//ibmJSONSecretUUID            = "15dd61cd-1f3b-02fc-ecd3-14b25c39302d"
		//ibmPlainSecretUUID       = "baa3a63f-f5ab-eea0-02d9-89556fc54cc5"
		REGION = "us-east"
	)

	var (
		secretManagerInstanceId = os.Getenv("SECRET_MANAGER_UUID")
		ibmJSONSecretUUID       = os.Getenv("SECRET_JSON_UUID")
		ibmPlainSecretUUID      = os.Getenv("SECRET_PLAIN_UUID")
		ibmPlain                LinkedSecretTest
		ibmJSON                 LinkedSecretTest
		ibmInvalidUUID          LinkedSecretTest
	)

	BeforeEach(func() {
		ibmJSON = LinkedSecretTest{
			name:      "ibm-example1",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "IBM",
				ProviderSecretFormat: "JSON",
				ProviderOptions:      map[string]string{"secretManagerInstanceId": secretManagerInstanceId, "secretId": ibmJSONSecretUUID, "region": REGION},
				SecretName:           "mysecret-ibm-example1",
				Schedule:             "@every 10s",
				Suspended:            false,
			},
		}

		ibmPlain = LinkedSecretTest{
			name:      "ibm-example2",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "IBM",
				ProviderSecretFormat: "PLAIN",
				ProviderOptions:      map[string]string{"secretManagerInstanceId": secretManagerInstanceId, "secretId": ibmPlainSecretUUID, "region": REGION},
				SecretName:           "mysecret-ibm-example2",
				Schedule:             "@every 10s",
				Suspended:            false,
			},
		}

		ibmInvalidUUID = LinkedSecretTest{
			name:      "ibm-example3",
			namespace: "default",
			spec: securityv1.LinkedSecretSpec{
				Provider:             "IBM",
				ProviderSecretFormat: "PLAIN",
				ProviderOptions:      map[string]string{"secretManagerInstanceId": "invalid-uuid-uuid-uuid-invalid34uuid", "secretId": ibmPlainSecretUUID, "region": REGION},
				SecretName:           "mysecret-ibm-example3",
				Schedule:             "@every 10s",
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
			Expect(ibmExample1.Spec.ProviderOptions["secretManagerInstanceId"]).Should(Equal(secretManagerInstanceId))
			Expect(ibmExample1.Spec.ProviderOptions["secretId"]).Should(Equal(ibmJSONSecretUUID))
			Expect(ibmExample1.Spec.ProviderOptions["region"]).Should(Equal(REGION))
			Expect(ibmExample1.Spec.SecretName).Should(Equal("mysecret-ibm-example1"))
			Expect(ibmExample1.Spec.Suspended).Should(Equal(false))
			Expect(ibmExample1.Spec.Schedule).Should(Equal("@every 10s"))

			// Check status
			Expect(ibmExample1.Status.CurrentSecret).Should(Equal("mysecret-ibm-example1"))
			Expect(ibmExample1.Status.CronJobID).Should(Equal(cron.EntryID(ibmExample1.Status.CronJobID)))
			Expect(ibmExample1.Status.CronJobStatus).Should(Equal("Scheduled"))
			Expect(ibmExample1.Status.CurrentSchedule).Should(Equal("@every 10s"))

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
			Expect(ibmExample2.Spec.ProviderOptions["secretManagerInstanceId"]).Should(Equal(secretManagerInstanceId))
			Expect(ibmExample2.Spec.ProviderOptions["secretId"]).Should(Equal(ibmPlainSecretUUID))
			Expect(ibmExample2.Spec.ProviderOptions["region"]).Should(Equal(REGION))
			Expect(ibmExample2.Spec.SecretName).Should(Equal("mysecret-ibm-example2"))
			Expect(ibmExample2.Spec.Suspended).Should(Equal(false))
			Expect(ibmExample2.Spec.Schedule).Should(Equal("@every 10s"))

			// Check status
			Expect(ibmExample2.Status.CurrentSecret).Should(Equal("mysecret-ibm-example2"))
			Expect(ibmExample2.Status.CronJobID).Should(Equal(cron.EntryID(ibmExample2.Status.CronJobID)))
			Expect(ibmExample2.Status.CronJobStatus).Should(Equal("Scheduled"))
			Expect(ibmExample2.Status.CurrentSchedule).Should(Equal("@every 10s"))

		})
	})

	Context("When creating Linkedsecret ibm-example3", func() {

		It("Should create Linkedsecret ibm-example3", func() {

			By("Creating Linkedsecret ibm-example3")
			ctx := context.Background()
			linkedSecret := &securityv1.LinkedSecret{
				TypeMeta:   v1.TypeMeta{Kind: "LinkedSecret", APIVersion: "linkedsecrets/api/v1"},
				ObjectMeta: v1.ObjectMeta{Name: ibmInvalidUUID.name, Namespace: ibmInvalidUUID.namespace},
				Spec:       ibmInvalidUUID.spec,
			}
			// Create new LinkeSecret
			Expect(k8sClient.Create(ctx, linkedSecret)).Should(Succeed())

			linkedSecretLookupKey := types.NamespacedName{Namespace: ibmInvalidUUID.namespace, Name: ibmInvalidUUID.name}
			ibmExample3 := &securityv1.LinkedSecret{}

			// Get linkedSecret
			Eventually(func() bool {
				err := k8sClient.Get(ctx, linkedSecretLookupKey, ibmExample3)
				if err != nil {
					return false
				}
				if ibmExample3.Status.CurrentSecret == "" {
					return false
				}
				return true
			}, TIMEOUT, INTERVAL).Should(BeTrue())

			// Check spec
			Expect(ibmExample3.Spec.Provider).Should(Equal("IBM"))
			Expect(ibmExample3.Spec.ProviderSecretFormat).Should(Equal("PLAIN"))
			Expect(ibmExample3.Spec.ProviderOptions["secretId"]).Should(Equal(ibmPlainSecretUUID))
			Expect(ibmExample3.Spec.ProviderOptions["region"]).Should(Equal(REGION))
			Expect(ibmExample3.Spec.SecretName).Should(Equal("mysecret-ibm-example3"))
			Expect(ibmExample3.Spec.Suspended).Should(Equal(false))
			Expect(ibmExample3.Spec.Schedule).Should(Equal("@every 10s"))

			// Check status
			Expect(ibmExample3.Status.CurrentSecret).Should(Equal("mysecret-ibm-example3"))
			Expect(ibmExample3.Status.CronJobStatus).Should(Equal("NotScheduled"))
			Expect(ibmExample3.Status.CurrentSchedule).Should(Equal("@every 10s"))

		})
	})

})
