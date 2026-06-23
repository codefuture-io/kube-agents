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

// Package auth provides authentication and authorization for kube-agents.
package auth

import (
	"context"
	"fmt"

	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Identity represents an authenticated user from a verified token.
type Identity struct {
	UserID   string
	Groups   []string
	AuthMode string // "serviceaccount" or "jwt"
}

// TokenValidator validates bearer tokens and returns an Identity.
type TokenValidator interface {
	Validate(ctx context.Context, token string) (*Identity, error)
}

// SATokenReviewer validates K8s ServiceAccount tokens via TokenReview API.
type SATokenReviewer struct {
	client kubernetes.Interface
}

// NewSATokenReviewer creates a ServiceAccount token validator.
func NewSATokenReviewer(client kubernetes.Interface) *SATokenReviewer {
	return &SATokenReviewer{client: client}
}

// Validate performs a TokenReview against the K8s API server.
func (s *SATokenReviewer) Validate(ctx context.Context, token string) (*Identity, error) {
	review := &authv1.TokenReview{
		Spec: authv1.TokenReviewSpec{Token: token},
	}
	result, err := s.client.AuthenticationV1().TokenReviews().Create(ctx, review, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("token review failed: %w", err)
	}
	if !result.Status.Authenticated {
		return nil, fmt.Errorf("token not authenticated: %s", result.Status.Error)
	}
	return &Identity{
		UserID:   result.Status.User.Username,
		Groups:   result.Status.User.Groups,
		AuthMode: "serviceaccount",
	}, nil
}

// JWTValidator validates self-issued JWT tokens (HMAC or RSA).
type JWTValidator struct {
	issuer    string
	secretKey []byte
}

// NewJWTValidator creates a JWT token validator.
func NewJWTValidator(issuer string, secretKey []byte) *JWTValidator {
	return &JWTValidator{issuer: issuer, secretKey: secretKey}
}

// Validate parses and validates a JWT token.
// Full JWT implementation pending addition of jwt library.
func (j *JWTValidator) Validate(ctx context.Context, token string) (*Identity, error) {
	_ = j.issuer
	_ = j.secretKey
	return nil, fmt.Errorf("JWT validation not yet implemented — use serviceaccount mode")
}
