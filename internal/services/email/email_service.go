package email

import (
	"context"
	"fmt"
	"net/smtp"
	"os"

	models "weather_subscription/internal/db/models"
)

type EmailService struct {
	from     string
	password string
	smtpHost string
	smtpPort string
	isLocal  bool
}

func NewEmailService(from, password, smtpHost, smtpPort string) *EmailService {
	// Check if we're in local environment
	isLocal := os.Getenv("ENV") == "local" || os.Getenv("ENV") == "development"

	return &EmailService{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		isLocal:  isLocal,
	}
}

func (s *EmailService) SendConfirmationEmail(to, token string) error {
	subject := "Confirm Your Weather Subscription"
	body := fmt.Sprintf(`
		Hello!

		Thank you for subscribing to our weather service. To confirm your subscription, please click the link below:

		http://localhost:8080/api/confirm/%s

		If you did not request this subscription, please ignore this email.

		Best regards,
		Weather Subscription Team
	`, token)

	err := s.sendEmail(to, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *EmailService) SendWeatherUpdate(ctx context.Context, email string, forecast *models.WeatherForecast) error {
	subject := fmt.Sprintf("Weather Update for %s", forecast.City)
	body := fmt.Sprintf(`
		<h2>Weather Update for %s</h2>
		<p>Current weather conditions:</p>
		<ul>
			<li>Temperature: %.1f°C</li>
			<li>Conditions: %s</li>
			<li>Humidity: %d%%</li>
			<li>Wind Speed: %.1f km/h</li>
		</ul>
		<p>Stay dry and have a great day!</p>
	`, forecast.City, forecast.Temperature, forecast.Description, forecast.Humidity, forecast.WindSpeed)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	var err error
	if s.isLocal {
		err = smtp.SendMail(addr, nil, s.from, []string{to}, []byte(message))
	} else {
		auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)
		err = smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message))
	}

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
