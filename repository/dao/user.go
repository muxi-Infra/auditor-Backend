package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/apikey"
	"muxi_auditor/repository/model"
	"time"
)

type UserDAO struct {
	DB *gorm.DB
}

// NewUserDAO 创建一个新的 UserDAO
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db}
}

func (d *UserDAO) Create(ctx context.Context, user *model.User) error {
	if err := d.DB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserDAO) Read(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := d.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 预计用不上
func (d *UserDAO) Update(ctx context.Context, user *model.User) error {
	if err := d.DB.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

// 预计用不上
func (d *UserDAO) Delete(ctx context.Context, id uint) error {
	if err := d.DB.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserDAO) List(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := d.DB.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (d *UserDAO) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := d.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (d *UserDAO) FindByProjectID(ctx context.Context, id int) ([]model.User, error) {
	var users []model.User
	err := d.DB.WithContext(ctx).Preload("Projects").Joins("JOIN user_projects ON user_projects.user_id = users.id").
		Where("user_projects.project_id = ? AND user_projects.role = ?", id, 1).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (d *UserDAO) FindByUserIDs(ctx context.Context, ids []uint) ([]model.User, error) {
	var users []model.User
	err := d.DB.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (d *UserDAO) GetResponse(ctx context.Context, users []model.User) ([]model.UserResponse, error) {
	var userResponses []model.UserResponse
	for _, user := range users {
		var projectPermits []model.ProjectPermit
		for _, project := range user.Projects {
			var userProject model.UserProject
			d.DB.WithContext(ctx).Where("user_id = ? AND project_id = ?", user.ID, project.ID).First(&userProject)

			projectPermits = append(projectPermits, model.ProjectPermit{
				ProjectID:   project.ID,
				ProjectName: project.ProjectName,
				ProjectRole: userProject.Role,
			})
		}

		userResponses = append(userResponses, model.UserResponse{
			Name:          user.Name,
			UserID:        user.ID,
			Avatar:        user.Avatar,
			ProjectPermit: projectPermits,
			Role:          user.UserRole,
		})
	}

	return userResponses, nil
}
func (d *UserDAO) ChangeProjectRole(ctx context.Context, userId int, projectPermit []model.ProjectPermit, role int) error {
	var user model.User
	err := d.DB.WithContext(ctx).Preload("Projects").Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return errors.New("未找到该用户")
	}
	user.UserRole = role
	var userProject model.UserProject
	for _, project := range projectPermit {
		userProject.Role = project.ProjectRole
		userProject.UserID = user.ID
		userProject.ProjectID = project.ProjectID
		d.DB.WithContext(ctx).Save(&userProject)
	}
	return nil
}
func (d *UserDAO) GetProjectList(ctx context.Context, logo string) ([]model.ProjectList, error) {
	var projects []model.Project
	if err := d.DB.WithContext(ctx).Where("logo = ?", logo).Find(&projects).Error; err != nil {
		return nil, errors.New("查询数据库错误")
	}
	var projectlist []model.ProjectList
	for _, project := range projects {
		projectlist = append(projectlist, model.ProjectList{
			ProjectId:   project.ID,
			ProjectName: project.ProjectName,
		})
	}
	return projectlist, nil
}
func (d *UserDAO) CreateProject(ctx context.Context, project *model.Project) error {
	if err := d.DB.WithContext(ctx).Create(project).Error; err != nil {
		return errors.New("创建项目失败")
	}
	key, err := apikey.GenerateAPIKey(project.ID)
	if err != nil {
		return errors.New("生成apikey失败")
	}
	project.Apikey = key
	if err := d.DB.WithContext(ctx).Save(project).Error; err != nil {
		return err
	}
	return nil

}
func (d *UserDAO) GetProjectDetails(ctx context.Context, id uint) (response.GetDetailResp, error) {
	var project model.Project
	err := d.DB.WithContext(ctx).Preload("Items").Preload("Users").First(&project, id).Error
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
	return re, nil
}
func (d *UserDAO) Select(ctx context.Context, req request.SelectReq) ([]model.Project, error) {
	query := d.DB.WithContext(ctx).Model(&model.Project{})

	if req.ProjectID != 0 {
		query = query.Where("id = ?", req.ProjectID)
	}
	if req.Tag != "" {
		query = query.Where("JSON_CONTAINS(tag, ?)", "\""+req.Tag+"\"")
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Auditor != "" {
		query = query.Where("auditor = ?", req.Auditor)
	}
	if req.Query != "" {
		query = query.Where("project_name LIKE ?", "%"+req.Query+"%")
	}
	var projects []model.Project
	if err := query.Find(&projects).Error; err != nil {
		return nil, errors.New("查询 Project 失败")
	}

	if len(projects) == 0 {
		return nil, nil
	}

	for i := range projects {
		var items []model.Item
		err := d.DB.
			Where("project_id = ?", projects[i].ID).
			Order("created_at DESC").
			Offset((req.Page-1)*req.PageSize).Limit(req.PageSize).
			Preload("Comments", func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at DESC").Limit(2)
			}).
			Find(&items).Error
		if err != nil {
			return nil, errors.New("查询 Items 失败")
		}
		projects[i].Items = items
	}

	return projects, nil
}
func (d *UserDAO) AuditItem(ctx context.Context, req request.AuditReq, id uint) error {
	var item model.Item
	err := d.DB.WithContext(ctx).Where("project_id = ? AND item_id = ?", req.ProjectID, req.ItemId).First(&item).Error
	if err != nil {
		return err
	}
	err = d.DB.WithContext(ctx).
		Model(&model.Item{}).
		Where("project_id = ? AND item_id = ?", req.ProjectID, req.ItemId).
		Updates(map[string]interface{}{
			"status": req.Status,
			"reason": req.Reason,
		}).Error

	if err != nil {
		return err
	}
	var history = model.History{
		UserID: id,
		ItemId: req.ItemId,
	}

	if err := d.DB.WithContext(ctx).Create(history).Error; err != nil {
		return err
	}
	return nil
}
func (d *UserDAO) SearchHistory(ctx context.Context, items *[]model.Item, id uint) error {
	var user model.User
	err := d.DB.WithContext(ctx).Preload("History").Where("id = ?", id).First(&user).Error
	if err != nil {
		return errors.New("未找到用户")
	}
	var itemIds []uint
	for _, h := range user.History {
		itemIds = append(itemIds, h.ItemId)
	}
	err = d.DB.WithContext(ctx).Where("id in ?", itemIds).Order("created_at DESC").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(2)
	}).Find(items).Error
	if err != nil {
		return err
	}
	return nil
}
func (d *UserDAO) Upload(ctx context.Context, req request.UploadReq, id uint, time time.Time) error {
	var project model.Project
	err := d.DB.WithContext(ctx).Where("project_id = ?", id).First(&project).Error
	if err != nil {
		return err
	}
	var item = model.Item{
		Status:     0,
		ProjectId:  project.ID,
		Auditor:    req.Auditor,
		Author:     req.Author,
		Tags:       req.Tags,
		PublicTime: time,
		Content:    req.Content.Topic.Content,
		Title:      req.Content.Topic.Title,
		Pictures:   req.Content.Topic.Pictures,
	}
	err = d.DB.WithContext(ctx).Create(&item).Error
	if err != nil {
		return err
	}
	var comment1 = model.Comment{
		Content:  req.Content.LastComment.Content,
		Pictures: req.Content.LastComment.Pictures,
		ItemId:   item.ID,
	}
	var comment2 = model.Comment{
		Content:  req.Content.NextComment.Content,
		Pictures: req.Content.NextComment.Pictures,
		ItemId:   item.ID,
	}
	err = d.DB.WithContext(ctx).Create(&comment1).Error
	if err != nil {
		return err
	}
	err = d.DB.WithContext(ctx).Create(&comment2).Error
	if err != nil {
		return err
	}
	return nil
}
