package repository

import (
	"fmt"
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
LEFT JOIN doctor ON question.doctor_id = doctor.id
WHERE question.id = ?
`

type interDoctor struct {
	doctorId         *int
	doctorFirstName  *string
	doctorMiddleName *string
	doctorLastName   *string
}

func (r *Repo) GetQuestion(questionId any) (model.Question, error) {
	var q model.Question
	var i interDoctor
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
		&i.doctorId,
		&i.doctorFirstName,
		&i.doctorMiddleName,
		&i.doctorLastName,
	); err != nil {
		return q, fmt.Errorf("query : %w", err)
	}
	if i.doctorId != nil {
		q.Doctor = &model.TrimDoctor{}
		q.Doctor.Id = *i.doctorId
		q.Doctor.FirstName = *i.doctorFirstName
		q.Doctor.MiddleName = i.doctorMiddleName
		q.Doctor.LastName = *i.doctorLastName
	}
	return q, nil
}

var allQuestionQuery = `
SELECT
question.id,
topic,
create_at,
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
LEFT JOIN doctor ON question.doctor_id = doctor.id
`

// Get all questions with following criteria
func (r *Repo) GetAllQuestion(criteria ...Criteria) ([]model.QuestionTopic, error) {
	queryString := attachCriteria(allQuestionQuery, criteria...)
	rows, err := r.db.Query(queryString + " ORDER BY CASE WHEN answer_at IS NOT NULL THEN answer_at ELSE create_at END DESC")
	if err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	defer rows.Close()
	res := []model.QuestionTopic{}
	for rows.Next() {
		var q model.QuestionTopic
		var i interDoctor
		if err := rows.Scan(
			&q.Id,
			&q.Topic,
			&q.CreateAt,
			&q.AnswerAt,
			&q.Patient.Id,
			&q.Patient.Hn,
			&q.Patient.FirstName,
			&q.Patient.MiddleName,
			&q.Patient.LastName,
			&q.Patient.Email,
			&q.Patient.Phone,
			&q.Patient.Verified,
			&i.doctorId,
			&i.doctorFirstName,
			&i.doctorMiddleName,
			&i.doctorLastName,
		); err != nil {
			return nil, fmt.Errorf("query : %w", err)
		}
		if i.doctorId != nil {
			q.Doctor = &model.TrimDoctor{}
			q.Doctor.Id = *i.doctorId
			q.Doctor.FirstName = *i.doctorFirstName
			q.Doctor.MiddleName = i.doctorMiddleName
			q.Doctor.LastName = *i.doctorLastName
		}
		res = append(res, q)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query : %w", err)
	}
	return res, nil
}

// var allQuestionTopicQuery = `
// SELECT
// question.id,
// topic,
// create_at,
// answer_at
// FROM question
// `

// // Get all questions with following criteria
// func (r *Repo) GetAllQuestionTopic(id int, criteria QueryCriteria) ([]model.QuestionTopic, error) {
// 	var queryString string
// 	switch criteria {
// 	case PATIENTID:
// 		queryString = allQuestionTopicQuery + " " + string(PATIENTID) + strconv.Itoa(id)
// 	case DOCTORID:
// 		queryString = allQuestionTopicQuery + " " + string(DOCTORID) + strconv.Itoa(id)
// 	case NONE:
// 		queryString = allQuestionTopicQuery
// 	default:
// 		return nil, fmt.Errorf("query : invalid criteria")
// 	}
// 	rows, err := r.db.Query(queryString + " ORDER BY create_at DESC")
// 	if err != nil {
// 		return nil, fmt.Errorf("query : %w", err)
// 	}
// 	defer rows.Close()
// 	res := []model.QuestionTopic{}
// 	for rows.Next() {
// 		var q model.QuestionTopic
// 		if err := rows.Scan(
// 			&q.Id,
// 			&q.Topic,
// 			&q.CreateAt,
// 			&q.AnswerAt,
// 		); err != nil {
// 			return nil, fmt.Errorf("query : %w", err)
// 		}
// 		res = append(res, q)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("query : %w", err)
// 	}
// 	return res, nil
// }

var createQuestionQuery = `
INSERT INTO question (patient_id, topic, question, create_at)
VALUES (?, ?, ?, ?)
`

func (r *Repo) CreateQuestion(patientId int, topic string, question string, createAt int) (int, error) {
	result, err := r.db.Exec(createQuestionQuery, patientId, topic, question, createAt)
	if err != nil {
		return -1, fmt.Errorf("exec : %w", err)
	}
	i, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("exec : %w", err)
	}
	lastId := int(i)
	return lastId, nil
}

var deleteQuestionQuery = `
DELETE FROM question
WHERE id = ?;
`

func (r *Repo) DeleteQuestion(questionId any) error {
	_, err := r.db.Exec(deleteQuestionQuery, questionId)
	if err != nil {
		return fmt.Errorf("exec : %w", err)
	}
	return nil
}
