package web_test

import (
	"testing"

	"github.com/PhasitWo/duchenne-server/repository"
)

// TODO write test
/* new method when using expecter struct
type-safe method replace .On
RunAndReturn to dynamically set a return value based on the input to the mock's call
*/
func TestCreateAppointment(t *testing.T) {
	repo := repository.MockRepo{}
	repo.EXPECT()
}