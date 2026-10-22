// Package mailer 定义邮件发送抽象。
// 低档实现为 SMTP，高档实现为邮件服务商 API。
package mailer

import "context"

// Message 是一封待发送的邮件。
type Message struct {
	To      string
	Subject string
	Text    string
	HTML    string
}

// Mailer 是邮件发送抽象。
type Mailer interface {
	Send(ctx context.Context, m Message) error
}
