package pkg

import "github.com/Rx-11/EDIS-A1/pkg/users"

var (
	UserRepo users.UserRepo
)

func init() {
	UserRepo = users.NewUserMySQLRepo()
}
