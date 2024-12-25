package env

var PaymentsVar = map[string]string{
	"RAZORPAY_SECRET_KEY": GetEnv("RAZORPAY_KEY_SECRET"),
	"RAZORPAY_KEY_ID":     GetEnv("RAZORPAY_KEY_ID"),
}
