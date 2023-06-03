/*
Copyright 2020 Layer5, Inc.

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
	"github.com/layer5io/meshkit/errors"
)

// Error codes
const (
	ErrReconcileCRCode    = "7001"
	ErrCheckHealthCode    = "7002"
	ErrCreateResourceCode = "7003"
	ErrGetResourceCode    = "7004"
	ErrUpdateResourceCode = "7005"
	ErrDeleteResourceCode = "7006"
)

// Error definitions

func ErrReconcileAdapter(err error) error {
	return errors.New(ErrReconcileCRCode, errors.Alert, []string{"Error during adaper reconciliation"}, []string{err.Error()}, []string{}, []string{})
}

func ErrCheckHealth(err error) error {
	return errors.New(ErrCheckHealthCode, errors.Alert, []string{"Error during health check"}, []string{err.Error()}, []string{}, []string{})
}

func ErrCreateAdapter(err error) error {
	return errors.New(ErrCreateResourceCode, errors.Alert, []string{"Error creating adapter"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGetAdapter(err error) error {
	return errors.New(ErrGetResourceCode, errors.Alert, []string{"Error getting adapter"}, []string{err.Error()}, []string{}, []string{})
}

func ErrUpdateAdapter(err error) error {
	return errors.New(ErrUpdateResourceCode, errors.Alert, []string{"Error updating adapter"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDeletingAdapter(err error) error {
	return errors.New(ErrDeleteResourceCode, errors.Alert, []string{"Error deleting adapter"}, []string{err.Error()}, []string{}, []string{})
}
