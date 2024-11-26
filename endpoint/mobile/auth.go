package mobile

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/PhasitWo/duchenne-server/auth"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/gin-gonic/gin"
)

type login struct {
	Hn        string `json:"hn" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

func (m *mobileHandler) Login(c *gin.Context) {
	var l login
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fetch patient from database
	storedPatient, err := m.repo.GetPatient(l.Hn)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// checking
	if !storedPatient.Verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unverified account"})
		return
	}
	if l.FirstName != storedPatient.FirstName || l.LastName != storedPatient.LastName {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
		return
	}
	// generate token
	token, err := auth.GenerateToken(storedPatient.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

type signup struct {
	Hn         string  `json:"hn" binding:"required"`
	FirstName  string  `json:"firstName" binding:"required"`
	MiddleName *string `json:"middleName" binding:"required"` // nullable
	LastName   string  `json:"lastName" binding:"required"`
	Phone      string  `json:"phone" binding:"required"`
	Email      string  `json:"email" binding:"required"`
}

func (m *mobileHandler) Signup(c *gin.Context) {
	var s signup
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fetch patient from database
	storedPatient, err := m.repo.GetPatient(s.Hn)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// checking
	if storedPatient.Verified { // already verified
		c.JSON(http.StatusConflict, gin.H{"error": "the account have been verified"})
		return
	}
	if s.FirstName != storedPatient.FirstName || s.LastName != storedPatient.LastName {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
		return
	}
	if storedPatient.MiddleName.Valid && *s.MiddleName != storedPatient.MiddleName.String {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
		return
	}
	// update patient info and mark patient as verified
	err = m.repo.UpdatePatient(
		model.Patient{
			Id:         storedPatient.Id,
			Hn:         storedPatient.Hn,
			FirstName:  storedPatient.FirstName,
			MiddleName: storedPatient.MiddleName,
			LastName:   storedPatient.LastName,
			Email:      sql.NullString{String: s.Email, Valid: true},
			Phone:      sql.NullString{String: s.Phone, Valid: true},
			Verified:   true,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// generate token
	token, err := auth.GenerateToken(storedPatient.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})

}

func (m *mobileHandler) Test(c *gin.Context) {
	err := m.repo.UpdatePatient(model.Patient{Id: 3, Hn: "test1", FirstName: "testasdasd", LastName: "ln3 change na", Email: sql.NullString{String: "email change na", Valid: true}, Phone: sql.NullString{String: "000", Valid: true}, Verified: true})
	fmt.Printf("%s", err.Error())
}
