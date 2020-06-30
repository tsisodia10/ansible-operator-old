package ansibleplaybookrun

import (
	"context"

	ansiblev1alpha1 "github.com/ansible-operator/ansible-operator/pkg/apis/ansible/v1alpha1"
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

var log = logf.Log.WithName("controller_ansibleplaybookrun")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AnsiblePlaybookRun Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAnsiblePlaybookRun{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ansibleplaybookrun-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AnsiblePlaybookRun
	err = c.Watch(&source.Kind{Type: &ansiblev1alpha1.AnsiblePlaybookRun{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Jobs and requeue the owner AnsiblePlaybookRun
	err = c.Watch(&source.Kind{Type: &v1beta1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ansiblev1alpha1.AnsiblePlaybookRun{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileAnsiblePlaybookRun implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAnsiblePlaybookRun{}

// ReconcileAnsiblePlaybookRun reconciles a AnsiblePlaybookRun object
type ReconcileAnsiblePlaybookRun struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AnsiblePlaybookRun object and makes changes based on the state read
// and what is in the AnsiblePlaybookRun.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAnsiblePlaybookRun) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AnsiblePlaybookRun")

	// Fetch the AnsiblePlaybookRun instance
	instance := &ansiblev1alpha1.AnsiblePlaybookRun{}
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

	// Define a new Job object
	job := newJobForCR(instance)

	// Set AnsiblePlaybookRun instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

// newJobForCR returns an ansible-runner pod with the same name/namespace as the cr
func newJobForCR(cr *ansiblev1alpha1.AnsiblePlaybookRun) *v1beta1.Job {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &v1beta1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: v1beta1.JobSpec{
			Containers: []v1beta1.Container{
				{
					Name:    "ansible-runner",
					Image:   "quay.io/tsisodia/ansible-runner",
					Command: []string{"sleep", "3600"},
				},
				Env: []v1beta1.EnvVar{},
				VolumeMounts: []v1beta1.VolumeMount{
				  v1beta1.VolumeMount{
					  Name: "extravars-volume"
					  MountPath: "/runner/env/extravars"
				  },
				  v1beta1.VolumeMount{
					  Name: "password-volume"
					  MountPath: "/runner/env/password"
				  },
				  v1beta1.VolumeMount{
					  Name: "sshkey-volume"
					  MountPath: "/runner/env/ssh_key"
				  },
				  v1beta1.VolumeMount{
					  Name: "inventory-volume"
					  MountPath: "/runner/inventory/hosts"
				  },
				  v1beta1.VolumeMount{
					  Name: "projectvars-volume"
					  MountPath: "/runner/project/roles/testrole/vars"
				  },
				  v1beta1.VolumeMount{
					  Name: "projectmeta-volume"
					  MountPath: "/runner/project/roles/testrole/meta"
				  },
				},

			},
		},
		Volumes: []v1beta1.Volume{
			v1beta1.Volume{
				Name: "extravars-volume"
				VolumeSource: v1beta1.VolumeSource{
					HostPath: &v1beta1.HostPathVolumeSource{
						Path: "/data"
					},
				},
			},
			v1beta1.Volume{
                                Name: "password-volume"
                                VolumeSource: v1beta1.VolumeSource{
                                        HostPath: &v1beta1.HostPathVolumeSource{
                                                Path: "/data"
                                        },
                                },
                        },
			v1beta1.Volume{
                                Name: "sshkey-volume"
                                VolumeSource: v1beta1.VolumeSource{
                                        HostPath: &v1beta1.HostPathVolumeSource{
                                                Path: "/data"
                                        },
                                },
                        },
			v1beta1.Volume{
                                Name: "inventory-volume"
                                VolumeSource: v1beta1.VolumeSource{
                                        HostPath: &v1beta1.HostPathVolumeSource{
                                                Path: "/data"
                                        },
                                },
                        },
			v1beta1.Volume{
                                Name: "projectvars-volume"
                                VolumeSource: v1beta1.VolumeSource{
                                        HostPath: &v1beta1.HostPathVolumeSource{
                                                Path: "/data"
                                        },
                                },
                        },
			v1beta1.Volume{
                                Name: "projectmeta-volume"
                                VolumeSource: v1beta1.VolumeSource{
                                        HostPath: &v1beta1.HostPathVolumeSource{
                                                Path: "/data"
                                        },
                                },
                        },

		},
	}
}
