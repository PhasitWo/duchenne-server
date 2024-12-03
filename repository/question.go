package repository

import (
	"fmt"
	"strconv"

	"github.com/PhasitWo/duchenne-server/model"
)

var questionQuery = `
SELECT
question.id,
topic,
question,
create_at,
answer,
answer_at,
patient_id,
patient.hn,
patient.first_name,
patient.middle_name,
patient.last_name,
patient.email,
patient.phone,
patient.verified,
doctor_id,
doctor.first_name,
doctor.middle_name,
doctor.last_name
FROM question
INNER JOIN patient ON question.patient_id = patient.id 
INNER JOIN doctor ON question.doctor_id = doctor.id
WHERE question.id = ?
`

func (r *Repo) GetQuestion(questionId any) (model.Question, error) {
	var q model.Question
	row := r.db.QueryRow(questionQuery, questionId)
	if err := row.Scan(
		&q.Id,
		&q.Topic,
		&q.Question,
		&q.CreateAt,
		&q.Answer,
		&q.AnswerAt,
		&q.Patient.Id,
		&q.Patient.Hn,
		&q.Patient.FirstName,
		&q.Patient.MiddleName,
		&q.Patient.LastName,
		&q.Patient.Email,
		&q.Patient.Phone,
		&q.Patient.Verified,
		&q.Doctor.Id,
		&q.Doctor.FirstName,
		&q.Doctor.MiddleName,
		&q.Doctor.LastName,
	); err != nil {
		return q, fmt.Errorf("query : %w", err)
	}
	return q, nil
}

var allQuestionQuery = `
SELECT
question.id,
topic,
question,
create_at,
answer,
answer_at,
patient_id,
patient.hn,
patient.first_name,
patient.middle_name,
patient.last_name,
patient.email,
patient.phone,
patient.verified,
doctor_id,
doctor.first_name,
doctor.middle_name,
doctor.last_name
FROM question
INNER JOIN patient ON question.patient_id = patient.id 
INNER JOIN doctor ON question.doctor_id = doctor.id
`

// Get all questions with following criteria
func (r *Repo) GetAllQuestion(id int, criteria QueryCriteria) ([]model.Question, error) {
	var queryString string
	switch criteria {
	case PATIENTID:
		queryString = allQuestionQuery + " " + string(PATIENTID) + strconv.Itoa(id)
	case DOCTORID:
		queryString = allQuestionQuery + " " + string(DOCTORID) + strconv.Itoa(id)
	case NONE:
		queryString = allQuestionQuery
	default:
		return nil, fmt.Errorf("query : invalid criteria")
	}
	rows, err := r.db.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	defer rows.Close()
	res := []model.Question{}
	for rows.Next() {
		var q model.Question
		if err := rows.Scan(
			&q.Id,
			&q.Topic,
			&q.Question,
			&q.CreateAt,
			&q.Answer,
			&q.AnswerAt,
			&q.Patient.Id,
			&q.Patient.Hn,
			&q.Patient.FirstName,
			&q.Patient.MiddleName,
			&q.Patient.LastName,
			&q.Patient.Email,
			&q.Patient.Phone,
			&q.Patient.Verified,
			&q.Doctor.Id,
			&q.Doctor.FirstName,
			&q.Doctor.MiddleName,
			&q.Doctor.LastName,
		); err != nil {
			return nil, fmt.Errorf("query : %w", err)
		}
		res = append(res, q)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	return res, nil
}
