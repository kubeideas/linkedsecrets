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

package main

import (
	"flag"
	"os"

	// pprof http tool
	"net/http"
	_ "net/http/pprof"

	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	securityv1 "linkedsecrets/api/v1"
	"linkedsecrets/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = securityv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {

	var metricsAddr string
	var enableLeaderElection bool
	var enablePProf bool

	//parse flags
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	//pprof http api option
	flag.BoolVar(&enablePProf, "enable-pprof", false, "Enable pprof http endpoint.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	//pprof http endpoint
	if enablePProf {
		go func() {
			log := ctrl.Log.WithName("controller-runtime").WithName("pprof")
			log.Info("pprof server is starting to listen", "addr", ":6060")
			log.Error(http.ListenAndServe("localhost:6060", nil), "Error")
		}()
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "dde0d8ed.kubeideas.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.LinkedSecretReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("LinkedSecret"),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("linkedsecret-controller"),
		Cronjob:  make(map[types.UID]*cron.Cron),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "LinkedSecret")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}

}
