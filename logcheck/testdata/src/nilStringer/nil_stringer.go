/*
Copyright 2026 The Kubernetes Authors.

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

package nilStringer

import (
	klog "k8s.io/klog/v2"
)

func nilUnsafeStringers() {
	var valueStringer *timestamp
	var wrapperStringer *wrappedTimestamp
	var promotedStringer *hasEmbeddedNilSafeStringer
	var pointerStringer *nilSafeStringer
	var plainValue timestamp

	// Calling String() on a nil *timestamp panics because the method has a
	// value receiver and thus dereferences the pointer.
	klog.Background().Info("Starting", "time", valueStringer) // want `The type \*nilStringer.timestamp implements fmt.Stringer through \(nilStringer.timestamp\).String, which panics for a nil pointer because calling it dereferences the pointer. klog logs the panic instead of the value, other logger implementations may crash. Wrap the value with klog.SafePtr.`
	klog.InfoS("Starting", "time", valueStringer)             // want `The type \*nilStringer.timestamp implements fmt.Stringer through \(nilStringer.timestamp\).String, which panics for a nil pointer because calling it dereferences the pointer. klog logs the panic instead of the value, other logger implementations may crash. Wrap the value with klog.SafePtr.`

	// The single-field wrapper struct exemption of the "value" check must
	// not apply here: *metav1.Time, the prototypical instance of this bug,
	// is exactly such a wrapper.
	klog.Background().Info("Starting", "time", wrapperStringer) // want `The type \*nilStringer.wrappedTimestamp implements fmt.Stringer through \(nilStringer.timestamp\).String, which panics for a nil pointer because calling it dereferences the pointer. klog logs the panic instead of the value, other logger implementations may crash. Wrap the value with klog.SafePtr.`

	// String has a nil-safe pointer receiver, but it is promoted from a
	// field that is embedded by value: selecting that field dereferences
	// the nil outer pointer before the method can check its receiver.
	klog.Background().Info("Starting", "value", promotedStringer) // want `The type \*nilStringer.hasEmbeddedNilSafeStringer implements fmt.Stringer through \(\*nilStringer.nilSafeStringer\).String, which panics for a nil pointer because calling it dereferences the pointer. klog logs the panic instead of the value, other logger implementations may crash. Wrap the value with klog.SafePtr.`

	// A pointer receiver declared on the type itself can (and here does)
	// handle nil.
	klog.Background().Info("Starting", "value", pointerStringer)

	// A non-pointer value cannot be nil.
	klog.Background().Info("Starting", "time", plainValue)

	// Taking the address of a value or calling new never yields a nil
	// pointer, with or without parentheses.
	klog.Background().Info("Starting", "time", &plainValue)
	klog.Background().Info("Starting", "time", (&plainValue))
	klog.Background().Info("Starting", "time", new(timestamp))

	// SafePtr handles nil pointers, that is the suggested fix.
	klog.Background().Info("Starting", "time", klog.SafePtr(valueStringer))
}

// timestamp implements fmt.Stringer with a value receiver, like time.Time
// does. Calling String() through a nil pointer panics.
type timestamp struct {
	seconds int64
	nanos   int32
}

func (t timestamp) String() string {
	return "some point in time"
}

// wrappedTimestamp mimics metav1.Time: a single-field wrapper struct that
// inherits the value receiver String() of the embedded type.
type wrappedTimestamp struct {
	timestamp
}

// nilSafeStringer implements fmt.Stringer with a pointer receiver that
// handles nil, like *net.IPNet does.
type nilSafeStringer struct {
	value string
}

func (s *nilSafeStringer) String() string {
	if s == nil {
		return "<nil>"
	}
	return s.value
}

// hasEmbeddedNilSafeStringer embeds nilSafeStringer by value, so the
// promoted String() dereferences a *hasEmbeddedNilSafeStringer before the
// nil check in the method has a chance to run.
type hasEmbeddedNilSafeStringer struct {
	nilSafeStringer
}
