package wpp

// func getTypeMessage(msg *event.Components) (typ string) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			typ = ""
// 		}
// 	}()

// 	return msg.Entry[0].Changes[0].Value.Messages[0].Type
// }

// func getSatusMessage(msg *event.Components) (status string) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			status = ""
// 		}
// 	}()

// 	return msg.Entry[0].Changes[0].Value.Statuses[0].Status
// }

// sent, delivered, read, failed, deleted, warning
func isVailidStatusMessage(status string) bool {
	if status == "read" || status == "delivered" || status == "sent" || status == "failed" ||
		status == "deleted" || status == "warning" {
		return true
	}

	return false
}
