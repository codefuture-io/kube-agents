/*
Copyright 2026 CodeFuture Authors

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

package auth

import (
	"context"
	"fmt"

	authzv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RBACChecker verifies whether an identity has permission to perform
// a specific operation on a Kubernetes resource.
type RBACChecker struct {
	client kubernetes.Interface
}

// NewRBACChecker creates an RBAC checker backed by K8s SubjectAccessReview.
func NewRBACChecker(client kubernetes.Interface) *RBACChecker {
	return &RBACChecker{client: client}
}

// AccessRequest describes a K8s access check.
type AccessRequest struct {
	Namespace    string
	Resource     string
	ResourceName string
	Verb         string // get, list, create, update, delete, watch
	APIGroup     string
}

// Check performs a SubjectAccessReview against the K8s API server.
func (r *RBACChecker) Check(ctx context.Context, identity *Identity, req AccessRequest) error {
	sar := &authzv1.SubjectAccessReview{
		Spec: authzv1.SubjectAccessReviewSpec{
			User:   identity.UserID,
			Groups: identity.Groups,
			ResourceAttributes: &authzv1.ResourceAttributes{
				Namespace: req.Namespace,
				Resource:  req.Resource,
				Name:      req.ResourceName,
				Verb:      req.Verb,
				Group:     req.APIGroup,
			},
		},
	}
	result, err := r.client.AuthorizationV1().SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("SubjectAccessReview failed: %w", err)
	}
	if !result.Status.Allowed {
		reason := result.Status.Reason
		if reason == "" {
			reason = "access denied by RBAC"
		}
		return fmt.Errorf("RBAC denied: %s", reason)
	}
	return nil
}
