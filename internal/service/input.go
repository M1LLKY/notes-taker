package service

type SignIn struct {
	Username string
	Password string
}

type CreateNote struct {
	Title   string
	Content string
	UserID  int
}

type UpdateNote struct {
	ID      int
	Title   string
	Content string
	UserID  int
}
