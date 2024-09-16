package user

type UserState struct {
	SomeValue string
}

type UserEntity struct {
	User  string    `isStructName:"true"`
	Email string    `db:"email" prisma:"@unique"`
	State UserState `db:"state"`
}
