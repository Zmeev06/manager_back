package repos

import (
	"stupidauth/models"
)

var UpdateUser = MakeUpdateFunc[models.User](Users)
var GetUser = MakeGetFunc[models.User](Users)
var AddUser = MakeAddFunc[models.User](Users)
