package notification

// Notifier sends notification
type Notifier interface {
	Notify(title, message string) error
}
