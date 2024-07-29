package initializer

import (
	"encoding/gob"

	"github.com/razorpay/razorpay-go"
)

var Client *razorpay.Client

func Initialize() {
	Envload()
	DBconnect()
	gob.Register(map[string]interface{}{})

	Client = razorpay.NewClient("RAZOR_PAY_KEY", "RAZOR_PAY_SECRET")
}
