package repository

import (
	"fmt"
	"github.com/PhasitWo/duchenne-server/model"
)

func (r *Repo) GetContent(contentID any) (model.Content, error) {
	var c model.Content
	err := r.db.Where("id = ?", contentID).First(&c).Error
	if err != nil {
		return c, fmt.Errorf("exec : %w", err)
	}
	return c, nil
}

func (r *Repo) GetAllContent(limit int, offset int, criteria ...Criteria) ([]model.Content, error) {
	res := []model.Content{}
	db := attachCriteria(r.db, criteria...)
	// omit body field
	err := db.Model(&model.Content{}).Omit("body").Limit(limit).Offset(offset).Order("`order` ASC").Find(&res).Error
	if err != nil {
		return res, fmt.Errorf("exec : %w", err)
	}
	return res, nil
}

func (r *Repo) CreateContent(content model.Content) (int, error) {
	err := r.db.Create(&content).Error
	if err != nil {
		return -1, fmt.Errorf("exec : %w", err)
	}
	return content.ID, nil
}

func (r *Repo) UpdateContent(content model.Content) error {
	result := r.db.Select("title", "body", "is_published", "order", "cover_image_url").Updates(&content)
	err := result.Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

func (r *Repo) DeleteContent(contentID any) error {
	err := r.db.Where("id = ?", contentID).Delete(&model.Content{}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
