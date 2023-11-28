package repos

import (
	"crypto/rsa"
)

var AddKey = MakeAddFunc[*rsa.PrivateKey](Keys)
var SetKey = MakeSetFunc[*rsa.PrivateKey](Keys)
var UpdateKey = MakeUpdateFunc[*rsa.PrivateKey](Keys)
var GetKey = MakeGetFunc[*rsa.PrivateKey](Keys)
