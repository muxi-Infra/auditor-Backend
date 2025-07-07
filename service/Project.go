package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/apikey"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/jwt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/dao"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
	"log"
	"net/http"
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

//这里的逻辑有点神秘了，但已经写成这样了，懒得改了，目前大概是有两个鉴权机制，一个是用来获取project_id,区分project的
//另一个是access_key机制，就和七牛云一样，这个是来确认调用方身份的。

func (s *ProjectService) Create(ctx context.Context, name string, url string, logo string, audioRule string, ids []uint) (uint, error) {

	users, err := s.userDAO.FindByUserIDs(ctx, ids)
	if err != nil {
		return 0, err
	}
	ac, se := apikey.GenerateKeyPair()

	project := model.Project{
		ProjectName: name,
		Logo:        logo,
		AudioRule:   audioRule,
		Users:       users,
		AccessKey:   ac,
		SecretKey:   se,
		HookUrl:     url,
	}
	key, err := s.userDAO.CreateProject(ctx, &project)
	if err != nil {

		return key, err
	}
	go func() {
		if err := s.ReturnApiKey("", url); err != nil {
			log.Println(err)
		}
	}()
	return key, nil
}

//给调用方指定接口发送密钥

func (s *ProjectService) ReturnSecretKey(ac string, se string, to string) error {
	var b = request.ReturnSecret{
		SecretKey: se,
		AccessKey: ac,
		Message:   "私钥只生成一次，请妥善保管，如遗失请重置",
	}
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, to, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil

}
func (s *ProjectService) ReturnApiKey(apiKey string, hookUrl string) error {
	var b = request.ReturnApiKey{
		ApiKey:  apiKey,
		Message: "私钥只生成一次，请妥善保管，如遗失请联系管理员重置",
	}
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, hookUrl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
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

// GetAllTags 获取某个项目中所有的Tags
func (s *ProjectService) GetAllTags(ctx context.Context, pid uint) error {

}
