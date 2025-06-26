package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/jwt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/dao"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
	"strconv"
)

type ProjectService struct {
	userDAO         dao.UserDAOInterface
	redisJwtHandler *jwt.RedisJWTHandler
}
type Count struct {
	AllCount     int
	CurrentCount int
}

func NewProjectService(userDAO dao.UserDAOInterface, redisJwtHandler *jwt.RedisJWTHandler) *ProjectService {
	return &ProjectService{userDAO: userDAO, redisJwtHandler: redisJwtHandler}
}
func (s *ProjectService) Create(ctx context.Context, name string, logo string, audioRule string, ids []uint) (uint, error) {

	users, err := s.userDAO.FindByUserIDs(ctx, ids)
	if err != nil {
		return 1111111, err
	}
	project := model.Project{
		ProjectName: name,
		Logo:        logo,
		AudioRule:   audioRule,
		Users:       users,
	}
	key, err := s.userDAO.CreateProject(ctx, &project)
	if err != nil {
		return key, err
	}
	return key, nil
}
func (s *ProjectService) GetProjectList(ctx context.Context) ([]model.ProjectList, error) {

	projects, err := s.userDAO.GetProjectList(ctx)
	if err != nil {
		return nil, err
	}
	var list []model.ProjectList
	for _, project := range projects {
		list = append(list, model.ProjectList{
			Id:   project.ID,
			Name: project.ProjectName,
		})
	}

	return list, nil
}
func (s *ProjectService) Detail(ctx context.Context, id uint) (response.GetDetailResp, error) {
	cacheKey := fmt.Sprintf("Datil_%s", strconv.Itoa(int(id)))
	r, err := s.redisJwtHandler.GetSByKey(ctx, cacheKey)
	if err == nil {
		var detailResp response.GetDetailResp
		if err := json.Unmarshal([]byte(r), &detailResp); err == nil {
			return detailResp, nil
		}
	}
	project, err := s.userDAO.GetProjectDetails(ctx, id)
	if err != nil {
		return response.GetDetailResp{}, err
	}
	countMap := map[int]int{
		0: 0,
		1: 0,
		2: 0,
	}
	for _, item := range project.Items {
		countMap[item.Status]++
	}
	//var users []model.UserResponse
	//for _, user := range project.Users {
	//	users = append(users, model.UserResponse{
	//		Name:   user.Name,
	//		UserID: user.ID,
	//		Avatar: user.Avatar,
	//	})
	//}

	re := response.GetDetailResp{
		TotalNumber:   countMap[0] + countMap[1] + countMap[2],
		CurrentNumber: countMap[0],
		Apikey:        project.Apikey,
		AuditRule:     project.AudioRule,
	}
	jsonData, _ := json.Marshal(re)
	s.redisJwtHandler.SetByKey(ctx, cacheKey, jsonData)
	return re, nil

}
func (s *ProjectService) Delete(ctx context.Context, cla jwt.UserClaims, projectId uint) error {
	uid := cla.Uid

	if cla.UserRule == 2 {
		err := s.userDAO.DeleteUserProject(ctx, projectId)
		if err != nil {
			return err
		}
		err = s.userDAO.DeleteProject(ctx, projectId)
		if err != nil {
			return err
		}
		return nil
	}
	role, err := s.userDAO.GetProjectRole(ctx, uid, projectId)
	if err != nil {
		return err
	}
	if role == 1 {
		err = s.userDAO.DeleteUserProject(ctx, projectId)
		if err != nil {
			return err
		}
		err = s.userDAO.DeleteProject(ctx, projectId)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("无权限")
}
func (s *ProjectService) Update(ctx context.Context, id uint, req request.UpdateProject) error {
	err := s.userDAO.UpdateProject(ctx, id, req)
	if err != nil {
		return err
	}
	return nil
}
func (s *ProjectService) GetUsers(ctx context.Context, id uint) ([]model.UserResponse, error) {
	users, err := s.userDAO.FindByProjectID(ctx, id)
	if err != nil {
		return nil, err
	}
	userResponse, err := s.userDAO.GetResponse(ctx, users, id)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}
