/*
Copyright 2019 Platformers @bonniernews.

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

package event

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Setup logging
func init() {
	logf.SetLogger(logf.ZapLogger(false))
}

var logger = logf.Log.WithName("k8s-event-logger")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileEvent{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("event-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Event
	err = c.Watch(&source.Kind{Type: &corev1.Event{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by Event - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &corev1.Event{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileEvent{}

// ReconcileEvent reconciles a Event object
type ReconcileEvent struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Event object and makes changes based on the state read
// and what is in the Event.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events/status,verbs=get;update;patch
func (r *ReconcileEvent) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Event instance
	instance := &corev1.Event{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	/* 	fmt.Printf("%v:%v:%v (%v):%v:%v/%v -- %v\n", instance.InvolvedObject.Kind,
	instance.Type,
	instance.Source.Component,
	instance.Source.Host,
	instance.Namespace,
	instance.InvolvedObject.Namespace,
	instance.InvolvedObject.Name,
	instance.Message) */

	logger.Info("Event", "type", instance.InvolvedObject.Kind,
		"eventLevel", instance.Type,
		"source", instance.Source.Component,
		"host", instance.Source.Host,
		"namespace", instance.Namespace,
		"objectNamespace", instance.InvolvedObject.Namespace,
		"objectName", instance.InvolvedObject.Name,
		"eventMessage", instance.Message)

	return reconcile.Result{}, nil
}
