package user

// AUTO GENERATED
// DO NOT EDIT

import (
	"fmt"
	"time"
)

type User struct {
	Id        int       `db:"id"`
	Email     string    `db:"email"`
	State     UserState `db:"state"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s User) String() string {
	return fmt.Sprint("User{ ",
		"Id: ", s.Id, ", ",
		"Email: ", s.Email, ", ",
		"State: ", s.State, ", ",
		"CreatedAt: ", s.CreatedAt, ", ",
		"UpdatedAt: ", s.UpdatedAt, ", ",
		"}")
}
