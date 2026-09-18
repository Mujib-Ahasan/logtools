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

package nilStringerDisabled

import (
	klog "k8s.io/klog/v2"
)

func nilUnsafeStringers() {
	var valueStringer *timestamp

	// Not reported, the nil-stringer check is disabled.
	klog.Background().Info("Starting", "time", valueStringer)
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
