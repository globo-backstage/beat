package auth;

type User interface {
	Email() string
}

type Authable interface {
	GetUser () User
}
