package internal

type User struct {
	Id	int	`json:"id,omitempty"`
	Name string	`json:"name,omitempty"`
	Password string	`json:"password,omitempty"`
}
