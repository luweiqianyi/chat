// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Create a customized prometheus.Registry,bind to a gin.Engine object, used by other modules to
// register some prometheus metrics to monitor

// by visit url path "/metrics" to view metrics which other modules registered to prometheus.Registry object

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var once sync.Once
var gRegistry *prometheus.Registry

func GetRegistry() *prometheus.Registry {
	once.Do(func() {
		gRegistry = prometheus.NewRegistry()
	})
	return gRegistry
}
