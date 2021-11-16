package app

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Instance struct {
	Ctx        context.Context
	Log        logr.Logger
	Controller string
	Request    string
	client.Object
	Error error
	*patch.Helper
}

func patchInstance(i Instance) {
	if err := i.validate(); err != nil {
		return
	}
	keysAndValues := []interface{}{
		"requestInfo", i.Request, "controller", i.Controller,
	}
	i.Log.Info("Patching object instance", keysAndValues...)
	var patchOpts []patch.Option
	if i.Error == nil {
		patchOpts = append(patchOpts, patch.WithStatusObservedGeneration{})
	}
	patchErr := i.Helper.Patch(i.Ctx, i.Object, patchOpts...)
	if patchErr != nil {
		i.Error = kerrors.NewAggregate([]error{patchErr, i.Error})
		i.Log.Info("Error patching object instance", keysAndValues...)
		return
	}
	i.Log.Info("Object instance patched", keysAndValues...)
}

func (i *Instance) validate() (err error) {
	if i.Controller == "" {
		return errors.New("Controller name empty")
	}
	if i.Request == "" {
		return errors.New("Object name empty")
	}
	if i.Log == nil {
		return errors.New("Log is nil")
	}
	if i.Object == nil {
		return errors.New("Object is nil")
	}
	return
}
