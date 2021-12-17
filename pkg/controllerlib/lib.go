package controllerlib

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InstanceOpts struct {
	Controller string
	Request    string
	client.Object
	Error error
	*patch.Helper
}

func PatchInstance(ctx context.Context, i InstanceOpts) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		log = ctrl.Log
	}

	if err := i.validate(); err != nil {
		return
	}
	keysAndValues := []interface{}{
		"requestInfo", i.Request, "controller", i.Controller,
	}
	log.Info("Patching object instance", keysAndValues...)
	var patchOpts []patch.Option
	if i.Error == nil {
		patchOpts = append(patchOpts, patch.WithStatusObservedGeneration{})
	}
	patchErr := i.Helper.Patch(ctx, i.Object, patchOpts...)
	if patchErr != nil {
		i.Error = kerrors.NewAggregate([]error{patchErr, i.Error})
		keysAndValues = []interface{}{
			"requestInfo", i.Request, "controller", i.Controller, "err", i.Error,
		}
		log.Info("Error patching object instance", keysAndValues...)
		return
	}
	log.Info("Object instance patched", keysAndValues...)
}

func (i *InstanceOpts) validate() (err error) {
	if i.Controller == "" {
		return errors.New("Controller name empty")
	}
	if i.Request == "" {
		return errors.New("Object name empty")
	}
	if i.Object == nil {
		return errors.New("Object is nil")
	}
	return
}
