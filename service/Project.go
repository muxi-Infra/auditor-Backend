package service

import (
	"context"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/dao"
	"muxi_auditor/repository/model"
)

type ProjectService struct {
	userDAO         *dao.UserDAO
	redisJwtHandler *jwt.RedisJWTHandler
}
type Count struct {
	AllCount     int
	CurrentCount int
}

func NewProjectService(userDAO *dao.UserDAO, redisJwtHandler *jwt.RedisJWTHandler) *ProjectService {
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
	return list, nil
}
func (s *ProjectService) Detail(ctx context.Context, id uint) (response.GetDetailResp, error) {
	detail, err := s.userDAO.GetProjectDetails(ctx, id)
	if err != nil {
		return response.GetDetailResp{}, err
	}

	return detail, nil
}
