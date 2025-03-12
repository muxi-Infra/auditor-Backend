package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/dao"
	"muxi_auditor/repository/model"
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
func (s *ProjectService) Create(ctx context.Context, name string, logo string, audioRule string, ids []uint) error {

	users, err := s.userDAO.FindByUserIDs(ctx, ids)
	if err != nil {
		return err
	}
	project := model.Project{
		ProjectName: name,
		Logo:        logo,
		AudioRule:   audioRule,
		Users:       users,
	}
	err = s.userDAO.CreateProject(ctx, &project)
	if err != nil {
		return err
	}
	return nil
}
func (s *ProjectService) GetProjectList(ctx context.Context, logo string) ([]model.ProjectList, error) {
	cacheKey := fmt.Sprintf("projectList_%s", logo)
	re, err := s.redisJwtHandler.GetSByKey(ctx, cacheKey)
	if err == nil {
		var list []model.ProjectList
		if err := json.Unmarshal([]byte(re), &list); err == nil {
			return list, nil
		}
	}
	projects, err := s.userDAO.GetProjectList(ctx, logo)
	if err != nil {
		return nil, err
	}
	var list []model.ProjectList
	for _, project := range projects {
		list = append(list, model.ProjectList{
			ProjectId:   project.ID,
			ProjectName: project.ProjectName,
		})
	}
	jsonData, _ := json.Marshal(list)
	s.redisJwtHandler.SetByKey(ctx, cacheKey, jsonData)
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
	var users []model.UserResponse
	for _, user := range project.Users {
		users = append(users, model.UserResponse{
			Name:   user.Name,
			UserID: user.ID,
			Avatar: user.Avatar,
		})
	}

	re := response.GetDetailResp{
		TotleNumber:   countMap[0] + countMap[1] + countMap[2],
		CurrentNumber: countMap[0],
		Apikey:        project.Apikey,
		AuditRule:     project.AudioRule,
		Members:       users,
	}
	jsonData, _ := json.Marshal(re)
	s.redisJwtHandler.SetByKey(ctx, cacheKey, jsonData)
	return re, nil

}
func (s *ProjectService) Delete(ctx context.Context, cla jwt.UserClaims, req request.DeleteProject) error {
	uid := cla.Uid
	role, err := s.userDAO.GetProjectRole(ctx, uid, req.ProjectId)
	if err != nil {
		return err
	}
	if cla.UserRule == 2 || role == 1 {
		err = s.userDAO.DeleteProject(ctx, req.ProjectId)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("无权限")
}
