// Copyright 2023 runningriven@gmail.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// metrics module of business: collect metric in ws_business.go

package wsBusiness

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

// define some business metrics to monitor
var (
	UserCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "UserCount",
		Help: "The number of user who logins in",
	}, []string{"UserCount"})

	UserDetail = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "UserDetail",
		Help: "User's detail information",
	}, []string{"AccountName", "websocketID"})
)

func RegisterMonitorIndicators(registry *prometheus.Registry) {
	registry.MustRegister(UserCount)
	registry.MustRegister(UserDetail)
}

func RefreshUsersInfoMetrics() {
	refreshUserCountMetric()
	refreshUserDetailMetric()
}

func refreshUserCountMetric() {
	if gUserManager != nil {
		UserCount.Reset()
		UserCount.WithLabelValues(strconv.Itoa(gUserManager.UserCount()))
	}
}

func refreshUserDetailMetric() {
	if gUserManager != nil {
		users := gUserManager.GetUsers()
		UserDetail.Reset()
		var websocketClientID string
		for i := 0; i < len(users); i++ {
			websocketClientID = users[i].UserInfo.Client.ID()
			UserDetail.WithLabelValues(users[i].AccountName, websocketClientID)
		}
	}
}
