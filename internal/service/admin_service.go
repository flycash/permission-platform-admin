package service

import "C"
import (
	"context"
	"fmt"

	"gitee.com/flycash/permission-platform-admin/internal/domain"
	permissionv1 "gitee.com/flycash/permission-platform/api/proto/gen/permission/v1"
	"google.golang.org/grpc/metadata"
)

type AdminService struct {
	rbacSvc    permissionv1.RBACServiceClient
	perm       permissionv1.PermissionServiceClient
	adminToken string
}

// CreateBusinessConfig 业务方接入
// 初始化业务方的权限，并且初始化业务管理员
// 而后授予接入者业务管理员角色
func (svc *AdminService) CreateBusinessConfig(ctx context.Context, businessConfig domain.BusinessConfig) error {
	// 使用“系统管理员”权限创建业务配置
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", svc.adminToken)
	resp, err := svc.createBusinessConfig(ctx, businessConfig)
	if err != nil {
		return err
	}

	// 通过 ctx 模拟业务方的身份，完成后续步骤
	bizID := resp.Id
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", resp.Token)

	// 初始化”初始业务资源“
	resources, err2 := svc.createInitialBusinessResources(ctx, bizID)
	if err2 != nil {
		return err2
	}

	// 为每种”初始业务资源“创建相应权限
	permissions, err3 := svc.createPermissionsForInitialBusinessResources(ctx, bizID, resources)
	if err3 != nil {
		return err3
	}

	// 创建”业务管理员“角色
	adminRole, err4 := svc.createBusinessAdminRole(ctx, bizID)
	if err4 != nil {
		return err4
	}

	// 授予”业务管理员“角色，“初始业务资源”的全部权限
	err5 := svc.grantRolePermissions(ctx, bizID, adminRole, permissions)
	if err5 != nil {
		return err5
	}

	// 赋予用户”业务管理员“角色
	err6 := svc.grantUserRole(ctx, bizID, businessConfig.OwnerID, adminRole)
	if err6 != nil {
		return err6
	}
	return err
}

func (svc *AdminService) createBusinessConfig(ctx context.Context, businessConfig domain.BusinessConfig) (*permissionv1.BusinessConfig, error) {
	resp, err := svc.rbacSvc.CreateBusinessConfig(ctx, &permissionv1.CreateBusinessConfigRequest{
		Config: &permissionv1.BusinessConfig{
			OwnerId:   businessConfig.OwnerID,
			OwnerType: businessConfig.OwnerType,
			Name:      businessConfig.Name,
			RateLimit: int32(businessConfig.RateLimit),
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.Config, nil
}

func (svc *AdminService) createInitialBusinessResources(ctx context.Context, bizID int64) ([]domain.Resource, error) {
	// 将管理平台的7张表，作为业务内部资源初始化，但使用预定义的Type、Key和Name
	systemResources := []domain.SystemTableResource{
		domain.ResourceTable,
		domain.PermissionTable,
		domain.RoleTable,
		domain.RoleInclusionTable,
		domain.RolePermissionTable,
		domain.UserRoleTable,
		domain.UserPermissionTable,
	}
	resources := make([]domain.Resource, 0, len(systemResources)+1)
	for i := range systemResources {
		res, err := svc.rbacSvc.CreateResource(ctx, &permissionv1.CreateResourceRequest{
			Resource: &permissionv1.Resource{
				BizId: bizID,
				Type:  systemResources[i].Type(),
				Key:   systemResources[i].KeyForBusinessAdmin(bizID),
				Name:  systemResources[i].String(),
			},
		})
		if err != nil {
			return nil, err
		}
		resources = append(resources, svc.toResourceDomain(res.Resource))
	}

	// ”账号管理“资源也作为业务内部资源初始化，但使用预定义的Type、Key和Name
	res, err := svc.rbacSvc.CreateResource(ctx, &permissionv1.CreateResourceRequest{
		Resource: &permissionv1.Resource{
			BizId: bizID,
			Type:  domain.ManagerAccountResource.Type(),
			Key:   domain.ManagerAccountResource.KeyForBusinessAdmin(bizID),
			Name:  domain.ManagerAccountResource.String(),
		},
	})
	if err != nil {
		return nil, err
	}
	resources = append(resources, svc.toResourceDomain(res.Resource))

	return resources, nil
}

func (svc *AdminService) toResourceDomain(r *permissionv1.Resource) domain.Resource {
	return domain.Resource{
		ID:          r.Id,
		BizID:       r.BizId,
		Type:        r.Type,
		Key:         r.Key,
		Name:        r.Name,
		Description: r.Description,
		Metadata:    r.Metadata,
	}
}

func (svc *AdminService) createPermissionsForInitialBusinessResources(ctx context.Context, bizID int64, resources []domain.Resource) ([]domain.Permission, error) {
	systemResourcePermissions := []domain.PermissionActionType{
		domain.PermissionActionRead,
		domain.PermissionActionWrite,
	}
	permissions := make([]domain.Permission, 0, len(resources)*len(systemResourcePermissions))
	for i := range resources {
		// 为每张赋予预定义的权限
		for j := range systemResourcePermissions {
			resp, err := svc.rbacSvc.CreatePermission(ctx, &permissionv1.CreatePermissionRequest{
				Permission: &permissionv1.Permission{
					BizId:        bizID,
					Name:         fmt.Sprintf("%s-%s", resources[i].Name, systemResourcePermissions[j].String()),
					Description:  fmt.Sprintf("%s-%s", resources[i].Name, systemResourcePermissions[j].String()),
					ResourceId:   resources[i].ID,
					ResourceType: resources[i].Type,
					ResourceKey:  resources[i].Key,
					Actions:      []string{systemResourcePermissions[j].String()},
				},
			})
			if err != nil {
				return nil, err

			}
			permissions = append(permissions, svc.toPermissionDomain(resp.Permission))
		}
	}
	return permissions, nil
}

func (svc *AdminService) toPermissionDomain(permission *permissionv1.Permission) domain.Permission {
	return domain.Permission{
		ID:          permission.Id,
		BizID:       permission.BizId,
		Name:        permission.Name,
		Description: permission.Description,
		Resource: domain.Resource{
			ID:   permission.ResourceId,
			Type: permission.ResourceType,
			Key:  permission.ResourceKey,
		},
		Action:   permission.Actions[0],
		Metadata: permission.Metadata,
	}
}

func (svc *AdminService) createBusinessAdminRole(ctx context.Context, bizID int64) (domain.Role, error) {
	resp, err := svc.rbacSvc.CreateRole(ctx, &permissionv1.CreateRoleRequest{
		Role: &permissionv1.Role{
			BizId:       bizID,
			Type:        domain.DefaultAccountRoleType,
			Name:        "业务管理员",
			Description: "具有业务内最高管理权限",
		},
	})
	if err != nil {
		return domain.Role{}, err

	}
	return domain.Role{
		ID:          resp.Role.Id,
		BizID:       resp.Role.BizId,
		Type:        resp.Role.Type,
		Name:        resp.Role.Name,
		Description: resp.Role.Description,
		Metadata:    resp.Role.Metadata,
	}, nil
}

func (svc *AdminService) grantRolePermissions(ctx context.Context, bizID int64, role domain.Role, permissions []domain.Permission) error {
	for i := range permissions {
		_, err := svc.rbacSvc.GrantRolePermission(ctx, &permissionv1.GrantRolePermissionRequest{
			RolePermission: &permissionv1.RolePermission{
				BizId:            bizID,
				RoleId:           role.ID,
				PermissionId:     permissions[i].ID,
				RoleName:         role.Name,
				RoleType:         role.Type,
				ResourceType:     permissions[i].Resource.Type,
				ResourceKey:      permissions[i].Resource.Key,
				PermissionAction: permissions[i].Action,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *AdminService) grantUserRole(ctx context.Context, bizID int64, userID int64, role domain.Role) error {
	_, err := svc.rbacSvc.GrantUserRole(ctx, &permissionv1.GrantUserRoleRequest{
		UserRole: &permissionv1.UserRole{
			BizId:    bizID,
			UserId:   userID,
			RoleId:   role.ID,
			RoleName: role.Name,
			RoleType: role.Type,
		},
	})
	return err
}
