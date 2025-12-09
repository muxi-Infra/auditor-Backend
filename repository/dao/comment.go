package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/repository/model"
)

type CommentDaoInterface interface {
	UpdateComments(ctx context.Context, id uint, lc *request.Comment, nc *request.Comment) error
}

type CommentDao struct {
	DB *gorm.DB
}

func NewCommentDao(db *gorm.DB) *CommentDao {
	return &CommentDao{
		DB: db,
	}
}

func (d *CommentDao) UpdateComments(ctx context.Context, tid uint, lc *request.Comment, nc *request.Comment) error {
	var comments []model.Comment
	d.DB.WithContext(ctx).Where("item_id = ?", tid).Find(&comments)
	if len(comments) != 2 {
		return errors.New("comments nums error")
	}

	updateComment(&comments[0], lc)
	updateComment(&comments[1], nc)

	for _, comment := range comments {
		err := d.DB.Save(&comment).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func updateComment(m *model.Comment, c *request.Comment) {
	if c.Content != "" {
		m.Content = c.Content
	}

	if len(c.Pictures) > 0 {
		m.Pictures = c.Pictures
	}
}
