package todo

type TodoEntity struct {
	Todo    string `isStructName:"true"`
	Name    string `db:"name"`
	Checked bool   `db:"checked"`
	State   string `db:"state" enum:"created,inReview,done"`
	UserId  int    `db:"user_id" prismaReference:"user"`
}

func CreateTodoEntity() TodoEntity {
	return TodoEntity{}
}
