package repo

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"id-maker/internal/entity"
	"time"
)

type SegmentRepo struct {
	*gorm.DB
}

// New -.
func New(db *gorm.DB) *SegmentRepo {
	return &SegmentRepo{db}
}

// GetList -.
// 从数据库获取数据库列表的实例信息
func (db *SegmentRepo) GetList() ([]entity.Segments, error) {
	var err error
	var s []entity.Segments
	if err = db.Find(&s).Error; err != nil {
		return s, fmt.Errorf("SegmentRepo - GetList - Find: %w", err)
	}
	return s, nil
}

// Add
// 向数据库中添加数据
func (db *SegmentRepo) Add(s *entity.Segments) error {
	var err error
	// 数据库中不存在这条数据才可以执行插入操作
	if s.BizTag == "" {
		return errors.New("nil record")
	}
	if errors.Is(db.Where("biz_tag = ?", s.BizTag).First(s).Error, gorm.ErrRecordNotFound) {
		if err = db.Create(s).Error; err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Tag Already Exist")
}

// GetNextId -.
// ag业务下更新ID号段最大值，并获取下一个可分配的号段区间
func (db *SegmentRepo) GetNextId(tag string) (*entity.Segments, error) {
	var (
		err error
		id  = &entity.Segments{}
		tx  = db.Begin()
	)
	if err = tx.Where("biz_tag = ?", tag).First(id).Error; err != nil {
		tx.Rollback()
		return id, fmt.Errorf("SegmentRepo - GetNextId - Get: %w", err)
	}
	if err = tx.Exec("update segments set max_id=max_id+step, updated_at = ? where biz_tag = ?", time.Now(), tag).Error; err != nil {
		tx.Rollback()
		return id, fmt.Errorf("SegmentRepo - GetNextId - Exec: %w", err)
	}
	tx.Commit()
	return id, nil
}
