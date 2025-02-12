// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package subnetmanager_test

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spidernet-io/spiderpool/pkg/constant"
	spiderpoolv2beta1 "github.com/spidernet-io/spiderpool/pkg/k8s/apis/spiderpool.spidernet.io/v2beta1"
	"github.com/spidernet-io/spiderpool/pkg/subnetmanager"
)

var _ = Describe("SubnetManager", Label("subnet_manager_test"), func() {
	Describe("New SubnetManager", func() {
		It("inputs nil client", func() {
			manager, err := subnetmanager.NewSubnetManager(nil, fakeAPIReader, mockRIPManager)
			Expect(err).To(MatchError(constant.ErrMissingRequiredParam))
			Expect(manager).To(BeNil())
		})

		It("inputs nil API reader", func() {
			manager, err := subnetmanager.NewSubnetManager(fakeClient, nil, mockRIPManager)
			Expect(err).To(MatchError(constant.ErrMissingRequiredParam))
			Expect(manager).To(BeNil())
		})

		It("inputs nil reserved-IP manager", func() {
			manager, err := subnetmanager.NewSubnetManager(fakeClient, fakeAPIReader, nil)
			Expect(err).To(MatchError(constant.ErrMissingRequiredParam))
			Expect(manager).To(BeNil())
		})
	})

	Describe("Test SubnetManager's method", func() {
		var ctx context.Context

		var count uint64
		var subnetName string
		var labels map[string]string
		var subnetT *spiderpoolv2beta1.SpiderSubnet

		BeforeEach(func() {
			ctx = context.TODO()

			atomic.AddUint64(&count, 1)
			subnetName = fmt.Sprintf("subnet-%v", count)
			labels = map[string]string{"foo": fmt.Sprintf("bar-%v", count)}
			subnetT = &spiderpoolv2beta1.SpiderSubnet{
				TypeMeta: metav1.TypeMeta{
					Kind:       constant.KindSpiderSubnet,
					APIVersion: fmt.Sprintf("%s/%s", constant.SpiderpoolAPIGroup, constant.SpiderpoolAPIVersion),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   subnetName,
					Labels: labels,
				},
				Spec: spiderpoolv2beta1.SubnetSpec{},
			}
		})

		var deleteOption *client.DeleteOptions

		AfterEach(func() {
			policy := metav1.DeletePropagationForeground
			deleteOption = &client.DeleteOptions{
				GracePeriodSeconds: pointer.Int64(0),
				PropagationPolicy:  &policy,
			}

			err := fakeClient.Delete(ctx, subnetT, deleteOption)
			Expect(client.IgnoreNotFound(err)).NotTo(HaveOccurred())

			err = tracker.Delete(
				schema.GroupVersionResource{
					Group:    constant.SpiderpoolAPIGroup,
					Version:  constant.SpiderpoolAPIVersion,
					Resource: "spidersubnets",
				},
				subnetT.Namespace,
				subnetT.Name,
			)
			Expect(client.IgnoreNotFound(err)).NotTo(HaveOccurred())
		})

		Describe("GetSubnetByName", func() {
			It("gets non-existent Subnet", func() {
				subnet, err := subnetManager.GetSubnetByName(ctx, subnetName, constant.IgnoreCache)
				Expect(apierrors.IsNotFound(err)).To(BeTrue())
				Expect(subnet).To(BeNil())
			})

			It("gets an existing Subnet through cache", func() {
				err := fakeClient.Create(ctx, subnetT)
				Expect(err).NotTo(HaveOccurred())

				subnet, err := subnetManager.GetSubnetByName(ctx, subnetName, constant.UseCache)
				Expect(err).NotTo(HaveOccurred())
				Expect(subnet).NotTo(BeNil())
				Expect(subnet).To(Equal(subnetT))
			})

			It("gets an existing Subnet through API Server", func() {
				err := tracker.Add(subnetT)
				Expect(err).NotTo(HaveOccurred())

				subnet, err := subnetManager.GetSubnetByName(ctx, subnetName, constant.IgnoreCache)
				Expect(err).NotTo(HaveOccurred())
				Expect(subnet).NotTo(BeNil())
				Expect(subnet).To(Equal(subnetT))
			})
		})

		Describe("ListSubnets", func() {
			It("failed to list Subnets due to some unknown errors", func() {
				patches := gomonkey.ApplyMethodReturn(fakeAPIReader, "List", constant.ErrUnknown)
				defer patches.Reset()

				err := tracker.Add(subnetT)
				Expect(err).NotTo(HaveOccurred())

				subnetList, err := subnetManager.ListSubnets(ctx, constant.IgnoreCache)
				Expect(err).To(MatchError(constant.ErrUnknown))
				Expect(subnetList).To(BeNil())
			})

			It("lists all Subnets through cache", func() {
				err := fakeClient.Create(ctx, subnetT)
				Expect(err).NotTo(HaveOccurred())

				subnetList, err := subnetManager.ListSubnets(ctx, constant.UseCache)
				Expect(err).NotTo(HaveOccurred())
				Expect(subnetList.Items).NotTo(BeEmpty())

				hasSubnet := false
				for _, subnet := range subnetList.Items {
					if subnet.Name == subnetName {
						hasSubnet = true
						break
					}
				}
				Expect(hasSubnet).To(BeTrue())
			})

			It("lists all Subnets through API Server", func() {
				err := tracker.Add(subnetT)
				Expect(err).NotTo(HaveOccurred())

				subnetList, err := subnetManager.ListSubnets(ctx, constant.IgnoreCache)
				Expect(err).NotTo(HaveOccurred())
				Expect(subnetList.Items).NotTo(BeEmpty())

				hasSubnet := false
				for _, subnet := range subnetList.Items {
					if subnet.Name == subnetName {
						hasSubnet = true
						break
					}
				}
				Expect(hasSubnet).To(BeTrue())
			})

			It("filters results by label selector", func() {
				err := tracker.Add(subnetT)
				Expect(err).NotTo(HaveOccurred())

				subnetList, err := subnetManager.ListSubnets(ctx, constant.IgnoreCache, client.MatchingLabels(labels))
				Expect(err).NotTo(HaveOccurred())
				Expect(subnetList.Items).NotTo(BeEmpty())

				hasSubnet := false
				for _, subnet := range subnetList.Items {
					if subnet.Name == subnetName {
						hasSubnet = true
						break
					}
				}
				Expect(hasSubnet).To(BeTrue())
			})

			It("filters results by field selector", func() {
				err := tracker.Add(subnetT)
				Expect(err).NotTo(HaveOccurred())

				subnetList, err := subnetManager.ListSubnets(ctx, constant.IgnoreCache, client.MatchingFields{metav1.ObjectNameField: subnetName})
				Expect(err).NotTo(HaveOccurred())
				Expect(subnetList.Items).NotTo(BeEmpty())

				hasSubnet := false
				for _, subnet := range subnetList.Items {
					if subnet.Name == subnetName {
						hasSubnet = true
						break
					}
				}
				Expect(hasSubnet).To(BeTrue())
			})
		})
	})
})
