package service

import (
	"context"
	merr "errors"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/jwt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/dao"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
	"sync"
)

const maxConcurrency = 10

type UserService struct {
	userDAO         *dao.UserDAO
	redisJwtHandler *jwt.RedisJWTHandler
}

func NewUserService(userDAO *dao.UserDAO, redisJwtHandler *jwt.RedisJWTHandler) *UserService {
	return &UserService{userDAO: userDAO, redisJwtHandler: redisJwtHandler}
}

//好的架构师做乘法，而我一直做加法😅

func (s *UserService) UpdateUserRole(ctx context.Context, userId uint, projectPermit []model.ProjectPermit, role int) error {
	user, err := s.userDAO.PPFUserByid(ctx, userId)
	if err != nil {
		return err
	}
	user.UserRole = role
	err = s.userDAO.Update(ctx, &user, userId)
	if err != nil {
		return err
	}
	for _, v := range projectPermit {
		_, err = s.userDAO.FindProjectByID(ctx, v.ProjectID)
		if err != nil {
			return err
		}
	}
	err = s.userDAO.ChangeProjectRole(ctx, user, projectPermit)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUsersRole 批量更新user role
func (s *UserService) UpdateUsersRole(ctx context.Context, li []request.UserRole) error {

	sem := make(chan struct{}, maxConcurrency)
	var (
		mu   sync.Mutex
		wait sync.WaitGroup
		errs []error
	)
	for _, v := range li {
		wait.Add(1)
		sem <- struct{}{} // 占用一个槽位
		go func(v request.UserRole) {
			defer wait.Done()
			defer func() { <-sem }() // 释放槽位
			user := model.User{UserRole: v.Role}
			mu.Lock()
			err := s.userDAO.Update(ctx, &user, v.Userid)
			if err != nil {
				errs = append(errs, err)
			}
			mu.Unlock()
		}(v)
	}
	wait.Wait()
	if len(errs) > 0 {
		return merr.Join(errs...)
	}
	return nil
}
func (s *UserService) UpdateMyInfo(ctx context.Context, req request.UpdateUserReq, id uint) error {
	if req.Name == "" && req.Avatar == "" {
		return merr.New(" name or avatar are required")
	}

	existingUser, err := s.userDAO.Read(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return merr.New("user not found")
	}
	var user model.User
	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	err = s.userDAO.Update(ctx, &user, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *UserService) GetMyInfo(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userDAO.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, merr.New("user not found")
	}
	return user, nil
}
func (s *UserService) GetUsers(ctx context.Context, req request.GetUsers) ([]response.UserAllInfo, error) {

	users, err := s.userDAO.GetUsers(ctx, req)
	if err != nil {
		return nil, err
	}
	projects, err := s.userDAO.GetProjectList(ctx)
	if err != nil {
		return nil, err
	}
	re, err := s.userDAO.GetUserProjectRoles(ctx, users, projects)
	if err != nil {
		return nil, err
	}
	return re, nil

}
func (s *UserService) GetUserInfo(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userDAO.Read(ctx, id)
	if err != nil {
		return &model.User{}, err
	}
	return user, nil

}
