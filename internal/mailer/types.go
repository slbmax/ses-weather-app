package mailer

type ConfirmationEmail struct {
	City      string
	Frequency string
	Token     string
}

type NotificationEmail struct {
	City        string
	Temperature float32
	Description string
	Humidity    uint8
	Frequency   string
}

type ConfirmationSuccessEmail struct {
	City      string
	Frequency string
	Token     string
}
