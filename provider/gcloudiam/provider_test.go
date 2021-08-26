package gcloudiam_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/odpf/guardian/domain"
	"github.com/odpf/guardian/mocks"
	"github.com/odpf/guardian/provider/gcloudiam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetType(t *testing.T) {
	t.Run("should return provider type name", func(t *testing.T) {
		expectedTypeName := domain.ProviderTypeMetabase
		crypto := new(mocks.Crypto)
		p := gcloudiam.NewProvider(expectedTypeName, crypto)

		actualTypeName := p.GetType()

		assert.Equal(t, expectedTypeName, actualTypeName)
	})
}

func TestGetResources(t *testing.T) {
	t.Run("should return one item of resource", func(t *testing.T) {
		crypto := new(mocks.Crypto)
		p := gcloudiam.NewProvider("", crypto)
		pc := &domain.ProviderConfig{
			Type: domain.ProviderTypeGCloudIAM,
			URN:  "test-project-id",
		}

		expectedResources := []*domain.Resource{
			{
				ProviderType: pc.Type,
				ProviderURN:  pc.URN,
				Type:         gcloudiam.ResourceTypeGcloudIam,
				URN:          pc.URN,
				Name:         fmt.Sprintf("%s - GCP IAM", pc.URN),
			},
		}

		actualResources, actualError := p.GetResources(pc)

		assert.Equal(t, expectedResources, actualResources)
		assert.Nil(t, actualError)
	})
}

func TestGrantAccess(t *testing.T) {
	t.Run("should return error if credentials is invalid", func(t *testing.T) {
		crypto := new(mocks.Crypto)
		p := gcloudiam.NewProvider("", crypto)

		pc := &domain.ProviderConfig{
			Credentials: "invalid-credentials",
			Resources: []*domain.ResourceConfig{
				{
					Type: "test-type",
					Roles: []*domain.RoleConfig{
						{
							ID: "test-role",
							Permissions: []interface{}{
								gcloudiam.PermissionConfig{
									Name: "test-permission-config",
								},
							},
						},
					},
				},
			},
		}
		a := &domain.Appeal{
			Resource: &domain.Resource{
				Type: "test-type",
			},
			Role: "test-role",
		}

		actualError := p.GrantAccess(pc, a)
		assert.Error(t, actualError)

	})

	t.Run("should return error if resource type is unknown", func(t *testing.T) {
		providerURN := "test-provider-urn"
		crypto := new(mocks.Crypto)
		client := new(mocks.GcloudIamClient)
		p := gcloudiam.NewProvider("", crypto)
		p.Clients = map[string]gcloudiam.GcloudIamClient{
			providerURN: client,
		}
		expectedError := errors.New("invalid resource type")

		pc := &domain.ProviderConfig{
			URN: providerURN,
			Resources: []*domain.ResourceConfig{
				{
					Type: "test-type",
					Roles: []*domain.RoleConfig{
						{
							ID: "test-role",
							Permissions: []interface{}{
								gcloudiam.PermissionConfig{
									Name: "test-permission-config",
								},
							},
						},
					},
				},
			},
		}
		a := &domain.Appeal{
			Resource: &domain.Resource{
				Type: "test-type",
			},
			Role: "test-role",
		}

		actualError := p.GrantAccess(pc, a)

		assert.EqualError(t, actualError, expectedError.Error())
	})

	t.Run("given database resource", func(t *testing.T) {
		t.Run("should return error if there is an error in granting the access", func(t *testing.T) {
			providerURN := "test-provider-urn"
			expectedError := errors.New("client error")
			crypto := new(mocks.Crypto)
			client := new(mocks.GcloudIamClient)
			p := gcloudiam.NewProvider("", crypto)
			p.Clients = map[string]gcloudiam.GcloudIamClient{
				providerURN: client,
			}
			client.On("GrantAccess", mock.Anything, mock.Anything).Return(expectedError).Once()

			pc := &domain.ProviderConfig{
				Resources: []*domain.ResourceConfig{
					{
						Type: gcloudiam.ResourceTypeRole,
						Roles: []*domain.RoleConfig{
							{
								ID: "test-role",
								Permissions: []interface{}{
									gcloudiam.PermissionConfig{
										Name: "test-permission-config",
									},
								},
							},
						},
					},
				},
				URN: providerURN,
			}
			a := &domain.Appeal{
				Resource: &domain.Resource{
					Type: gcloudiam.ResourceTypeRole,
					URN:  "999",
					Name: "test-role",
				},
				Role: "test-role",
			}

			actualError := p.GrantAccess(pc, a)

			assert.EqualError(t, actualError, expectedError.Error())
		})

		t.Run("should return nil error if granting access is successful", func(t *testing.T) {
			providerURN := "test-provider-urn"
			crypto := new(mocks.Crypto)
			client := new(mocks.GcloudIamClient)
			expectedResource := &gcloudiam.Role{
				Name: "test-role",
			}
			expectedUser := "test@email.com"
			p := gcloudiam.NewProvider("", crypto)
			p.Clients = map[string]gcloudiam.GcloudIamClient{
				providerURN: client,
			}
			client.On("GrantAccess", expectedResource, expectedUser).Return(nil).Once()

			pc := &domain.ProviderConfig{
				Resources: []*domain.ResourceConfig{
					{
						Type: gcloudiam.ResourceTypeRole,
						Roles: []*domain.RoleConfig{
							{
								ID: "allow",
								Permissions: []interface{}{
									gcloudiam.PermissionConfig{
										Name: "allow",
									},
								},
							},
						},
					},
				},
				URN: providerURN,
			}
			a := &domain.Appeal{
				Resource: &domain.Resource{
					Type: gcloudiam.ResourceTypeRole,
					URN:  "test-role",
				},
				Role:       "viewer",
				User:       expectedUser,
				ResourceID: 999,
				ID:         999,
			}

			actualError := p.GrantAccess(pc, a)

			assert.Nil(t, actualError)
		})
	})
}

func TestRevokeAccess(t *testing.T) {
	t.Run("should return error if resource type is unknown", func(t *testing.T) {
		providerURN := "test-provider-urn"
		crypto := new(mocks.Crypto)
		client := new(mocks.GcloudIamClient)
		p := gcloudiam.NewProvider("", crypto)
		p.Clients = map[string]gcloudiam.GcloudIamClient{
			providerURN: client,
		}
		expectedError := errors.New("invalid resource type")

		pc := &domain.ProviderConfig{
			URN: providerURN,
			Resources: []*domain.ResourceConfig{
				{
					Type: "test-type",
					Roles: []*domain.RoleConfig{
						{
							ID: "test-role",
							Permissions: []interface{}{
								gcloudiam.PermissionConfig{
									Name: "test-permission-config",
								},
							},
						},
					},
				},
			},
		}
		a := &domain.Appeal{
			Resource: &domain.Resource{
				Type: "test-type",
			},
			Role: "test-role",
		}

		actualError := p.RevokeAccess(pc, a)

		assert.EqualError(t, actualError, expectedError.Error())
	})

	t.Run("given database resource", func(t *testing.T) {
		t.Run("should return error if there is an error in granting the access", func(t *testing.T) {
			providerURN := "test-provider-urn"
			expectedError := errors.New("client error")
			crypto := new(mocks.Crypto)
			client := new(mocks.GcloudIamClient)
			p := gcloudiam.NewProvider("", crypto)
			p.Clients = map[string]gcloudiam.GcloudIamClient{
				providerURN: client,
			}
			client.On("RevokeAccess", mock.Anything, mock.Anything).Return(expectedError).Once()

			pc := &domain.ProviderConfig{
				Resources: []*domain.ResourceConfig{
					{
						Type: gcloudiam.ResourceTypeRole,
						Roles: []*domain.RoleConfig{
							{
								ID: "test-role",
								Permissions: []interface{}{
									gcloudiam.PermissionConfig{
										Name: "test-permission-config",
									},
								},
							},
						},
					},
				},
				URN: providerURN,
			}
			a := &domain.Appeal{
				Resource: &domain.Resource{
					Type: gcloudiam.ResourceTypeRole,
					URN:  "999",
					Name: "test-role",
				},
				Role: "test-role",
			}

			actualError := p.RevokeAccess(pc, a)

			assert.EqualError(t, actualError, expectedError.Error())
		})

		t.Run("should return nil error if revoking access is successful", func(t *testing.T) {
			providerURN := "test-provider-urn"
			crypto := new(mocks.Crypto)
			client := new(mocks.GcloudIamClient)
			expectedResource := &gcloudiam.Role{
				Name: "test-role",
			}
			expectedUser := "test@email.com"
			p := gcloudiam.NewProvider("", crypto)
			p.Clients = map[string]gcloudiam.GcloudIamClient{
				providerURN: client,
			}
			client.On("RevokeAccess", expectedResource, expectedUser).Return(nil).Once()

			pc := &domain.ProviderConfig{
				Resources: []*domain.ResourceConfig{
					{
						Type: gcloudiam.ResourceTypeRole,
						Roles: []*domain.RoleConfig{
							{
								ID: "allow",
								Permissions: []interface{}{
									gcloudiam.PermissionConfig{
										Name: "allow",
									},
								},
							},
						},
					},
				},
				URN: providerURN,
			}
			a := &domain.Appeal{
				Resource: &domain.Resource{
					Type: gcloudiam.ResourceTypeRole,
					URN:  "test-role",
				},
				Role:       "viewer",
				User:       expectedUser,
				ResourceID: 999,
				ID:         999,
			}

			actualError := p.RevokeAccess(pc, a)

			assert.Nil(t, actualError)
		})
	})
}
