// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package dao

import (
	"errors"
	"testing"

	"github.com/derailed/k9s/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsStandardGroup(t *testing.T) {
	uu := map[string]struct {
		gv string
		e  bool
	}{
		"core": {gv: "v1", e: true},
		"apps": {gv: "apps/v1", e: true},
		"batch": {gv: "batch/v1", e: true},
		"networking": {gv: "networking.k8s.io/v1", e: true},
		"storage": {gv: "storage.k8s.io/v1", e: true},
		"rbac": {gv: "rbac.authorization.k8s.io/v1", e: true},
		"flowcontrol": {gv: "flowcontrol.apiserver.k8s.io/v1", e: true},
		"cluster-api": {gv: "cluster.x-k8s.io/v1beta2", e: false},
		"cluster-api-infra": {gv: "infrastructure.cluster.x-k8s.io/v1beta2", e: false},
		"cluster-api-addons": {gv: "addons.cluster.x-k8s.io/v1beta1", e: false},
		"traefik": {gv: "traefik.io/v1alpha1", e: false},
		"monitoring": {gv: "monitoring.coreos.com/v1", e: false},
	}

	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			assert.Equal(t, u.e, isStandardGroup(u.gv))
		})
	}
}

func TestMetaFor(t *testing.T) {
	uu := map[string]struct {
		gvr *client.GVR
		err error
		e   metav1.APIResource
	}{
		"xray-gvr": {
			gvr: client.XGVR,
			e: metav1.APIResource{
				Name:         "xrays",
				Kind:         "XRays",
				SingularName: "xray",
				Categories:   []string{k9sCat},
			},
		},

		"xray": {
			gvr: client.NewGVR("xrays"),
			e: metav1.APIResource{
				Name:         "xrays",
				Kind:         "XRays",
				SingularName: "xray",
				Categories:   []string{k9sCat},
			},
		},

		"policy": {
			gvr: client.NewGVR("policy"),
			e: metav1.APIResource{
				Name:       "policies",
				Kind:       "Rules",
				Namespaced: true,
				Categories: []string{k9sCat},
			},
		},

		"helm": {
			gvr: client.NewGVR("helm"),
			e: metav1.APIResource{
				Name:       "helm",
				Kind:       "Helm",
				Namespaced: true,
				Verbs:      []string{"delete"},
				Categories: []string{helmCat},
			},
		},

		"toast": {
			gvr: client.NewGVR("blah"),
			err: errors.New("no resource meta defined for\n \"blah\""),
		},
	}

	m := NewMeta()
	require.NoError(t, m.LoadResources(nil))
	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			meta, err := m.MetaFor(u.gvr)
			assert.Equal(t, u.err, err)
			if err == nil {
				assert.Equal(t, &u.e, meta)
			}
		})
	}
}
