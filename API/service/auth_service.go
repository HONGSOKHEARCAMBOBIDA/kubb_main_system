package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, input request.AuthRequest, c *gin.Context) (*response.AuthResponse, error)
	RefreshToken(ctx context.Context, input request.RefreshTokenRequest, c *gin.Context) (*response.AuthResponse, error)
	Register(ctx context.Context, input request.RegisterRequest) error
	Update(ctx context.Context, id int, input request.UpdateUserRequest) error
	GetUser(ctx context.Context, id int, pf request.Pagination, filter map[string]string) ([]response.UserResponse, *model.PaginationMetadata, error)
}

type authservice struct {
	db *gorm.DB
}

func NewAuthService() AuthService {
	return &authservice{
		db: config.DB,
	}
}

var requiredPermissions = []string{
	"user.view", "user.create", "user.update", "user.delete",
}

func (s *authservice) Login(ctx context.Context, input request.AuthRequest, c *gin.Context) (*response.AuthResponse, error) {

	key := "login_attempt:" + input.Email
	attempts, _ := utils.Redis.Get(utils.Ctx, key).Int()
	if attempts >= 5 {
		return nil, errors.New("អ្នកព្យាយាមចូលច្រើនពេក សូមព្យាយាមម្តងទៀតក្រោយ 10 នាទី")
	}

	var user model.User
	if err := s.db.Select("id,email,password, role_id, is_active, name_kh,name_en").
		Preload("Role").
		Where("email = ? AND is_active = 1", input.Email).
		First(&user).Error; err != nil {

		return nil, errors.New("ព័ត៌មានមិនត្រឹមត្រូវ ឬ អ្នកប្រើប្រាស់ត្រូវបានបិទគណនី")
	}

	var settings []model.Setting
	if err := s.db.Where("`key` IN ?", []string{
		"ACCESS_TOKEN_EXPIRE_HOURS",
		"REFRESH_TOKEN_EXPIRE_DAYS",
	}).Find(&settings).Error; err != nil {
		return nil, errors.New("Setting Not Found")
	}

	settingMap := make(map[string]string)
	for _, s := range settings {
		settingMap[s.Key] = s.Value
	}

	accesstoken, err := strconv.Atoi(settingMap["ACCESS_TOKEN_EXPIRE_HOURS"])
	if err != nil {
		return nil, errors.New("Bad request")
	}

	refreshtoken, err := strconv.Atoi(settingMap["REFRESH_TOKEN_EXPIRE_DAYS"])
	if err != nil {
		return nil, errors.New("Bad request")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	); err != nil {
		utils.Redis.Incr(utils.Ctx, key)
		utils.Redis.Expire(utils.Ctx, key, 10*time.Minute)
		return nil, errors.New("ព័ត៌មានមិនត្រឹមត្រូវ")
	}
	utils.Redis.Del(utils.Ctx, key)

	accessExpiry := time.Now().Add(time.Duration(accesstoken) * time.Hour)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role_id": user.RoleId,
		"exp":     accessExpiry.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenStr, err := accessToken.SignedString(utils.Jwtkey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshTokenStr := hex.EncodeToString(refreshTokenBytes)
	tokenPrefix := refreshTokenStr[:16]
	hashedRefresh := utils.HashToken(refreshTokenStr)
	if err := s.db.Where("user_id = ?", user.ID).Delete(&model.Session{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete session")
	}

	refreshExpiry := time.Now().Add(time.Duration(refreshtoken) * 24 * time.Hour)
	session := model.Session{
		UserID:       uint(user.ID),
		RefreshToken: string(hashedRefresh),
		TokenPrefix:  tokenPrefix,
		ExpiresAt:    refreshExpiry,
	}
	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	var permissions []model.Permission
	if err := s.db.WithContext(ctx).
		Table("permissions p").
		Select("p.name AS name").
		Joins("JOIN role_has_permission rhp ON rhp.permission_id = p.id").
		Where("rhp.role_id = ? AND p.name IN ?", user.RoleId, requiredPermissions).
		Scan(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	resp := &response.AuthResponse{
		Name:         user.NameEN,
		Email:        user.Email,
		RefreshToken: refreshTokenStr,
		AccessToken:  accessTokenStr,
		Role:         []model.Role{user.Role},
		Permissions:  permissions,
	}

	return resp, nil
}

func (s *authservice) RefreshToken(ctx context.Context, input request.RefreshTokenRequest, c *gin.Context) (*response.AuthResponse, error) {
	if len(input.RefreshToken) < 16 {
		return nil, errors.New("Invalid refresh token")
	}
	prefix := input.RefreshToken[:16]

	var session model.Session
	err := s.db.Where("token_prefix = ?", prefix).First(&session).Error
	if err != nil {
		return nil, errors.New("Invalid or expired refresh token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("Invalid or expired refresh token")
	}

	if !utils.VerifyToken(session.RefreshToken, input.RefreshToken) {
		return nil, errors.New("Invalid or expired refresh token")
	}

	var settings []model.Setting
	if err := s.db.Where("`key` IN ?", []string{
		"ACCESS_TOKEN_EXPIRE_HOURS",
		"REFRESH_TOKEN_EXPIRE_DAYS",
	}).Find(&settings).Error; err != nil {
		return nil, err
	}

	settingMap := make(map[string]string)
	for _, s := range settings {
		settingMap[s.Key] = s.Value
	}

	accesstoken, err := strconv.Atoi(settingMap["ACCESS_TOKEN_EXPIRE_HOURS"])
	if err != nil {
		return nil, err
	}

	refreshtoken, err := strconv.Atoi(settingMap["REFRESH_TOKEN_EXPIRE_DAYS"])
	if err != nil {
		return nil, errors.New("Bad request")
	}

	accessExpiry := time.Now().Add(time.Duration(accesstoken) * time.Minute)
	refreshExpiry := time.Now().Add(time.Duration(refreshtoken) * 24 * time.Hour)
	newRefreshBytes := make([]byte, 32)
	if _, err := rand.Read(newRefreshBytes); err != nil {
		return nil, errors.New("failed to generate refresh token")
	}
	newRefreshStr := hex.EncodeToString(newRefreshBytes)
	newHash := utils.HashToken(newRefreshStr)
	newPrefix := newRefreshStr[:16]

	if err := s.db.Model(&session).Updates(model.Session{
		RefreshToken: newHash,
		TokenPrefix:  newPrefix,
		ExpiresAt:    refreshExpiry,
	}).Error; err != nil {
		return nil, err
	}

	var user model.User
	if err := s.db.Select("id,email, role_id, is_active, name_kh,name_en").Preload("Role").Where("id = ?", session.UserID).First(&user).Error; err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleId,
		"exp":     accessExpiry.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(utils.Jwtkey)
	if err != nil {
		return nil, err
	}

	var permissions []model.Permission
	if err := s.db.WithContext(ctx).
		Table("permissions p").
		Select("p.name AS name").
		Joins("JOIN role_has_permission rhp ON rhp.permission_id = p.id").
		Where("rhp.role_id = ? AND p.name IN ?", user.RoleId, requiredPermissions).
		Scan(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	return &response.AuthResponse{
		Name:         user.NameEN,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: newRefreshStr,
		Role:         []model.Role{user.Role},
		Permissions:  permissions,
	}, nil
}

func (s *authservice) Register(ctx context.Context, input request.RegisterRequest) error {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	committed := false
	defer func() {
		if r := recover(); r != nil || !committed {
			tx.Rollback()
		}
	}()

	var role model.Role
	if err := tx.First(&role, input.RoleID).Error; err != nil {
		return err
	}

	var lastUser model.User
	err := tx.Order("id DESC").First(&lastUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	nextID := lastUser.ID + 1

	email := helper.GenerateEmail(input.NameEN, nextID)
	passwordHash := utils.HasPassword("12345678")

	newuser := model.User{
		NameKh:   input.NameKh,
		NameEN:   input.NameEN,
		Email:    email,
		Password: passwordHash,
		Level:    role.Level,
		RoleId:   input.RoleID,
		Gender:   input.Gender,
		Dob:      input.Dob,
		Isactive: true,
	}
	if err := tx.Create(&newuser).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	committed = true
	return nil
}

func applyCommonFilter(query *gorm.DB, filter map[string]string) *gorm.DB {
	for key, value := range filter {
		if value == "" {
			continue
		}
		switch key {
		case "name":
			query = query.Where("u.name_kh LIKE ? OR u.name_en LIKE ?", "%"+value+"%", "%"+value+"%")
		case "role_id":
			query = query.Where("u.role_id =?", value)
		}
	}
	return query
}

func (s *authservice) GetUser(ctx context.Context, id int, pf request.Pagination, filter map[string]string) ([]response.UserResponse, *model.PaginationMetadata, error) {
	var users []response.UserResponse
	var user model.User
	if err := s.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, nil, err
	}
	offset := (pf.Page - 1) * pf.PageSize
	userquery := s.db.WithContext(ctx).Table("user u").
		Select(`
		u.id AS id,
		u.name_kh AS name_kh,
		u.name_en AS name_en,
		u.email AS email,
		r.id AS role_id,
		r.module_name AS role_name,
		u.gender AS gender,
		u.dob AS dob,
		u.is_active AS is_active
	`).
		Joins("LEFT JOIN role r ON r.id = u.role_id")

	userquery = helper.ApplyAccessGetUser(userquery, s.db, user)
	userquery = applyCommonFilter(userquery, filter)
	userquery = userquery.Order("id DESC")
	var totalCount int64
	countQuery := userquery.Session(&gorm.Session{})
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}
	if err := userquery.Offset(offset).Limit(pf.PageSize).Scan(&users).Error; err != nil {
		return nil, nil, err
	}
	for i := range users {
		users[i].Dob = helper.FormatDate(users[i].Dob)
	}
	return users, helper.BuildPaginationMeta(pf, totalCount), nil
}

func (s *authservice) Update(ctx context.Context, id int, input request.UpdateUserRequest) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user not found")
			}
			return err
		}

		updates := make(map[string]interface{})

		if input.NameKh != nil {
			updates["name_kh"] = *input.NameKh
		}

		if input.NameEN != nil {
			updates["name_en"] = *input.NameEN
		}

		if input.Gender != nil {
			updates["gender"] = *input.Gender
		}

		if input.Dob != nil {
			updates["dob"] = *input.Dob
		}

		if input.RoleID != nil {
			var role model.Role
			if err := tx.First(&role, *input.RoleID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("role not found")
				}
				return err
			}
			updates["role_id"] = *input.RoleID
		}

		if len(updates) == 0 {
			return errors.New("no data to update")
		}
		result := tx.Model(&user).Updates(updates)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("user was not updated")
		}

		return nil
	})
}
