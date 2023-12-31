package setup

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func DeleteOrUpdate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc:  func(e event.CreateEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return true },
		UpdateFunc:  func(e event.UpdateEvent) bool { return true },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}

func CreateOrUpdate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc:  func(e event.CreateEvent) bool { return true },
		DeleteFunc:  func(e event.DeleteEvent) bool { return false },
		UpdateFunc:  func(e event.UpdateEvent) bool { return true },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}

func CreateOrUpdateOrDelete() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc:  func(e event.CreateEvent) bool { return true },
		DeleteFunc:  func(e event.DeleteEvent) bool { return true },
		UpdateFunc:  func(e event.UpdateEvent) bool { return true },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}

func CreateOrDelete() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc:  func(e event.CreateEvent) bool { return true },
		DeleteFunc:  func(e event.DeleteEvent) bool { return true },
		UpdateFunc:  func(e event.UpdateEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}
