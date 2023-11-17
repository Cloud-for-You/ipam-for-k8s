/*
Copyright 2023.

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

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ipamv1 "github.com/Cloud-for-You/ipam-for-k8s/api/v1"
	pkgSubnet "github.com/Cloud-for-You/ipam-for-k8s/pkg/subnet"
)

const (
	subnetFinalizer string = "ipam.cfy.cz/finalizer"
)

// SubnetReconciler reconciles a Subnet object
type SubnetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ipam.cfy.cz,resources=subnets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ipam.cfy.cz,resources=subnets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ipam.cfy.cz,resources=subnets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Subnet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *SubnetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	log := log.FromContext(ctx)
	log.Info("Verify if a CRD of Subnet exists")

	// Fetch the Subnet instance
	subnet := &ipamv1.Subnet{}
	if err := r.Get(ctx, req.NamespacedName, subnet); err != nil {
		log.Error(err, "unable to fetch Subnet")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if subnet.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(subnet, subnetFinalizer) {
			controllerutil.AddFinalizer(subnet, subnetFinalizer)
			if err := r.Update(ctx, subnet); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	} else {
		if controllerutil.ContainsFinalizer(subnet, subnetFinalizer) {
			if err := r.finalizeSubnet(subnet); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(subnet, subnetFinalizer)
			if err := r.Update(ctx, subnet); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Nase logika kodu
	totalCount := 0
	usedCount := 0
	reservedCount := 10

	for _, ipRange := range subnet.Spec.UsableIPs {
		count, err := pkgSubnet.GetUsedIPsInSubnet(ipRange)
		if err != nil {
			fmt.Printf("Error processing IP range %s: %v\n", ipRange, err)
			continue
		}
		fmt.Printf("IP range %s has %d addresses\n", ipRange, count)
		totalCount += count
	}

	// Update the status of the Subnet instance
	if err := r.updateStatus(subnet, int(totalCount), int(usedCount), int(reservedCount)); err != nil {
		log.Error(err, "unable to update Subnet status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SubnetReconciler) updateStatus(subnet *ipamv1.Subnet, totalAddress int, usedAddress int, reservedAddress int) error {
	subnet.Status.TotalAddresses = totalAddress
	subnet.Status.UsedAddresses = usedAddress
	subnet.Status.ReserverdAddresses = reservedAddress
	subnet.Status.FreeAddresses = totalAddress - usedAddress - reservedAddress

	// Save the updated status
	if err := r.Status().Update(context.Background(), subnet); err != nil {
		return err
	}

	return nil
}

func (r *SubnetReconciler) finalizeSubnet(m *ipamv1.Subnet) error {
	log.Log.Info("Successfuly finalize subnet")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SubnetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ipamv1.Subnet{}).
		Complete(r)
}
