package mail

import (
	"context"

	"github.com/mailgun/mailgun-go/v3"
)

var _ Sender = (*Mailgun)(nil)

// Mailgun implements Sender interface.
type Mailgun struct {
	mg mailgun.Mailgun
}

// NewMailgun creates new Mailgun.
func NewMailgun(mg mailgun.Mailgun) *Mailgun {
	return &Mailgun{
		mg: mg,
	}
}

// Send implements Sender interface.
func (m *Mailgun) Send(ctx context.Context, message *Mail) (*Response, error) {
	msg := m.mg.NewMessage(message.From, message.Subject, message.Text, message.To...)

	res, id, err := m.mg.Send(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &Response{
		Message: res,
		ID:      id,
	}, nil
}
