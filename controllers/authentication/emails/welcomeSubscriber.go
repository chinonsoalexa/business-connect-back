package emails

import (
	OrderEmail "business-connect/controllers/authentication"
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/joho/godotenv"
)

func TodacWelcomeEmail(subscriberEmail string) error {

	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Println(envErr)
		return fmt.Errorf("failed to load .env file: %w", envErr)
	}

	config := OrderEmail.EmailConfig{
		Name:              os.Getenv("ADMIN_EMAIL_SENDER_NAME"),
		FromEmailAddress:  os.Getenv("ADMIN_EMAIL_SENDER_ACCOUNT"),
		FromEmailPassword: os.Getenv("ADMIN_EMAIL_SENDER_PASSWORD"),
	}

	sender := OrderEmail.NewGmailSender(config)
	subject := "Welcome to Shopsphere Africa!"

	htmlTemplate := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Welcome to Shopsphere Africa</title>
	</head>
	<body style="margin: 20px auto; font-family: Arial, sans-serif; background-color: #f4f4f4;">
		<table align="center" border="0" cellpadding="0" cellspacing="0" style="width: 100%; max-width: 600px; background-color: #ffffff; box-shadow: 0px 0px 14px -4px rgba(0, 0, 0, 0.27); border-radius: 10px; padding: 30px; margin-bottom: 20px;">
			<tr>
				<td style="text-align: center;">
					<img src="https://shopsphereafrica.com/image/catalog/logo.png" alt="Welcome" style="margin-bottom: 30px; border-radius: 10px; width: 100%; max-width: 560px;">
				</td>
			</tr>
			<tr>
				<td style="text-align: left; color: #717171;">
					<h4 style="color: #333333;">Hi Dear,</h4>
					<p>Welcome to Shopsphere Africa! We're excited that you’ve joined us in discovering quality products across a wide range of categories.</p>
					<p>Whether you're looking to upgrade your home, office, wardrobe, or workspace, you’ve made the right choice signing up with us. From essential household items to electronics, sports gear, and agricultural tools — we’ve got you covered.</p>
					<p>Here are a few things to help you get started:</p>
					<ul>
						<li>Explore our extensive collection of <a href="https://shopsphereafrica.com/shop?page=1" style="color: #0066cc;">products</a> including household essentials, textile & decor, sports & recreation, apparel accessories, transportation, electronics, computers, services, furniture, food, and office supplies.</li>
						<li>Stay informed on product highlights and shopping tips by visiting our <a href="https://shopsphereafrica.com/blog?page=1" style="color: #0066cc;">blog</a>.</li>
						<li>Need help? Just reply to this email or <a href="https://wa.me/+2347025455850?text=Hello,+I+just+joined+Shopsphere+Africa+and+would+like+more+information." style="color: #0066cc;">contact us</a>.</li>
					</ul>
					<p>Stay tuned for updates, exclusive deals, and tips from Shopsphere Africa.</p>
					<p>Thank you for choosing us. We’re here to support your shopping journey every step of the way.</p>
				</td>
			</tr>
			<tr>
				<td style="text-align: center; padding: 30px 0;">
					<h4 style="color: #2d2d2d; margin-bottom: 20px;">Follow Us</h4>
					
					<a href="https://x.com/Shopsphereafric?t=LHeC9C4ywm9f5GO1kqutTw&s=09" style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/x.png" alt="twitter" style="width: 24px; height: 24px;"></a>
					<a href="https://www.instagram.com/todac_wellness/" style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/instagram.png" alt="instagram" style="width: 24px; height: 24px;"></a>
					<a href="https://wa.me/+2347025455850?text=Hello,+I+just+subscribed+to+Todac+Wellness+Store+and+would+like+more+information." style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/whatsapp1.png" alt="whatsapp" style="width: 24px; height: 24px;"></a>
				</td>
			</tr>
			<tr>
				<td style="text-align: center; padding: 10px 30px 30px 30px; background-color: #f4f4f4;">
					<p style="font-size: 13px; margin: 0;">© {{.Year}} Shopsphere Africa.</p>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`
	currentYear := time.Now().Year()

	// Sample data
	data := struct {
		Year int
	}{
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

	to := []string{subscriberEmail}

	emailSendErr := sender.SendEmail(subject, content, to, nil, nil, nil)
	if emailSendErr != nil {
		fmt.Println(emailSendErr)
		return fmt.Errorf("failed to send email: %w", emailSendErr)
	}

	return nil
}
