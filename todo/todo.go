package todo

type TodoState string

const (
	TodoStateCreated  TodoState = "created"
	TodoStateInReview TodoState = "inReview"
	TodoStateDone     TodoState = "done"
)

type Todo struct {
	Id      int       `db:"id"`
	Name    string    `db:"name"`
	Checked bool      `db:"checked"`
	State   TodoState `db:"state" enum:"created,inReview,done"`
}
