package repository

import (
	"fmt"
	"time"

	"github.com/PhasitWo/duchenne-server/model"
)

func (r *Repo) GetQuestion(questionId any) (model.SafeQuestion, error) {
	var q model.SafeQuestion
	err := r.db.Model(&model.Question{}).Joins("Doctor").Preload("Patient").Where("questions.id = ?", questionId).First(&q).Error
	if err != nil {
		return q, fmt.Errorf("exec : %w", err)
	}
	return q, nil
}

// Get all questions with following criteria
func (r *Repo) GetAllQuestion(limit int, offset int, criteria ...Criteria) ([]model.QuestionTopic, error) {
	res := []model.QuestionTopic{}
	db := attachCriteria(r.db, criteria...)
	err := db.Model(&model.Question{}).Joins("Doctor").Preload("Patient").Limit(limit).Offset(offset).Order("COALESCE(answer_at, create_at)  DESC").Find(&res).Error
	if err != nil {
		return res, fmt.Errorf("exec : %w", err)
	}
	return res, nil
}

func (r *Repo) CreateQuestion(patientId int, topic string, question string, createAt int) (int, error) {
	q := &model.Question{PatientID: patientId, Topic: topic, Question: question, CreateAt: createAt, DoctorID: nil}
	err := r.db.Create(&q).Error
	if err != nil {
		return -1, fmt.Errorf("exec : %w", err)
	}
	return q.ID, nil
}

func (r *Repo) UpdateQuestionAnswer(questionId int, answer string, doctorId int) error {
	now := int(time.Now().Unix())
	err := r.db.Updates(&model.Question{ID: questionId, Answer: &answer, DoctorID: &doctorId, AnswerAt: &now}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}

func (r *Repo) DeleteQuestion(questionId any) error {
	err := r.db.Where("id = ?", questionId).Delete(&model.Question{}).Error
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
