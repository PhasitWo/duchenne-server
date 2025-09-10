package repository

import (
	"fmt"

	"github.com/PhasitWo/duchenne-server/model"
	"gorm.io/gorm/clause"
)

func (r *Repo) GetConsentById(consentId any) (model.Consent, error) {
	var c model.Consent
	err := r.db.Where("id = ?", consentId).First(&c).Error
	if err != nil {
		return c, fmt.Errorf("exec : %w", err)
	}
	return c, nil
}

func (r *Repo) GetConsentBySlug(slug string) (model.Consent, error) {
	var c model.Consent
	err := r.db.Where("slug = ?", slug).First(&c).Error
	if err != nil {
		return c, fmt.Errorf("exec : %w", err)
	}
	return c, nil
}

func (r *Repo) UpsertConsent(consent model.Consent) (string, error) {
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "slug"}},
		DoUpdates: clause.AssignmentColumns([]string{"body"}),
	}).Create(&consent).Error
	if err != nil {
		return "", fmt.Errorf("exec : %w", err)
	}
	return consent.Slug, nil
}

func (r *Repo) DeleteConsentById(consentID any) error {
	err := r.db.Where("id = ?", consentID).Delete(&model.Consent{}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

func (r *Repo) DeleteConsentBySlug(slug string) error {
	err := r.db.Where("slug = ?", slug).Delete(&model.Consent{}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
