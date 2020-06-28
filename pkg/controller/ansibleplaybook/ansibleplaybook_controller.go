package ansibleplaybook

import (
	"context"

	appv1alpha1 "github.com/ansible-operator/ansible-operator/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_ansibleplaybook")


func AnsibleConfigMap(cr *appv1alpha1.AnsiblePlaybook, namespace string) (*corev1.ConfigMap, controllerutil.MutateFn) {


	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ansible-configmap",
			Namespace: cr.ObjectMeta.Namespace,
		},
	}

	f := func() error {
		if err := controllerutil.SetControllerReference(cr, configMap, scheme); err != nil {
			return err
		}
		addAppLabel(cr.Spec.AppName, &configMap.ObjectMeta)
		configMap.Data = data
		return nil
	}

	return configMap, f
}


func newSecrets(c *CoreV1Client, namespace string) *secrets {
	return &secrets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}


// Create takes the representation of a secret and creates it.  Returns the server's representation of the secret, and an error, if there is any.
func (c *secrets) Create(secret *v1.Secret) (result *v1.Secret, err error) {
	result = &v1.Secret{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("secrets").
		Body(secret).
		Do().
		Into(result)
	return
}


func DefaultSecret(cr *appv1alpha1.AnsiblePlaybook) *corev1.Secret {
	labels := map[string]string{
		"app": cr.Spec.AppName,
	}
	secret := map[string]string{
		"username":   /runner/env/envvars,
		"password":   /runner/env/password,
		"hostname":   /runner/inventory/hosts,
		"ssh_key" :   /runner/env/ssh_key,
		"type"    :   /runner/env/envvars,
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      postgresqlSecretName(cr),
			Namespace: cr.ObjectMeta.Namespace,
			Labels:    labels,
		},
		StringData: secret,
	}
}

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AnsiblePlaybook Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAnsiblePlaybook{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ansibleplaybook-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AnsiblePlaybook
	err = c.Watch(&source.Kind{Type: &appv1alpha1.AnsiblePlaybook{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner AnsiblePlaybook
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.AnsiblePlaybook{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileAnsiblePlaybook implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAnsiblePlaybook{}

// ReconcileAnsiblePlaybook reconciles a AnsiblePlaybook object
type ReconcileAnsiblePlaybook struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}



// func (cr *Controller) getName() string {
// 	return "provisioner"
// }



// func (cr *Controller) handleAnsiblePlaybookAdd(ansibleplaybook *types.AnsiblePlaybook) error {
// 	config, err := getConfigStr(ansibleplaybook)
// 	if err != nil {
// 		return err
// 	}
// 	// Compare applied vs current config, and only run update when there are changes
// 	if config == ansibleplaybook.Status.AppliedConfig {
// 		return nil
// 	}

// 	logrus.Infof("AnsiblePlaybook [%s] is updated; provisioning...", ansibleplaybook.Name)
// 	// Add finalizer and other init fields
// 	if err := cr.initialize(ansibleplaybook, cr.getName()); err != nil {
// 		return fmt.Errorf("error initializing ansibleplaybook %s %v", ansibleplaybook.Name, err)
// 	}

// 	// Provision the ansibleplaybook
// 	_, err = types.AnsiblePlaybookConditionProvisioned.Do(ansibleplaybook, func() (runtime.Object, error) {
// 		// this is the place where cluster provisioning backend logic is being invoked
// 		return ansibleplaybook, provisionAnsiblePlaybook(ansibleplaybook)
// 	})

// 	if err != nil {
// 		return fmt.Errorf("error provisioning ansibleplaybook %s %v", ansibleplaybook.Name, err)
// 	}
// 	// Update ansibleplaybook with applied spec
// 	if err := cr.updateAppliedConfig(ansibleplaybook, config); err != nil {
// 		return fmt.Errorf("error updating ansibleplaybook %s %v", ansibleplaybook.Name, err)
// 	}
// 	logrus.Infof("Successfully provisioned ansible playbook %v", ansibleplaybook.Name)
// 	return nil
// }


// func getConfigStr(ansibleplaybook *types.AnsiblePlaybook) (string, error) {
// 	b, err := ioutil.ReadFile(ansibleplaybook.Spec.ConfigPath)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(b), nil
// }

// func (cr *Controller) initialize(ansibleplaybook *types.AnsiblePlaybook, finalizerKey string) error {
// 	//set finalizers
// 	metadata, err := meta.Accessor(ansibleplaybook)
// 	if err != nil {
// 		return err
// 	}
// 	if containsString(metadata.GetFinalizers(), finalizerKey) {
// 		return nil
// 	}
// 	finalizers := metadata.GetFinalizers()
// 	finalizers = append(finalizers, finalizerKey)
// 	metadata.SetFinalizers(finalizers)
// 	for i := 0; i < 3; i++ {
// 	//	_, err = c.clusterClient.ClusterprovisionerV1alpha1().AnsiblePlaybook().Update(ansibleplaybook)
// 		if err == nil {
// 			return err
// 		}
// 	}
// 	return nil
// }


// Reconcile reads that state of the cluster for a AnsiblePlaybook object and makes changes based on the state read
// and what is in the AnsiblePlaybook.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAnsiblePlaybook) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AnsiblePlaybook")

	// Fetch the AnsiblePlaybook instance
	instance := &appv1alpha1.AnsiblePlaybook{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// // Define a new Pod object
	// pod := newPodForCR(instance)

	// // Set AnsiblePlaybook instance as the owner and controller
	// if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Check if this Pod already exists
	// found := &corev1.Pod{}
	// err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	// 	err = r.client.Create(context.TODO(), pod)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}

	// 	// Pod created successfully - don't requeue
	// 	return reconcile.Result{}, nil
	// } else if err != nil {
	// 	return reconcile.Result{}, err
	// }

    // Check if the job already exists, if not create a new job.
	found := &v1beta1.CronJob{{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil {
	 	if errors.IsNotFound(err) {
			// Define and create a new deployment.
			dep := r.newJobForCR(instance)
			if err = r.client.Create(context.TODO(), dep); err != nil {
				return reconcile.Result{}, err
			}
			reqLogger.Info("Creating a new Job", "CronJob.Namespace", dep.Namespace, "CronJob.Name", dep.Name)
			return reconcile.Result{Requeue: true}, nil
		} else {
			return reconcile.Result{}, err			
		}
	}



    // Ensure the repository is the same as the spec.
	repo := instance.Spec.Repository
	if *found.Spec.Repository != repo {
		reqLogger.Info("Updating repository...", "Repository.Namespace", repo.Namespace, "Repository.Name", repo.Name)
		found.Spec.Repository = &repo
		if err = r.client.Update(context.TODO(), found); err != nil {
			return reconcile.Result{}, err
		}
		
		return reconcile.Result{Requeue: true}, nil
	}



	// Ensure the ansible_playbook is the same as the spec.
	ap := app.Spec.AnsiblePlaybook
	if *found.Spec.ANsiblePLaybook != ap {
		reqLogger.Info("Updating ansible playbook...")
		found.Spec.AnsiblePlaybook = &ap
		if err = r.client.Update(context.TODO(), found); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}


	// Ensure the type is the same as the spec.
	type := app.Spec.Type
	if *found.Spec.Type != type {
		reqLogger.Info("Updating ansible type...")
		found.Spec.Type = &type
		if err = r.client.Update(context.TODO(), found); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}


	// Ensure the URL is the same as the spec.
	url := app.Spec.URL
	if *found.Spec.URL != url {
		reqLogger.Info("Updating ansible URL...")
		found.Spec.URL = &url
		if err = r.client.Update(context.TODO(), found); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}


	// Ensure the PlaybookContent is the same as the spec.
	pc := app.Spec.PlaybookContent
	if *found.Spec.PlaybookContent != pc {
		reqLogger.Info("Updating ansible playbook content...")
		found.Spec.PlaybookContent = &pc
		if err = r.client.Update(context.TODO(), found); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}


	// Ensure the playbook_name is the same as the spec.
	pn := app.Spec.PlaybookName
	if *found.Spec.PlaybookName != pn {
		reqLogger.Info("Updating ansible playbook name...")
		found.Spec.PlaybookName = &pn
		if err = r.client.Update(context.TODO(), found); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	// Update the App status with the pod names.
	// List the pods for this app's job.
	jobList := &v1beta1.CronJob{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(labelsForApp(instance.Name)),
	}
	if err = r.client.List(context.TODO(), jobList, listOpts...); err != nil {
		return reconcile.Result{}, err
	}


    isAnsiblePLaybookMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isAnsiblePlaybookMarkedToBeDeleted {
		if contains(instance.GetFinalizers(), apFinalizer) {
			// Run finalization logic for memcachedFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeAnsiblePLaybook(reqLogger, instance); err != nil {
				return reconcile.Result{}, err
			}

			// Remove memcachedFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(instance, apFinalizer)
			err := r.client.Update(context.TODO(), memcached)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

    func (r *ReconcileMemcached) finalizeMemcached(reqLogger logr.Logger, m *appv1alpha1.AnsiblePlaybook) error {
		// TODO(user): Add the cleanup steps that the operator
		// needs to do before the CR can be deleted. Examples
		// of finalizers include performing backups and deleting
		// resources that are not owned by this CR, like a PVC.
		reqLogger.Info("Successfully finalized ansible playbook")
		return nil
	}
	
	func (r *ReconcileMemcached) addFinalizer(reqLogger logr.Logger, m *appv1alpha1.AnsiblePlaybook) error {
		reqLogger.Info("Adding Finalizer for the Memcached")
		controllerutil.AddFinalizer(m, memcachedFinalizer)
	
		// Update CR
		err := r.client.Update(context.TODO(), m)
		if err != nil {
			reqLogger.Error(err, "Failed to update Memcached with finalizer")
			return err
		}
		return nil
	}
	
	func contains(list []string, s string) bool {
		for _, v := range list {
			if v == s {
				return true
			}
		}
		return false
	}








	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newJobForCR(cr *appv1alpha1.AnsiblePlaybook) *v1beta1.CronJob {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: v1beta1.CronJob{
			Containers: []v1beta1.Container{
				{
					Name:    "ansible-playbook",
					Image:   "quay.io/tsisodia/ansible-runner",
					Command: []string{"sleep", "3600"},
				},
			Env: []v1beta1.EnvVar{
				v1beta1.EnvVar{
				Name:  "RUNNER_PLAYBOOK",
				Value: /runner/env/envvars,
			},
			v1beta1.EnvVar{
                                Name:  "PROJECT_TYPE",
                                Value: /runner/env/envvars,
                        },
			v1beta1.EnvVar{
                                Name:  "USERNAME",
                                Value: /runner/env/envvars,
                        },

                        },

		        VolumeMounts: []v1beta1.VolumeMount{
			v1beta1.VolumeMount{Name: "configmap-volume", MountPath: "/runner/env/extravars"},
			v1beta1.VolumeMount{Name: "secrets-volume", MountPath: "/runner/env/passwords"},
		   },



		},
	}
}
