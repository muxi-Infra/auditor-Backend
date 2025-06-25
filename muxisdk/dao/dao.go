package dao

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Key [2]any
type Dao interface {
	Get(ctx context.Context, val any, fields ...Key) error
	PreGet(ctx context.Context, val any, pre string, fields ...Key) error
	Create(ctx context.Context, val any) error
	Update(ctx context.Context, val any, fields ...Key) error
	Delete(ctx context.Context, val any, fields ...Key) error
	GetByID(ctx context.Context, val any, id uint) error
	UpdateByID(ctx context.Context, val any, id uint) error
	DeleteByID(ctx context.Context, val any, id uint) error
	DeleteBack(ctx context.Context, val any, id uint) error
	DeleteGet(ctx context.Context, val any) (any, error)
}

var structPtrTypeCache sync.Map

// fields的第一个数据需要为string
func isStructPtr(val any) bool {
	typ := reflect.TypeOf(val)
	if typ == nil {
		return false
	}

	// 尝试从缓存中读取
	if cached, ok := structPtrTypeCache.Load(typ); ok {
		return cached.(bool)
	}

	// 真实判断逻辑
	isPtrStruct := typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct

	// 缓存起来
	structPtrTypeCache.Store(typ, isPtrStruct)

	return isPtrStruct
}
func check(val any) error {
	ok := isStructPtr(val)
	if !ok {
		return errors.New("val 不是指向结构体的指针")
	}
	return nil
}

//以下都是mysql的实现逻辑
//todo 增添其他数据库的支持

func (o *OrmClient) Get(ctx context.Context, val any, fields ...Key) error {
	err := check(val)
	if err != nil {
		return nil
	}
	query := o.DB.WithContext(ctx).Model(val)
	for _, field := range fields {
		key, ok := field[0].(string)
		if !ok {
			return errors.New("field key 必须为字符串")
		}
		query = query.Where(key+" = ?", field[1])
	}
	err = query.First(val).Error
	if err != nil {
		return err
	}
	return nil
}

//需要preload的Get方法

func (o *OrmClient) PreGet(ctx context.Context, val any, pres []string, fields ...Key) error {
	err := check(val)
	if err != nil {
		return err
	}
	query := o.DB.WithContext(ctx).Model(val)
	for _, field := range fields {
		key, ok := field[0].(string)
		if !ok {
			return errors.New("field key 必须为字符串")
		}
		query = query.Where(key+" = ?", field[1])
	}
	for _, pre := range pres {
		if pre != "" {
			query = query.Preload(pre)
		}
	}
	err = query.First(val).Error
	if err != nil {
		return err
	}
	return nil
}
func (o *OrmClient) Create(ctx context.Context, val any) error {
	err := check(val)
	if err != nil {
		return err
	}
	if err = o.DB.WithContext(ctx).Create(val).Error; err != nil {
		return err
	}
	return nil
}
func (o *OrmClient) Update(ctx context.Context, val any, fields ...Key) error {
	err := check(val)
	if err != nil {
		return err
	}
	query := o.DB.WithContext(ctx).Model(val)
	for _, field := range fields {
		key, ok := field[0].(string)
		if !ok {
			return errors.New("field key 必须为字符串")
		}
		query = query.Where(key+" = ?", field[1])

	}
	err = query.Updates(val).Error
	if err != nil {
		return err
	}
	return nil
}
func (o *OrmClient) Delete(ctx context.Context, val any, fields ...Key) error {
	err := check(val)
	if err != nil {
		return err
	}
	query := o.DB.WithContext(ctx).Model(val)
	for _, field := range fields {
		key, ok := field[0].(string)
		if !ok {
			return errors.New("field key 必须为字符串")
		}
		query = query.Where(key+" = ?", field[1])
	}
	err = query.Delete(val).Error
	if err != nil {
		return err
	}
	return nil
}
func (o *OrmClient) GetByID(ctx context.Context, val any, id uint) error {
	err := check(val)
	if err != nil {
		return err
	}
	if err = o.DB.WithContext(ctx).First(val, id).Error; err != nil {
		return err
	}
	return nil
}
func (o *OrmClient) UpdateByID(ctx context.Context, val any, id uint) error {
	err := check(val)
	if err != nil {
		return err
	}
	if err = o.DB.WithContext(ctx).Where("id = ?", id).Updates(val).Error; err != nil {
		return err
	}
	return nil
}
func (o *OrmClient) DeleteByID(ctx context.Context, val any, id uint) error {
	err := check(val)
	if err != nil {
		return err
	}
	if err = o.DB.WithContext(ctx).Where("id = ?", id).Delete(val).Error; err != nil {
		return err
	}
	return nil
}
func (o *OrmClient) DeleteGet(ctx context.Context, val any) (any, error) {
	err := check(val)
	if err != nil {
		return nil, err
	}

	valType := reflect.TypeOf(val)
	if valType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("val must be a pointer to a struct")
	}

	// 获取实际类型（去掉指针）
	elemType := valType.Elem()

	// 创建一个 []elemType 的切片
	sliceType := reflect.SliceOf(elemType)
	slicePtr := reflect.New(sliceType) // 指针类型 *sliceType

	// 执行查询
	err = o.DB.WithContext(ctx).
		Unscoped().
		Where("deleted_at IS NOT NULL").
		Find(slicePtr.Interface()).Error
	if err != nil {
		return nil, err
	}
	return slicePtr.Elem().Interface(), nil
}
func (o *OrmClient) DeleteBack(ctx context.Context, val any, id uint) error {
	err := check(val)
	if err != nil {
		return err
	}
	if err = o.DB.WithContext(ctx).Model(val).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	return nil
}
