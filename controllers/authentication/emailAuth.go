package authentication

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/smtp"
	"os"
	"time"

	dbFunc "business-connect/database/dbHelpFunc"
	Data "business-connect/models"

	"github.com/joho/godotenv"
	"github.com/jordan-wright/email"
)

const (
	// smtpAuthAddress   = "smtp.gmail.com"
	// smtpServerAddress = "smtp.gmail.com:587"

	smtpAuthAddress   = "smtp.zoho.com"
	smtpServerAddress = "smtp.zoho.com:587"

	// smtpAuthAddress   = "smtppro.zoho.com"
	// smtpServerAddress = "smtppro.zoho.com:587"
)

var (
	emailSendErr error
	randError    error
	envErr       error
	otp          string
)

type EmailSender interface {
	SendEmail(subject, content string, to, cc, bcc []string, attachFiles []string) error
}

type EmailConfig struct {
	Name              string
	FromEmailAddress  string
	FromEmailPassword string
}

type GmailSender struct {
	Config EmailConfig
}

func NewGmailSender(config EmailConfig) EmailSender {
	return &GmailSender{
		Config: config,
	}
}

func (sender *GmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.Config.Name, sender.Config.FromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s : %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.Config.FromEmailAddress, sender.Config.FromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return value
}

func EmailVerification(name, sendTo string) error {
	var (
		newOTP Data.OTP
	)

	if os.Getenv("RENDER") == "" {
		// Local development only
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Failed to load .env file: %v\n", err)
		}
	}

	config := EmailConfig{
		Name:              mustEnv("EMAIL_SENDER_NAME"),
		FromEmailAddress:  mustEnv("EMAIL_SENDER_ACCOUNT"),
		FromEmailPassword: mustEnv("EMAIL_SENDER_PASSWORD"),
	}

	otp, randError = EmailOTPGeneratorNumber(6)
	if randError != nil {
		return randError
	}

	sender := NewGmailSender(config)
	subject := "Business Connect Email Verification"
	// HTML template with %s as a placeholder for the OTP
	htmlTemplate := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Business Connect Email Verification</title>
	</head>
	<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; text-align: center;">
		<div style="max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); text-align: left;">
			<img src="https://payuee.shop/assets/images/logo.png" alt="Business Connect Logo" style="display: block; margin: 0 auto; max-width: 100%;">
			<h1 style="color: #333333; margin-bottom: 20px; text-align: center;">Verify Your Email Address</h1>

			<p style="color: #777777;">Hi {{.Name}},</p>
			<p style="color: #777777;">We're excited to have you on board! Just one more step to activate your Business Connect account:</p>
			<p style="color: #777777;">Please click the button below to verify your email address and start exploring Business Connect's features.</p>

			<a href="#" style="display: block; margin: 0 auto; padding: 15px; background-color: #007bff; color: #ffffff; text-decoration: none; font-size: 24px; border-radius: 6px; width: 200px; text-align: center;">{{.Token}}</a>

			<p style="color: #777777; margin-top: 20px;">For security reasons, this link will expire in 24 hours. If you have any issues, please contact us at <a href="mailto:support@businessconnectt.com">support@businessconnectt.com</a>.</p>
			<p style="color: #777777;">Thanks,<br>The Business Connect Team</p>
		</div>
	</body>
	</html>
    `
	// Sample data
	data := struct {
		Name string
		Token  string
	}{
		Name: name,
		Token:  otp,
	}

	// Create a new template and parse the HTML
	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		fmt.Println(err)
	}

	// Execute the template with the provided data
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		fmt.Println(err)
	}

	// Convert the buffer to a string to get the final HTML content
	content := body.String()

	// Format the HTML template with the actual OTP
	// url := "https://shopsphereafrica.com/?user=" + sendTo +"&token=" + otp
	// content := fmt.Sprintf(htmlTemplate, name, url)

	to := []string{sendTo}

	emailSendErr = sender.SendEmail(subject, content, to, nil, nil, nil)
	if emailSendErr != nil {
		fmt.Println(emailSendErr)
		return fmt.Errorf("failed to send email: %w", emailSendErr)
	}

	// let's check if the user has a previous otp stored
	oldOTP, getErr := dbFunc.DBHelper.GetOTPByEmail(sendTo)

	if getErr != nil {
		if getErr.Error() == "otp not found by email" {
			var usersOTP Data.OTP
			usersOTP.Email = sendTo
			usersOTP.OTP = otp
			usersOTP.CreatedAT = time.Now().Add(60 * time.Minute).Unix()
			usersOTP.EmailVerification = true
			usersOTP.PasswordReset = false

			saveErr := dbFunc.DBHelper.CreateOTP(usersOTP)
			if saveErr != nil {
				fmt.Println(saveErr)
				return fmt.Errorf("failed to save email otp to db: %w", saveErr)
			}
			return nil
		} else {
			return fmt.Errorf("failed to retrieve email otp from db: %w", getErr)
		}
	}

	// if user exist let's update the users otp to a new one
	newOTP.OTP = otp
	newOTP.CreatedAT = time.Now().Add(60 * time.Minute).Unix()
	newOTP.MaxTry = 0
	log.Println("new otp to update: ", newOTP)
	updateErr := dbFunc.DBHelper.UpdateExistingOTP(newOTP, oldOTP.CustomID)
	if updateErr != nil {
		return fmt.Errorf("failed to save email otp to db: %w", updateErr)
	}

	return nil
}

func ForgotPasswordEmailVerification(name, sendTo string) error {
	var (
		newOTP Data.OTP
	)

	if os.Getenv("RENDER") == "" {
		// Local development only
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Failed to load .env file: %v\n", err)
		}
	}

	config := EmailConfig{
		Name:              mustEnv("EMAIL_SENDER_NAME"),
		FromEmailAddress:  mustEnv("EMAIL_SENDER_ACCOUNT"),
		FromEmailPassword: mustEnv("EMAIL_SENDER_PASSWORD"),
	}

	otp, randError = EmailOTPGenerator(30)
	if randError != nil {
		return randError
	}

	sender := NewGmailSender(config)
	subject := "Business Connect Admin Password Reset"
	// HTML template with %s as a placeholder for the OTP
	htmlTemplate := `
	<!DOCTYPE html>
	<html lang="en">

	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta name="description" content="Business Connect admin is super flexible, powerful, clean &amp; modern responsive bootstrap 5 admin template with unlimited possibilities.">
		<meta name="keywords" content="admin template, Business Connect admin template, dashboard template, flat admin template, responsive admin template, web app">
		<meta name="author" content="payuee">
		<link rel="icon" href="https://shopsphereafrica.com/image/catalog/cart.png" type="image/x-icon">
		<link rel="shortcut icon" href="https://shopsphereafrica.com/image/catalog/cart.png" type="image/x-icon">
		<title>Reset Your Business Connect - Admin Password</title>
	</head>
	<body style="margin: 30px auto; width: 650px; font-family: Work Sans, sans-serif; background-color: #f6f7fb; display: block; padding: 0 12px;">
		<div style="max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); text-align: left;">
		<table style="width: 100%;">
			<tbody>
				<tr>
				<td>
					<table style="background-color: #f6f7fb; width: 100%;">
					<tbody>
						<tr>
						<td>
							<table style="margin: 0 auto; margin-bottom: 30px;">
							<tbody>
								<tr style="display: flex; align-items: center; justify-content: space-between; width: 650px;">
								<td><img style="max-width: 100%;" src="https://shopsphereafrica.com/images/logo.png" alt=""></td>
								</tr>
							</tbody>
							</table>
						</td>
						</tr>
					</tbody>
					</table>
					<table style="margin: 0 auto; background-color: #fff; border-radius: 8px;">
					<tbody>
						<tr>
						<td style="padding: 30px;"> 
							<h6 style="font-weight: 600; font-size: 16px; margin: 0 0 18px 0;">Password Reset</h6>
							<p style="font-size: 13px; line-height: 1.7; letter-spacing: 0.7px; margin-top: 0;">Hi {{.Name}},</p>
							<p style="font-size: 13px; line-height: 1.7; letter-spacing: 0.7px; margin-top: 0;">We've received a request to reset your Business Connect Admin password. If you didn't make this request, please ignore this email. Your account is still secure.</p>
							<p style="text-align: center;"><a href="{{.URL}}" style="padding: 10px; background-color: #5C61F2; color: #fff; display: inline-block; border-radius: 30px; font-weight: 700; padding: 0.6rem 1.75rem;">Reset Password</a></p>
							<p style="font-size: 13px; line-height: 1.7; letter-spacing: 0.7px; margin-top: 0;">If you can't click the button, please copy and paste the following link into your browser:</p>
								<textarea readonly style="display: block; margin: 0 auto; padding: 10px; background-color: #f9f9f9; border: 1px solid #ccc; border-radius: 6px; width: 100%; resize: none; font-size: 14px;">{{.URL}}</textarea>
							<p style="font-size: 13px; line-height: 1.7; letter-spacing: 0.7px; margin-top: 0;">If you remember your password you can safely ignore this email.</p>
							<p style="font-size: 13px; line-height: 1.7; letter-spacing: 0.7px; margin-top: 0;">Good luck! Hope it works.</p>
							<p style="font-size: 13px; line-height: 1.7; letter-spacing: 0.7px; margin-top: 0;">For security reasons, this link will expire in an hour.</p>
							<p style="margin-bottom: 0; font-size: 13px; line-height: 1.7; letter-spacing: 0.7px; margin-top: 0;">
							Regards,<br>Business Connect Team</p>
						</td>
						</tr>
					</tbody>
					</table>
					<table align="center" border="0" cellpadding="0" cellspacing="0" style="width: 100%; max-width: 600px; background-color: #ffffff; box-shadow: 0px 0px 14px -4px rgba(0, 0, 0, 0.27); border-radius: 10px; padding: 30px; margin-bottom: 20px;">
						<tr>
							<td style="text-align: center; padding: 30px 0;">
								<h4 style="color: #2d2d2d; margin-bottom: 20px;">Follow Us</h4>
								
								<a href="https://x.com/Shopsphereafric?t=LHeC9C4ywm9f5GO1kqutTw&s=09" style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/x.png" alt="twitter" style="width: 24px; height: 24px;"></a>
								<a href="https://www.instagram.com/officialshopsphereafrica?igsh=aGppbjNzZmRpa3lv" style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/instagram.png" alt="instagram" style="width: 24px; height: 24px;"></a>
								<a href="https://wa.me/2348027830748?text=Hello%20ShopSphere%20Africa%2C%20I'm%20interested%20in%20your%20products%20and%20would%20like%20to%20know%20more." style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/whatsapp1.png" alt="whatsapp" style="width: 24px; height: 24px;"></a>
							</td>
						</tr>
						<tr>
							<td style="text-align: center; padding: 10px 30px 30px 30px; background-color: #f4f4f4;">
								<p style="font-size: 13px; margin: 0;">© {{.Year}} Business Connect.</p>
							</td>
						</tr>
					</table>
			</td>
			</tr>
			</tbody>
		</table>
		</div>
	</body>
	</html>
    `
	currentYear := time.Now().Year()
	// Sample data
	data := struct {
		Name string
		URL  string
		Year int
	}{
		Name: name,
		URL:  "https://shopsphereafrica.com/new_password.html?user=" + sendTo + "&token=" + otp,
		Year: currentYear,
	}

	// Create a new template and parse the HTML
	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}

	// Execute the template with the provided data
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		panic(err)
	}

	// Convert the buffer to a string to get the final HTML content
	content := body.String()

	to := []string{sendTo}

	emailSendErr = sender.SendEmail(subject, content, to, nil, nil, nil)
	if emailSendErr != nil {
		fmt.Println(emailSendErr)
		return fmt.Errorf("failed to send email: %w", emailSendErr)
	}

	// let's check if the user has a previous otp stored
	oldOTP, getErr := dbFunc.DBHelper.GetOTPByEmail(sendTo)

	if getErr != nil {
		if getErr.Error() == "otp not found by email" {
			var usersOTP Data.OTP
			usersOTP.Email = sendTo
			usersOTP.OTP = otp
			usersOTP.CreatedAT = time.Now().Add(60 * time.Minute).Unix()
			usersOTP.EmailVerification = false
			usersOTP.PasswordReset = true

			saveErr := dbFunc.DBHelper.CreateOTP(usersOTP)
			if saveErr != nil {
				fmt.Println(saveErr)
				return fmt.Errorf("failed to save email otp to db: %w", saveErr)
			}
			return nil
		} else {
			return fmt.Errorf("failed to retrieve email otp from db: %w", getErr)
		}
	}

	// if user exist let's update the users otp to a new one
	newOTP.OTP = otp
	newOTP.CreatedAT = time.Now().Add(60 * time.Minute).Unix()
	newOTP.MaxTry = 0
	// log.Println("new otp to update: ", newOTP)
	updateErr := dbFunc.DBHelper.UpdateExistingOTP(newOTP, oldOTP.CustomID)
	if updateErr != nil {
		return fmt.Errorf("failed to save email otp to db: %w", updateErr)
	}

	return nil
}

func MagicLinkEmailVerification(name, sendTo string) error {
	var (
		newOTP Data.OTP
	)

	if os.Getenv("RENDER") == "" {
		// Local development only
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Failed to load .env file: %v\n", err)
		}
	}


	config := EmailConfig{
		Name:              mustEnv("EMAIL_SENDER_NAME"),
		FromEmailAddress:  mustEnv("EMAIL_SENDER_ACCOUNT"),
		FromEmailPassword: mustEnv("EMAIL_SENDER_PASSWORD"),
	}

	otp, randError = EmailOTPGenerator(30)
	if randError != nil {
		return randError
	}

	sender := NewGmailSender(config)
	subject := "Business Connect Magic Login Link"
	// HTML template with %s as a placeholder for the OTP
	htmlTemplate := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Magic Login Link</title>
	</head>
	<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; text-align: center;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); text-align: left;">
		<img src="https://payuee.shop/assets/images/logo.png" alt="Business Connect Logo" style="display: block; margin: 0 auto; max-width: 100%;">
		<h1 style="color: #333333; margin-bottom: 20px; text-align: center;">Magic Login Link</h1>
	
		<p style="color: #777777;">Hi {{.Name}},</p>
		<p style="color: #777777;">You requested a magic login link to access your Business Connect account. This link will automatically log you into your Business Connect account.</p>
		<p style="color: #777777;">To login, please click the button below:</p>
	
		<a href="{{.URL}}" style="display: block; margin: 0 auto; padding: 15px; background-color: #007bff; color: #ffffff; text-decoration: none; font-size: 24px; border-radius: 6px; width: 200px; text-align: center;">Login</a>
	
		<p style="color: #777777;">If you can't click the button, please copy and paste the following link into your browser:</p>
		<textarea readonly style="display: block; margin: 0 auto; padding: 10px; background-color: #f9f9f9; border: 1px solid #ccc; border-radius: 6px; width: 100%; resize: none; font-size: 14px;">{{.URL}}</textarea>
	
		<p style="color: #777777; margin-top: 20px;">For security reasons, this link will expire in 5 minutes.</p>
		<p style="color: #777777;">If you didn't request this login link or need assistance, please ignore this email or contact us at <a href="mailto:support@businessconnectt.com">support@businessconnectt.com</a>.</p>
		<p style="color: #777777;">Thanks,<br>The Business Connect Team</p>
	</div>
	</body>
	</html>
    `
	// Sample data
	data := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  "https://payuee.shop/dashboard/sign-in?user=" + sendTo + "&magic-code=" + otp,
	}

	// Create a new template and parse the HTML
	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}

	// Execute the template with the provided data
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		panic(err)
	}

	// Convert the buffer to a string to get the final HTML content
	content := body.String()

	to := []string{sendTo}

	emailSendErr = sender.SendEmail(subject, content, to, nil, nil, nil)
	if emailSendErr != nil {
		fmt.Println(emailSendErr)
		return fmt.Errorf("failed to send email: %w", emailSendErr)
	}

	// let's check if the user has a previous otp stored
	oldOTP, getErr := dbFunc.DBHelper.GetOTPByEmail(sendTo)

	if getErr != nil {
		if getErr.Error() == "otp not found by email" {
			var usersOTP Data.OTP
			usersOTP.Email = sendTo
			usersOTP.OTP = otp
			usersOTP.CreatedAT = time.Now().Add(5 * time.Minute).Unix()
			usersOTP.EmailVerification = false
			usersOTP.PasswordReset = true

			saveErr := dbFunc.DBHelper.CreateOTP(usersOTP)
			if saveErr != nil {
				fmt.Println(saveErr)
				return fmt.Errorf("failed to save email otp to db: %w", saveErr)
			}
			return nil
		} else {
			return fmt.Errorf("failed to retrieve email otp from db: %w", getErr)
		}
	}

	// if user exist let's update the users otp to a new one
	newOTP.OTP = otp
	newOTP.CreatedAT = time.Now().Add(5 * time.Minute).Unix()
	newOTP.MaxTry = 0
	log.Println("new otp to update: ", newOTP)
	updateErr := dbFunc.DBHelper.UpdateExistingOTP(newOTP, oldOTP.CustomID)
	if updateErr != nil {
		return fmt.Errorf("failed to save email otp to db: %w", updateErr)
	}

	return nil
}

func SendEmailToSubscribers(Subject, content, sendTo string) error {

	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Println(envErr)
		return fmt.Errorf("failed to load .env file: %w", envErr)
	}

	config := EmailConfig{
		// Name:              os.Getenv("ADMIN_EMAIL_SENDER_NAME"),
		FromEmailAddress:  os.Getenv("EMAIL_SENDER_ACCOUNT"),
		FromEmailPassword: os.Getenv("EMAIL_SENDER_PASSWORD"),
	}

	sender := NewGmailSender(config)
	subject := Subject

	// HTML template with {{.Content}} as a placeholder for the content
	htmlTemplate := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Business Connect Email Notification</title>
	</head>
	<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; text-align: center;">
		<div style="max-width: 600px; margin: 0 auto; padding: 20px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); text-align: left;">
			{{.Content}}
			<table align="center" border="0" cellpadding="0" cellspacing="0" style="width: 100%; max-width: 600px; background-color: #ffffff; box-shadow: 0px 0px 14px -4px rgba(0, 0, 0, 0.27); border-radius: 10px; padding: 30px; margin-bottom: 20px;">
				<tr>
					<td style="text-align: center; padding: 30px 0;">
						<h4 style="color: #2d2d2d; margin-bottom: 20px;">Follow Us</h4>
						
						<a href="https://x.com/Shopsphereafric?t=LHeC9C4ywm9f5GO1kqutTw&s=09" style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/x.png" alt="twitter" style="width: 24px; height: 24px;"></a>
						a href="https://www.instagram.com/officialshopsphereafrica?igsh=aGppbjNzZmRpa3lv" style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/instagram.png" alt="instagram" style="width: 24px; height: 24px;"></a>
						<a href="https://wa.me/+2349067009539?text=Hello,+I+am+one+of+your+email+subscriber+and+i+just+recieved+an+email+and+I+would+like+to+know+more+about+it." style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/whatsapp1.png" alt="whatsapp" style="width: 24px; height: 24px;"></a>
					</td>
				</tr>
				<tr>
					<td style="text-align: center; padding: 10px 30px 30px 30px; background-color: #f4f4f4;">
						<p style="font-size: 13px; margin: 0;">© {{.Year}} Business Connect.</p>
					</td>
				</tr>
			</table>
		</div>
	</body>
	</html>`

	currentYear := time.Now().Year()

	// Sample data with template.HTML to avoid escaping the content
	data := struct {
		Content template.HTML
		Year    int
	}{
		Content: template.HTML(content), // Marks content as safe HTML
		Year:    currentYear,
	}

	// Create a new template and parse the HTML
	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Execute the template with the provided data
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Convert the buffer to a string to get the final HTML content
	bodyContent := body.String()

	to := []string{sendTo}

	emailSendErr := sender.SendEmail(subject, bodyContent, to, nil, nil, nil)
	if emailSendErr != nil {
		fmt.Println(emailSendErr)
		return fmt.Errorf("failed to send email: %w", emailSendErr)
	}

	return nil
}

func EmailOTPGenerator(length int) (string, error) {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	max := big.NewInt(int64(len(chars)))
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = chars[n.Int64()]
	}
	return string(b), nil
}

func EmailOTPGeneratorNumber(length int) (string, error) {
	const chars = "0123456789"

	max := big.NewInt(int64(len(chars)))
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = chars[n.Int64()]
	}
	return string(b), nil
}