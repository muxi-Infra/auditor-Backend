package dao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"muxi_auditor/api/request"
	"muxi_auditor/pkg/apikey"
	"muxi_auditor/repository/model"
	"strings"
	"time"
)

type UserDAOInterface interface {
	Create(ctx context.Context, user *model.User) error
	Read(ctx context.Context, id uint) (*model.User, error)
	Update(ctx context.Context, user *model.User, id uint) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByProjectID(ctx context.Context, id uint) ([]model.User, error)
	FindByUserIDs(ctx context.Context, ids []uint) ([]model.User, error)
	GetResponse(ctx context.Context, users []model.User) ([]model.UserResponse, error)
	PPFUserByid(ctx context.Context, id uint) (model.User, error)
	ChangeProjectRole(ctx context.Context, user model.User, projectPermit []model.ProjectPermit) error
	GetProjectList(ctx context.Context) ([]model.Project, error)
	CreateProject(ctx context.Context, project *model.Project) (uint, error)
	GetProjectDetails(ctx context.Context, id uint) (model.Project, error)
	Select(ctx context.Context, req request.SelectReq) ([]model.Item, error)
	AuditItem(ctx context.Context, ItemId uint, Status int, Reason string, id uint) error
	SelectItemById(ctx context.Context, id uint) (model.Item, error)
	SearchHistory(ctx context.Context, items *[]model.Item, id uint) error
	Upload(ctx context.Context, req request.UploadReq, id uint, time time.Time) error
	GetProjectRole(ctx context.Context, uid uint, pid uint) (int, error)
	DeleteProject(ctx context.Context, pid uint) error
	DeleteUserProject(ctx context.Context, pid uint) error
	RollBack(ItemId uint, Status int, Reason string) error
	UpdateProject(ctx context.Context, id uint, req request.UpdateProject) error
}
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
func (d *UserDAO) Update(ctx context.Context, user *model.User, id uint) error {
	if err := d.DB.WithContext(ctx).Where("id =?", id).Updates(user).Error; err != nil {
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
func (d *UserDAO) DeleteProject(ctx context.Context, pid uint) error {
	err := d.DB.Where("project_id = ?", pid).Delete(&model.Item{}).Error
	if err != nil {
		return err
	}
	if err = d.DB.WithContext(ctx).Where("ID=?", pid).Delete(&model.Project{}).Error; err != nil {
		return err
	}

	return nil
}
func (d *UserDAO) DeleteUserProject(ctx context.Context, pid uint) error {
	if err := d.DB.WithContext(ctx).Where("project_id=?", pid).Delete(&model.UserProject{}).Error; err != nil {
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
func (d *UserDAO) FindByProjectID(ctx context.Context, id uint) ([]model.User, error) {
	var users []model.User
	err := d.DB.WithContext(ctx).Preload("Projects").Joins("JOIN user_projects ON user_projects.user_id = users.id").
		Where("user_projects.project_id = ? ", id).Find(&users).Error
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
func (d *UserDAO) PPFUserByid(ctx context.Context, id uint) (model.User, error) {
	var user model.User
	err := d.DB.WithContext(ctx).Preload("Projects").Where("id = ?", id).First(&user).Error
	if err != nil {
		return model.User{}, errors.New("未找到该用户")
	}
	return user, nil
}

func (d *UserDAO) ChangeProjectRole(ctx context.Context, user model.User, projectPermit []model.ProjectPermit) error {

	var userProject model.UserProject
	for _, project := range projectPermit {
		userProject.Role = project.ProjectRole
		userProject.UserID = user.ID
		userProject.ProjectID = project.ProjectID
		err := d.DB.WithContext(ctx).Save(&userProject).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *UserDAO) GetProjectList(ctx context.Context) ([]model.Project, error) {
	var projects []model.Project
	if err := d.DB.WithContext(ctx).Find(&projects).Error; err != nil {
		return nil, errors.New("查询数据库错误")
	}

	return projects, nil
}
func (d *UserDAO) CreateProject(ctx context.Context, project *model.Project) (uint, error) {
	if err := d.DB.WithContext(ctx).Create(project).Error; err != nil {
		return project.ID, errors.New("创建项目失败")
	}
	key, err := apikey.GenerateAPIKey(project.ID)
	if err != nil {
		return project.ID, errors.New("生成apikey失败")
	}
	project.Apikey = key
	if err := d.DB.WithContext(ctx).Save(project).Error; err != nil {
		return project.ID, err
	}
	return project.ID, nil

}
func (d *UserDAO) GetProjectDetails(ctx context.Context, id uint) (model.Project, error) {
	var project model.Project
	err := d.DB.WithContext(ctx).Preload("Items").Preload("Users").First(&project, id).Error
	if err != nil {
		return model.Project{}, err
	}
	return project, nil

}
func (d *UserDAO) FindProjectByID(ctx context.Context, id uint) (model.Project, error) {
	var project model.Project
	err := d.DB.WithContext(ctx).Where("id = ?", id).First(&project).Error
	if err != nil {
		return model.Project{}, errors.New(fmt.Sprintf("该project: projectid=%d 不存在", id))
	}
	return project, nil
}
func (d *UserDAO) Select(ctx context.Context, req request.SelectReq) ([]model.Item, error) {
	query1 := d.DB.WithContext(ctx).Model(&model.Project{})
	query2 := d.DB.WithContext(ctx).Model(&model.Item{})

	hasFilters := req.ProjectID != 0 || len(req.Tags) > 0 || len(req.Statuses) > 0 ||
		len(req.Auditors) > 0 || len(req.RoundTime) > 0 || req.Query != ""

	if !hasFilters {
		return nil, nil
	}

	if req.ProjectID != 0 {
		query1 = query1.Where("id = ?", req.ProjectID)
		query2 = query2.Where("project_id = ?", req.ProjectID) // 这里补充 project_id 过滤，避免查出所有 items
	}
	if len(req.Tags) > 0 {
		tagConditions := make([]string, 0)
		for _, tag := range req.Tags {
			tagConditions = append(tagConditions, fmt.Sprintf("JSON_CONTAINS(tags, '\"%s\"')", tag))
		}
		query2 = query2.Where(strings.Join(tagConditions, " OR "))
	}
	if len(req.Statuses) > 0 {
		query2 = query2.Where("status IN (?)", req.Statuses)
	}
	if len(req.Auditors) > 0 {
		query2 = query2.Where("auditor IN (?)", req.Auditors)
	}

	// 处理 RoundTime 查询
	if len(req.RoundTime) > 0 {
		var conditions []string
		var values []interface{}

		for _, rt := range req.RoundTime {
			if len(rt) == 2 {
				unixTimestamp1 := int64(rt[0])
				unixTimestamp2 := int64(rt[1])
				if unixTimestamp1 > 1e10 {
					unixTimestamp1 /= 1000
				}
				if unixTimestamp2 > 1e10 {
					unixTimestamp2 /= 1000
				}
				t1 := time.Unix(unixTimestamp1, 0)
				t2 := time.Unix(unixTimestamp2, 0)

				conditions = append(conditions, "(created_at BETWEEN ? AND ?)")
				values = append(values, t1, t2)
			}
		}

		if len(conditions) > 0 {
			queryStr := strings.Join(conditions, " OR ")
			query2 = query2.Where(queryStr, values...)
		}
	}

	// 如果有查询条件
	if req.Query != "" {
		query1 = query1.Where("project_name LIKE ?", "%"+req.Query+"%")
		query2 = query2.Where("title LIKE ?", "%"+req.Query+"%")
	}

	// 查询 Project 并使用 Preload 加载关联的 Items，并加入过滤条件
	var projects []model.Project
	if err := query1.Preload("Items", func(db *gorm.DB) *gorm.DB {
		// 在 Preload 中为 Items 添加查询条件，避免查出不符合条件的数据
		query2 := db.Order("created_at DESC").Limit(2)

		// 将 RoundTime 等条件传递给 Items 查询
		if len(req.RoundTime) > 0 {
			var conditions []string
			var values []interface{}
			for _, rt := range req.RoundTime {
				if len(rt) == 2 {
					unixTimestamp1 := int64(rt[0])
					unixTimestamp2 := int64(rt[1])
					if unixTimestamp1 > 1e10 {
						unixTimestamp1 /= 1000
					}
					if unixTimestamp2 > 1e10 {
						unixTimestamp2 /= 1000
					}
					t1 := time.Unix(unixTimestamp1, 0)
					t2 := time.Unix(unixTimestamp2, 0)

					conditions = append(conditions, "(created_at BETWEEN ? AND ?)")
					values = append(values, t1, t2)
				}
			}
			if len(conditions) > 0 {
				query2 = query2.Where(strings.Join(conditions, " OR "), values...)
			}
		}

		// 如果有其他查询条件，例如 Tags 或 Statuses，也可以在此处添加
		if len(req.Tags) > 0 {
			tagConditions := make([]string, 0)
			for _, tag := range req.Tags {
				tagConditions = append(tagConditions, fmt.Sprintf("JSON_CONTAINS(tags, '\"%s\"')", tag))
			}
			query2 = query2.Where(strings.Join(tagConditions, " OR "))
		}

		if len(req.Statuses) > 0 {
			query2 = query2.Where("status IN (?)", req.Statuses)
		}

		// 返回修改后的 query2
		return query2
	}).Find(&projects).Error; err != nil {
		fmt.Println(err)
		return nil, errors.New("查询 Project 失败")
	}

	var items []model.Item
	if err := query2.Find(&items).Error; err != nil {
		return nil, errors.New("查询 Item 失败")
	}

	// 避免重复数据
	itemMap := make(map[uint]model.Item)
	for _, item := range items {
		itemMap[item.ID] = item
	}
	for _, project := range projects {
		for _, item := range project.Items {
			itemMap[item.ID] = item
		}
	}

	// 转换 map 为 slice
	finalItems := make([]model.Item, 0, len(itemMap))
	for _, item := range itemMap {
		finalItems = append(finalItems, item)
	}

	return finalItems, nil
}

func (d *UserDAO) AuditItem(ctx context.Context, ItemId uint, Status int, Reason string, id uint) error {
	var item model.Item
	err := d.DB.WithContext(ctx).Where(" id = ?", ItemId).First(&item).Error
	if err != nil {
		return err
	}
	err = d.DB.WithContext(ctx).
		Model(&model.Item{}).
		Where(" id = ?", ItemId).
		Updates(map[string]interface{}{
			"status": Status,
			"reason": Reason,
		}).Error

	if err != nil {
		return err
	}
	var history = model.History{
		UserID: id,
		ItemId: ItemId,
	}

	if err := d.DB.WithContext(ctx).Create(&history).Error; err != nil {
		return err
	}

	return nil
}
func (d *UserDAO) RollBack(ItemId uint, Status int, Reason string) error {
	err := d.DB.
		Model(&model.Item{}).
		Where(" id = ?", ItemId).
		Updates(map[string]interface{}{
			"status": Status,
			"reason": Reason,
		}).Error

	if err != nil {
		return err
	}
	return nil
}
func (d *UserDAO) SelectItemById(ctx context.Context, id uint) (model.Item, error) {
	var item model.Item
	err := d.DB.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return model.Item{}, errors.New("获取item失败")
	}
	return item, nil
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
	err := d.DB.WithContext(ctx).Where("id = ?", id).First(&project).Error

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
		HookUrl:    req.HookUrl,
		HookId:     req.Id,
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
func (d *UserDAO) GetProjectRole(ctx context.Context, uid uint, pid uint) (int, error) {
	var project model.UserProject

	err := d.DB.WithContext(ctx).Where("user_id = ? AND project_id = ?", uid, pid).First(&project).Error

	if err != nil {
		return 1, err
	}

	return project.Role, nil
}
func (d *UserDAO) UpdateProject(ctx context.Context, id uint, req request.UpdateProject) error {
	var project model.Project
	err := d.DB.WithContext(ctx).Where("id =?", id).First(&project).Error
	if err != nil {
		return errors.New("project不存在")
	}
	project.AudioRule = req.AudioRule
	project.Logo = req.Logo
	err = d.DB.WithContext(ctx).Save(&project).Error
	if err != nil {
		return errors.New("更新project失败")
	}
	return nil
}
