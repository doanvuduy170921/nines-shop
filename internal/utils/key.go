package utils

import "fmt"

func GenRedisKeyTemptForPendingOrder(pOrderID int32) string {
	return fmt.Sprintf("otp_tempt:order:%d", pOrderID)
}
