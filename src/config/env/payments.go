package env

import "os"

var PaymentsVar = map[string]string{
	"RAZORPAY_SECRET_KEY": os.Getenv("RAZORPAY_KEY_SECRET"),
	"RAZORPAY_KEY_ID":     os.Getenv("RAZORPAY_KEY_ID"),
}
