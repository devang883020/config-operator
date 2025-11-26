package controllers

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	configv1alpha1 "github.com/devang/config-auto-update-operator/api/v1alpha1"
)

type AppConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *AppConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	var appConfig configv1alpha1.AppConfig
	if err := r.Get(ctx, req.NamespacedName, &appConfig); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appConfig.Name + "-config",
			Namespace: req.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, cm, func() error {
		cm.Data = appConfig.Spec.ConfigData
		return nil
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	dep := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      appConfig.Spec.TargetDeployment,
		Namespace: req.Namespace,
	}, dep)

	if err == nil {
		if dep.Spec.Template.Annotations == nil {
			dep.Spec.Template.Annotations = map[string]string{}
		}

		dep.Spec.Template.Annotations["config-reloaded"] = time.Now().Format(time.RFC3339)
		r.Update(ctx, dep)
	}

	appConfig.Status.LastApplied = metav1.Now().Time
	r.Status().Update(ctx, &appConfig)

	return ctrl.Result{}, nil
}

func (r *AppConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&configv1alpha1.AppConfig{}).
		Complete(r)
}
