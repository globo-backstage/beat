package auth;

type DraftAuthentication struct {}
type DraftUser struct {}

func (d *DraftAuthentication) GetUser() User  {
	return &DraftUser{};
}

func (u *DraftUser) Email () string  {
	return "anonymous@example.net"
}

