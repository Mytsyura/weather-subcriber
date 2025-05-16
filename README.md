# Weather Subscription Service

A service that allows users to subscribe to weather updates for specific cities. Users can receive weather updates either daily or hourly, depending on their preference.

## Features

- User subscription management (create, confirm, unsubscribe)
- Weather updates via email
- Configurable update frequency (daily/hourly)
- Local email testing with Mailhog
- PostgreSQL database for data persistence
- RESTful API endpoints

## Prerequisites

- Go 1.23 or higher
- PostgreSQL
- Mailhog (for local email testing)
- Weather API key (from weatherapi.com)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Mytsyura/weather-subcriber.git
cd weather-subcriber
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL:
```bash
# Create database
createdb weather_subscription

# Run migrations to create tables and set up permissions
psql -U postgres -d weather_subscription -f internal/db/migrations/001_create_subscriptions_table.sql
```

4. Install Mailhog (for local email testing):
```bash
# Using Homebrew (macOS)
brew install mailhog

# Start Mailhog
mailhog
```

## Configuration

1. Update `config/config.yaml` with your specific details:
```yaml
dbAddr: "localhost:5432"
dbName: "weather_subscription"
dbUser: "postgres"
dbPass: "postgres"
serverPort: "8080" 
weather_api:
  key: {set_up_your_key} // register an account to get a key here: https://www.weatherapi.com/
email:
  # Local development settings (used when ENV=local)
  local:
    from: "noreply@weather-subscription.com"
    smtpHost: "localhost"
    smtpPort: "1025"
  # Production settings (used when ENV=production)
  production:
    from: "${EMAIL_FROM}"
    password: "${EMAIL_PASSWORD}"
    smtpHost: "${EMAIL_SMTP_HOST}"
    smtpPort: "${EMAIL_SMTP_PORT}"
```

2. Set environment variables:
```bash
export ENV=local
export WEATHER_API_KEY=your-weather-api-key
```

## Running the Application

1. Start the application:
```bash
go run main.go
```

The server will start on port 8080 (or the port specified in your config).

## Using the Application

1. Open your web browser and navigate to:
```
http://localhost:8080
```

2. You'll see a user-friendly interface with two tabs:
   - **Subscribe**: Enter your email, select a city, and choose your preferred update frequency (daily/hourly)
   - **Unsubscribe**: Enter your email to cancel your weather subscription

3. After subscribing:
   - You'll receive a confirmation email
   - Click the confirmation link in the email to activate your subscription
   - Start receiving weather updates according to your chosen frequency

## Email Testing

When running in local environment (ENV=local):
- Emails are sent to Mailhog
- View emails at http://localhost:8025
- No SMTP authentication required
