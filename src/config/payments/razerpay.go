package payments

import (
	"dine-server/src/config/env"

	"github.com/razorpay/razorpay-go"
)

var RazorpayClient *razorpay.Client

func init() {
	// Replace with your Razorpay API Key ID and Secret
	RazorpayClient = razorpay.NewClient(env.PaymentsVar["RAZORPAY_KEY_ID"], env.PaymentsVar["RAZORPAY_SECRET_KEY"])
}
