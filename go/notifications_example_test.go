package forge

func ExampleNotificationListOptions_String() {
	_ = (*NotificationListOptions).String
}

func ExampleNotificationListOptions_GoString() {
	_ = (*NotificationListOptions).GoString
}

func ExampleNotificationRepoMarkOptions_String() {
	_ = (*NotificationRepoMarkOptions).String
}

func ExampleNotificationRepoMarkOptions_GoString() {
	_ = (*NotificationRepoMarkOptions).GoString
}

func ExampleNotificationMarkOptions_String() {
	_ = (*NotificationMarkOptions).String
}

func ExampleNotificationMarkOptions_GoString() {
	_ = (*NotificationMarkOptions).GoString
}

func ExampleNotificationService_List() {
	_ = (*NotificationService).List
}

func ExampleNotificationService_Iter() {
	_ = (*NotificationService).Iter
}

func ExampleNotificationService_NewAvailable() {
	_ = (*NotificationService).NewAvailable
}

func ExampleNotificationService_ListRepo() {
	_ = (*NotificationService).ListRepo
}

func ExampleNotificationService_IterRepo() {
	_ = (*NotificationService).IterRepo
}

func ExampleNotificationService_MarkNotifications() {
	_ = (*NotificationService).MarkNotifications
}

func ExampleNotificationService_MarkRepoNotifications() {
	_ = (*NotificationService).MarkRepoNotifications
}

func ExampleNotificationService_MarkRead() {
	_ = (*NotificationService).MarkRead
}

func ExampleNotificationService_GetThread() {
	_ = (*NotificationService).GetThread
}

func ExampleNotificationService_MarkThreadRead() {
	_ = (*NotificationService).MarkThreadRead
}
