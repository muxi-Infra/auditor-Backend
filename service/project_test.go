package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/dao"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
	"testing"
	"time"
)

type RedisJWTHandler struct {
	cmd *redis.Client
}
type projectService struct {
	userDAO         *dao.UserDAO
	RedisJWTHandler *RedisJWTHandler
}

func (r *RedisJWTHandler) GetSByKey(ctx context.Context, cacheKey string) (string, error) {
	re, err := r.cmd.Get(ctx, cacheKey).Result()
	if err != nil {
		return "", err
	}
	return re, nil
}
func (r *RedisJWTHandler) SetByKey(ctx context.Context, cacheKey string, list []byte) error {

	err := r.cmd.Set(ctx, cacheKey, list, time.Hour).Err()
	return err
}

func setupProjectService() *projectService {

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // 你的 Redis 地址
	})

	redisHandler := &RedisJWTHandler{cmd: redisClient}

	db, err := gorm.Open(mysql.Open("root:chenhaoqi318912@tcp(127.0.0.1:3306)/muxiAuditor?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	userDAO := &dao.UserDAO{DB: db}

	return &projectService{
		userDAO:         userDAO,
		RedisJWTHandler: redisHandler,
	}
}

//func init(db *gorm.DB) {
//	// 自动迁移：确保数据库表与模型结构一致
//	if err := db.AutoMigrate(&model.User{}, &model.Project{}, &model.Item{}); err != nil {
//		log.Fatalf("Failed to migrate database: %v", err)
//	}
//
//	// 插入用户数据
//	users := []model.User{
//		{Name: "User 1", Avatar: "logo.png"},
//		{Name: "User 2", Avatar: "logo.png"},
//		{Name: "User 3", Avatar: "logo.png"},
//	}
//	for i := range users {
//		if err := db.Where("name = ?", users[i].Name).FirstOrCreate(&users[i]).Error; err != nil {
//			log.Fatalf("Failed to insert user data: %v", err)
//		}
//	}
//
//	// 查询已有用户
//	var user1, user2, user3 model.User
//	db.Where("name = ?", "User 1").First(&user1)
//	db.Where("name = ?", "User 2").First(&user2)
//	db.Where("name = ?", "User 3").First(&user3)
//
//	// 插入项目数据
//	projects := []model.Project{
//		{Model: gorm.Model{ID: 1},
//			ProjectName: "Project A",
//			Logo:        "logo.png",
//			AuditRule:   "Audio Rule A",
//			Users:       []model.User{user1, user2, user3}, // 直接使用已有的用户
//			Items: []model.Item{
//				{Author: "Item 1", Status: 1, PublicTime: time.Now()},
//				{Auditor: "Item 2", Status: 0, PublicTime: time.Now()},
//			},
//			Apikey: "apikey123",
//		}, {
//			Model:       gorm.Model{ID: 2},
//			ProjectName: "Project B",
//			Logo:        "logo.png",
//			AuditRule:   "Audio Rule B",
//			Users: []model.User{
//				{Name: "User 4", Avatar: "logo.png"},
//			},
//			Items: []model.Item{
//				{Author: "Item 3", Status: 1, PublicTime: time.Now()},
//			},
//		},
//	}
//
//	for _, project := range projects {
//		if err := db.Create(&project).Error; err != nil {
//			log.Fatalf("Failed to insert project data: %v", err)
//		}
//	}
//}

func (s *projectService) DatilWithRedis(ctx context.Context) (response.GetDetailResp, error) {
	id := s.Find(ctx)
	cacheKey := fmt.Sprintf("Datil_%s", strconv.Itoa(int(id)))
	r, err := s.RedisJWTHandler.GetSByKey(ctx, cacheKey)
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
		AuditRule:     project.AuditRule,
		Members:       users,
	}
	jsonData, _ := json.Marshal(re)
	s.RedisJWTHandler.SetByKey(ctx, cacheKey, jsonData)
	return re, nil

}
func (s *projectService) Find(ctx context.Context) uint {
	var project model.Project
	if err := s.userDAO.DB.Where("project_name = ?", "Project A").First(&project).Error; err != nil {
		return 0
	}
	return project.ID
}
func (s *projectService) DatilNoRedis(ctx context.Context) (response.GetDetailResp, error) {
	id := s.Find(ctx)
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
		AuditRule:     project.AuditRule,
		Members:       users,
	}
	return re, nil
}

func BenchmarkDatilWithRedis(b *testing.B) {
	s := setupProjectService()

	for i := 0; i < b.N; i++ {
		_, err := s.DatilWithRedis(context.Background())
		if err != nil {
			log.Fatalf("Failed to insert user data: %v", err)
		}
	}
}

func BenchmarkDatilService(b *testing.B) {
	s := setupProjectService()

	for i := 0; i < b.N; i++ {
		_, err := s.DatilNoRedis(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}
}
