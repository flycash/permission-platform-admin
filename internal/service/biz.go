package service

import "C"
import (
	"context"
	"errors"
	"fmt"
	"gitee.com/flycash/permission-platform-admin/internal/domain"
	permissionv1 "gitee.com/flycash/permission-platform/api/proto/gen/permission/v1"
	"google.golang.org/grpc/metadata"
)

type AdminService struct {
	bizSvc     permissionv1.BizServiceClient
	rbacSvc    permissionv1.RBACServiceClient
	perm       permissionv1.PermissionServiceClient
	adminToken string
}

func (svc *AdminService) ListBizByUid(ctx context.Context, uid int64) (domain.BizConfig, error) {
	// 我查询 uid 用户有类型为 biz 的有权限的资源
	// biz 是 resource type
	// select * from roles where uid = xx
	// select permissions where role_id IN () AND type = 'biz'
	// select resource from resources where rid = 'xxx'
	rs0, err := svc.rbacSvc.ListResourcesByType(ctx, uid, "biz")
	if err != nil {

	}
	var rs1 []permissionv1.Resource
	var bizIDs []int64
	for _, r := range rs1 {
		bizID := r.Metadata["bizID"]
		bizID = append(bizIDs, bizID)
	}
	bizCofigs := svc.bizSvc.ListBizs(bizIDs)
	return bizCofigs, nil
}

// CreateBiz 创建一个 Biz
func (svc *AdminService) CreateBiz(ctx context.Context, biz domain.BizConfig) error {
	// biz_id = 1 的 token
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", svc.adminToken)
	resp, err := svc.bizSvc.Create(ctx, &permissionv1.CreateRequest{
		Name: biz.Name,
		//Owner:     123,
		//OwnerType: "personal",
	})

	// 模拟业务方的身份去初始化最开始的资源、权限和角色
	// 后续都使用这个业务的 token
	// biz_id = 2 的token
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", resp.Token)

	// 初始化各种资源
	// role, resource, permission
	// grant: 授权与撤销授权
	const resourceType = "amdin"
	resources := [3]string{"role", "resource", "permission"}
	for _, resource := range resources {
		key := "gitee.com/flycash/permission-platform-admin/%d/%s"
		// key := "gitee.com/flycash/permission-platform-admin/2/role"
		res, err := svc.rbacSvc.CreateResource(ctx, &permissionv1.CreateResourceRequest{
			Resource: &permissionv1.Resource{
				Key:  fmt.Sprintf(key, resources, resp.Id),
				Type: resourceType,
				Name: "角色操作",
				// 理论上是从 token 解析的
				BizId: 2,
			},
		})
		if err != nil {
			return err
		}

		// 创建初始化的各种权限，增删改查初始一编
		perm, err := svc.rbacSvc.CreatePermission(ctx, &permissionv1.CreatePermissionRequest{
			Permission: &permissionv1.Permission{
				ResourceId:   res.GetResource().Id,
				ResourceKey:  key,
				ResourceType: resourceType,
				BizId:        2,
			},
		})
		if err != nil {
			return err
		}
	}

	perm, err := svc.rbacSvc.CreatePermission(ctx, &permissionv1.CreatePermissionRequest{
		Permission: &permissionv1.Permission{
			// 2 是选择的 biz_id
			ResourceKey:  "gitee.com/flycash/permission-platform-admin/%d/permission",
			ResourceType: resourceType,
			BizId:        2,
			Actions:      []string{"GRANT"},
		},
	})

	res, err := svc.rbacSvc.CreateResource(ctx, &permissionv1.CreateResourceRequest{
		Resource: &permissionv1.Resource{
			Key:  fmt.Sprintf("gitee.com/flycash/permission-platform-admin/%d", resp.Id),
			Type: "biz",
			Name: "角色操作",
			// 理论上是从 token 解析的
			BizId: 2,
		},
	})

	roleResp, err := svc.rbacSvc.CreateRole(ctx, &permissionv1.CreateRoleRequest{
		Role: &permissionv1.Role{
			Type:  "custom",
			Name:  "业务管理员",
			BizId: 2,
		},
	})

	if err != nil {
		return err
	}

	return err
}

func (svc *AdminService) BizGrantRole(ctx context.Context, bizID int64, uid int64) {
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", svc.adminToken)
	bizCfg := svc.bizSvc.GetBiz(ctx, bizID)

	// token 用来做业务方之间的隔离
	// 主动调用 GetPermission 来鉴权
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", bizCfg.Token)
	p, err := svc.perm.CheckPermission(ctx, &permissionv1.CheckPermissionRequest{
		Uid: uid,
		Permission: &permissionv1.Permission{
			ResourceKey: fmt.Sprintf("gitee.com/flycash/permission-platform-admin/%d/permission", bizID),
			Actions:     []string{"GRANT"},
		},
	})
	if err == nil && p.Allowed {
		roleResp, err := svc.rbacSvc.CreateRole(ctx, &permissionv1.CreateRoleRequest{
			Role: &permissionv1.Role{
				Type:  "custom",
				Name:  "渣渣辉",
				BizId: bizID,
			},
		})
		return roleResp, err
	}
}

// 业务方创建自己的角色
func (svc *AdminService) BizCreateRole(ctx context.Context, bizID int64, uid int64) {
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", svc.adminToken)
	bizCfg := svc.bizSvc.GetBiz(ctx, bizID)

	// token 用来做业务方之间的隔离
	// 主动调用 GetPermission 来鉴权
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", bizCfg.Token)
	p, err := svc.perm.CheckPermission(ctx, &permissionv1.CheckPermissionRequest{
		Uid: uid,
		Permission: &permissionv1.Permission{
			ResourceKey: fmt.Sprintf("gitee.com/flycash/permission-platform-admin/%d/role", bizID),
			Actions:     []permissionv1.ActionType{permissionv1.ActionType_WRITE},
		},
	})
	if err == nil && p.Allowed {
		roleResp, err := svc.rbacSvc.CreateRole(ctx, &permissionv1.CreateRoleRequest{
			Role: &permissionv1.Role{
				Type:  "custom",
				Name:  "渣渣辉",
				BizId: bizID,
			},
		})
		return roleResp, err
	}

}

func (svc *AdminService) MyBiz(ctx context.Context, uid int64) {
	const token = "xxx"
	// 假如说你这里的 BizID  = 3
	// 所有的权限校验都局限在 biz_id 内，就是查询一定带 biz_id = xxx 这个条件
	// biz_id 来源从 token 解析
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", token)
	p, err := svc.perm.CheckPermission(ctx, &permissionv1.CheckPermissionRequest{
		Uid: uid,
		Permission: &permissionv1.Permission{
			// ResourceKey => biz_id = 4 的，你能访问到吗？
			ResourceKey: "/mybiz/table/user_tab",
			Actions:     []permissionv1.ActionType{permissionv1.ActionType_WRITE},
		},
	})
	if err == nil && p.Allowed {
		return svc.dao.GetUserInfo()
	}
	return errors.New("你没有权限")
}
