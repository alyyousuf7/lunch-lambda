package notification

import (
	"fmt"
)

var _ Notifier = &Console{}

type Console struct{}

func (n *Console) Notify(title, message string) error {
	fmt.Printf("Title: %s\nMessage:\n%s\n", title, message)

	return nil
}
