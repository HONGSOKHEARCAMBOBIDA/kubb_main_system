package service

import (
	"context"
	"errors"
	"fmt"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"

	"gorm.io/gorm"
)

type RoleService interface {
	CreateRoleHasPermission(ctx context.Context, input request.CreateRolePermissionInput) error
	DeleteRoleHasPermission(ctx context.Context, input request.DeleteRolePermissionsInput) error
	GetRolePermission(ctx context.Context, id int) ([]response.PermissionWithAssignedRole, error)
	UpdateRole(ctx context.Context, id int, input request.RoleRequestUpdate) error
	GetRole(ctx context.Context, id int) ([]model.Role, error)
}

type roleservice struct {
	db *gorm.DB
}

func NewRoleService() RoleService {
	return &roleservice{
		db: config.DB,
	}
}

func (s *roleservice) CreateRoleHasPermission(ctx context.Context, input request.CreateRolePermissionInput) error {
	if len(input.PermissionIDs) == 0 {
		return errors.New("permission_ids cannot be empty")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var rolePermissions []model.RoleHasPermission
	for _, permissionID := range input.PermissionIDs {
		rolePermissions = append(rolePermissions, model.RoleHasPermission{
			RoleID:       uint(input.RoleID),
			PermissionID: uint(permissionID),
		})
	}

	if err := tx.Create(&rolePermissions).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *roleservice) DeleteRoleHasPermission(ctx context.Context, input request.DeleteRolePermissionsInput) error {
	if len(input.PermissionIDs) == 0 {
		return errors.New("permission_ids cannot be empty")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.
		Where("role_id = ? AND permission_id IN ?", input.RoleID, input.PermissionIDs).
		Delete(&model.RoleHasPermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *roleservice) GetRolePermission(ctx context.Context, id int) ([]response.PermissionWithAssignedRole, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid role id: %d", id)
	}

	var permissions []response.PermissionWithAssignedRole
	err := s.db.WithContext(ctx).Table("permissions p").
		Select(`
            p.id AS id,
            p.name AS name,
            CASE 
                WHEN role_has_permission.permission_id IS NULL THEN false 
                ELSE true 
            END AS assigned
        `).
		Joins(`
            LEFT JOIN role_has_permission 
            ON p.id = role_has_permission.permission_id 
            AND role_has_permission.role_id = ?
        `, id).
		Order("p.id ASC").
		Scan(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *roleservice) UpdateRole(ctx context.Context, id int, input request.RoleRequestUpdate) error {
	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.ModuleName != nil {
		updates["module_name"] = *input.ModuleName
	}
	if len(updates) == 0 {
		return errors.New(" no field to update")
	}
	result := s.db.WithContext(ctx).Model(&model.Role{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *roleservice) GetRole(ctx context.Context, id int) ([]model.Role, error) {
	var role []model.Role
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	roleQuery := s.db.WithContext(ctx).Table("role r").
		Select(`
		 r.id AS id,
		 r.name AS name,
		 r.module_name AS module_name
	`)
	roleQuery = helper.ApplyAccessGetRole(roleQuery, s.db, user.Role, user)

	if err := roleQuery.Scan(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}
