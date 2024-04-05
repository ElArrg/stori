package models

const (
	FileSystemSourceType = "file-system"
)

type Channel string

const (
	EmailChannel Channel = Channel("email")
)

type NotificationsSettings struct {
	ID        string
	AccountID string  // AccountID is the reference to the account which is affected by this settings
	Channel   Channel // Channel indicates the notification channel to apply this settings
	Enabled   bool    // Enabled indicates if this channel is enabled for sending notifications to the account
}

type Template struct {
	ID         string
	Operation  string  // Operation is the kind of operation event for the template
	Channel    Channel // Channel is where this template should be sent
	Source     string  // Source is where the actual template content should be found
	SourceType string  // SourceType is the type of repository where the template content is stored
	Active     bool
}
