package todo

// AUTO GENERATED
// DO NOT EDIT

import (
	"fmt"
	"time"
)

type TodoState string

const (
	TodoStateCreated  TodoState = "created"
	TodoStateInReview TodoState = "inReview"
	TodoStateDone     TodoState = "done"
)

type Todo struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Checked   bool      `db:"checked"`
	State     TodoState `db:"state"`
	UserId    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s Todo) String() string {
	return fmt.Sprint("Todo{ ",
		"Id: ", s.Id, ", ",
		"Name: ", s.Name, ", ",
		"Checked: ", s.Checked, ", ",
		"State: ", s.State, ", ",
		"UserId: ", s.UserId, ", ",
		"CreatedAt: ", s.CreatedAt, ", ",
		"UpdatedAt: ", s.UpdatedAt, ", ",
		"}")
}
