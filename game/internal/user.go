package internal

import "fmt"

type User struct {
	Id	int	"_id"
	AccID	string
	Name string	`json:"name,omitempty"`
	Password string	`json:"password,omitempty"`
	CreatedTime	int
	LastLoginTime	int
}

func (data *User)initValue(accID string)error  {
	userID, err := mongoDBNextSeq(C_USERS)
	if err != nil{
		return fmt.Errorf("get next users id error: %v", err)
	}
	data.Id = userID
	data.AccID = accID
	return nil
}


