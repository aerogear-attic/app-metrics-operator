package controller

import (
	"github.com/aerogear/app-metrics-operator/pkg/controller/appmetricsconfig"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, appmetricsconfig.Add)
}
