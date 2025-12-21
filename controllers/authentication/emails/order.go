package emails

import (
	OrderEmail "business-connect/controllers/authentication"
	Data "business-connect/models"
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/joho/godotenv"
)

func ShopsphereConfirmationEmail(OrderHistoryBody Data.OrderHistory, RowEmailDataProducts []Data.ProductOrder) error {

	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Println(envErr)
		return fmt.Errorf("failed to load .env file: %w", envErr)
	}

	if OrderHistoryBody.CustomerStreetAddress2 == "" {
		OrderHistoryBody.CustomerStreetAddress2 = OrderHistoryBody.CustomerStreetAddress1
	}

	config := OrderEmail.EmailConfig{
		Name:              os.Getenv("ADMIN_EMAIL_SENDER_NAME"),
		FromEmailAddress:  os.Getenv("ADMIN_EMAIL_SENDER_ACCOUNT"),
		FromEmailPassword: os.Getenv("ADMIN_EMAIL_SENDER_PASSWORD"),
	}

	sender := OrderEmail.NewGmailSender(config)
	subject := "Shopsphere Africa Order Confirmation"

	htmlTemplate := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Shopsphere Africa Admin</title>
	</head>
	<body style="margin: 20px auto; font-family: Arial, sans-serif; background-color: #f4f4f4;">
		<table align="center" border="0" cellpadding="0" cellspacing="0" style="width: 100%; max-width: 600px; background-color: #ffffff; box-shadow: 0px 0px 14px -4px rgba(0, 0, 0, 0.27); border-radius: 10px; padding: 30px; margin-bottom: 20px;">
			<tr>
				<td style="text-align: center;">
					<img src="https://shopsphereafrica.com/image/catalog/logo.png" alt="" style="margin-bottom: 30px; border-radius: 10px; width: 100%; max-width: 560px;">
				</td>
			</tr>
			<tr>
				<td style="text-align: left; color: #717171;">
					<h4 style="color: #333333;">Hi {{.Name}},</h4>
					<p>Your order has been successfully processed and is on its way!</p>
					<p>Transaction ID: {{.OrderID}}</p>
					<p>You can track your order using the link below:</p>
					<p><a href="https://shopsphereafrica.com/track-order?OrderID={{.OrderID}}" style="color: #0066cc;" target="_blank">Track My Order</a></p>
				</td>
			</tr>
			<tr>
				<td style="padding: 15px 0;">
					<table cellpadding="0" cellspacing="0" border="0" style="width: 100%;">
						<tr>
							<td style="background-color: rgba(62,95,206,0.02); border: 1px solid #eeeeee; padding: 15px; width: 50%; border-radius: 10px;">
								<h5 style="font-size: 16px; font-weight: 600; color: #000; line-height: 1.2; margin: 0 0 13px 0;">Your Shipping Address 1:</h5>
								<p style="font-size: 14px; color: #717171; line-height: 1.5; margin: 0;">{{.CustomerAddress1}}</p>
							</td>
							<td style="width: 30px;"><img src="https://shopsphereafrica.com/admin/images/email-template/space.jpg" alt=" " height="25" width="30"></td>
							<td style="background-color: rgba(62,95,206,0.02); border: 1px solid #eeeeee; padding: 15px; width: 50%; border-radius: 10px;">
								<h5 style="font-size: 16px; font-weight: 600; color: #000; line-height: 1.2; margin: 0 0 13px 0;">Your Shipping Address 2:</h5>
								<p style="font-size: 14px; color: #717171; line-height: 1.5; margin: 0;">{{.CustomerAddress2}}</p>
							</td>
						</tr>
					</table>
				</td>
			</tr>
			<tr>
				<td>
					<table border="0" cellpadding="0" cellspacing="0" style="width: 100%; margin-bottom: 30px;">
						<tr>
							<th style="text-align: left; font-size: 14px; padding-bottom: 10px;">PRODUCT</th>
							<th style="text-align: left; font-size: 14px; padding-left: 15px; padding-bottom: 10px;">DESC</th>
							<th style="text-align: left; font-size: 14px; padding-bottom: 10px;">QTY</th>
							<th style="text-align: right; font-size: 14px; padding-bottom: 10px;">PRICE</th>
						</tr>
						{{range .RowEmailDataProducts}}
						<tr style="border-top: 1px solid #eeeeee;">
							<td style="padding-top: 10px; text-align: center;">
								<img src="https://shopsphereafrica.com/image/{{.Image1}}" alt="{{.Title}}" style="width: 80px;">
							</td>
							<td style="padding-left: 15px;">
								<h5 style="font-size: 1em; margin: 15px 0; color: #333333;">{{.Title}}</h5>
							</td>
							<td style="padding-left: 15px;">
								<p style="font-size: 0.8em; color: #444; margin: 15px 0;">QTY: {{.Quantity}}</p>
							</td>
							<td style="text-align: right;">
								<p style="font-size: 1em; color: #444; margin: 15px 0;"><b>₦{{.OrderCost}}</b></p>
							</td>
						</tr>
						{{end}}
						<tr style="border-top: 1px solid #eeeeee;">
							<td colspan="2" style="padding-top: 15px; font-size: 14px;">SUBTOTAL:</td>
							<td colspan="2" style="padding-top: 15px; text-align: right; font-size: 14px;"><b>₦{{.SubTotal}}</b></td>
						</tr>
						<tr>
							<td colspan="2" style="padding-top: 10px; font-size: 14px;">SHIPPING:</td>
							<td colspan="2" style="padding-top: 10px; text-align: right; font-size: 14px;"><b>₦{{.ShippingCharge}}</b></td>
						</tr>
						<tr>
							<td colspan="2" style="padding-top: 10px; font-size: 14px;">DISCOUNT:</td>
							<td colspan="2" style="padding-top: 10px; text-align: right; font-size: 14px;"><b>₦{{.Discount}}</b></td>
						</tr>
						<tr>
							<td colspan="2" style="padding-top: 15px; font-size: 14px;">TOTAL:</td>
							<td colspan="2" style="padding-top: 15px; text-align: right; font-size: 14px;"><b>₦{{.Total}}</b></td>
						</tr>
					</table>
				</td>
			</tr>
			<tr>
				<td style="text-align: center; padding: 30px 0;">
					<h4 style="color: #2d2d2d; margin-bottom: 20px;">Follow Us</h4>
					
					<a href="https://x.com/Shopsphereafric?t=LHeC9C4ywm9f5GO1kqutTw&s=09" style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/x.png" alt="twitter" style="width: 24px; height: 24px;"></a>
					<a href="https://www.instagram.com/todac_wellness/" style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/instagram.png" alt="instagram" style="width: 24px; height: 24px;"></a>
					<a href="https://wa.me/+2347025455850?text=Hello,+I+just+placed+an+order+and+would+like+more+information.+My+Order+ID+is+{{.OrderID}}" style="margin: 0 10px;"><img src="https://shopsphereafrica.com/admin/images/social/whatsapp1.png" alt="whatsapp" style="width: 24px; height: 24px;"></a>
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
		OrderID              uint
		Name                 string
		RowEmailDataProducts []Data.ProductOrder
		CustomerAddress1     string
		CustomerAddress2     string
		SubTotal             string
		ShippingCharge       string
		Discount             string
		Total                string
		Year                 int
	}{
		OrderID:              OrderHistoryBody.ID,
		Name:                 OrderHistoryBody.CustomerFName + " " + OrderHistoryBody.CustomerSName,
		RowEmailDataProducts: RowEmailDataProducts,
		CustomerAddress1:     OrderHistoryBody.CustomerStreetAddress1,
		CustomerAddress2:     OrderHistoryBody.CustomerStreetAddress2,
		SubTotal:             formatNaira(OrderHistoryBody.OrderSubTotalCost),
		ShippingCharge:       formatNaira(OrderHistoryBody.ShippingCost),
		Discount:             formatNaira(OrderHistoryBody.OrderDiscount),
		Total:                formatNaira(OrderHistoryBody.OrderSubTotalCost + OrderHistoryBody.ShippingCost),
		Year:                 currentYear,
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

	to := []string{OrderHistoryBody.CustomerEmail}

	emailSendErr := sender.SendEmail(subject, content, to, nil, nil, nil)
	if emailSendErr != nil {
		fmt.Println(emailSendErr)
		return fmt.Errorf("failed to send email: %w", emailSendErr)
	}

	return nil
}

// Add commas to the integer part of the number
func addCommas(n string) string {
	var result strings.Builder

	// Start from the rightmost digit and add commas every three digits
	for i, digit := range n {
		if i > 0 && (len(n)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}

	return result.String()
}

// Format the amount with commas, two decimal places, and the Naira symbol
func formatNaira(amount float64) string {
	// Convert the amount to a string with two decimal places
	amountStr := fmt.Sprintf("%.2f", amount)

	// Split the string into integer and fractional parts
	parts := strings.Split(amountStr, ".")
	intPart := parts[0]
	fracPart := parts[1]

	// Add commas to the integer part
	intPartWithCommas := addCommas(intPart)

	// Combine the integer part, fractional part, and Naira symbol
	return intPartWithCommas + "." + fracPart
}
