// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package ippoolmanager_test

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	k8stesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	electionmock "github.com/spidernet-io/spiderpool/pkg/election/mock"
	"github.com/spidernet-io/spiderpool/pkg/ippoolmanager"
	spiderpoolv2beta1 "github.com/spidernet-io/spiderpool/pkg/k8s/apis/spiderpool.spidernet.io/v2beta1"
	reservedipmanagermock "github.com/spidernet-io/spiderpool/pkg/reservedipmanager/mock"
)

var mockCtrl *gomock.Controller
var mockLeaderElector *electionmock.MockSpiderLeaseElector
var mockRIPManager *reservedipmanagermock.MockReservedIPManager

var scheme *runtime.Scheme
var fakeClient client.Client
var tracker k8stesting.ObjectTracker
var fakeAPIReader client.Reader
var ipPoolManager ippoolmanager.IPPoolManager
var ipPoolWebhook *ippoolmanager.IPPoolWebhook

func TestIPPoolManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IPPoolManager Suite", Label("ippoolmanager", "unitest"))
}

var _ = BeforeSuite(func() {
	scheme = runtime.NewScheme()
	err := spiderpoolv2beta1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	fakeClient = fake.NewClientBuilder().
		WithScheme(scheme).
		WithIndex(&spiderpoolv2beta1.SpiderIPPool{}, metav1.ObjectNameField, func(raw client.Object) []string {
			ipPool := raw.(*spiderpoolv2beta1.SpiderIPPool)
			return []string{ipPool.GetObjectMeta().GetName()}
		}).
		WithIndex(&spiderpoolv2beta1.SpiderIPPool{}, "spec.default", func(raw client.Object) []string {
			ipPool := raw.(*spiderpoolv2beta1.SpiderIPPool)
			return []string{strconv.FormatBool(*ipPool.Spec.Default)}
		}).
		Build()

	tracker = k8stesting.NewObjectTracker(scheme, k8sscheme.Codecs.UniversalDecoder())
	fakeAPIReader = fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjectTracker(tracker).
		WithIndex(&spiderpoolv2beta1.SpiderIPPool{}, metav1.ObjectNameField, func(raw client.Object) []string {
			ipPool := raw.(*spiderpoolv2beta1.SpiderIPPool)
			return []string{ipPool.GetObjectMeta().GetName()}
		}).
		WithIndex(&spiderpoolv2beta1.SpiderIPPool{}, "spec.default", func(raw client.Object) []string {
			ipPool := raw.(*spiderpoolv2beta1.SpiderIPPool)
			return []string{strconv.FormatBool(*ipPool.Spec.Default)}
		}).
		Build()

	mockLeaderElector = electionmock.NewMockSpiderLeaseElector(mockCtrl)
	mockRIPManager = reservedipmanagermock.NewMockReservedIPManager(mockCtrl)
	ipPoolManager, err = ippoolmanager.NewIPPoolManager(
		ippoolmanager.IPPoolManagerConfig{},
		fakeClient,
		fakeAPIReader,
		mockRIPManager,
	)
	Expect(err).NotTo(HaveOccurred())

	ipPoolWebhook = &ippoolmanager.IPPoolWebhook{
		Client:             fakeClient,
		APIReader:          fakeAPIReader,
		EnableIPv4:         true,
		EnableIPv6:         true,
		EnableSpiderSubnet: true,
	}
})
