package ai

import (
	"context"
	"fmt"
	"net/http"

	"os"
	"strings"

	// dbFunc "business-connect/database/dbHelpFunc"
	// cookieNul "business-connect/middleware"
	Data "business-connect/models"
	// payueeTrans "business-connect/payueeTrans"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func GetVendorProductDescriptionAI(ctx *fiber.Ctx) error {

	// Define the struct type
	type AiDescriptionBody struct {
		Title string
	}

	// Initialize an instance of AiDescriptionBody
	var aiDescriptionBody AiDescriptionBody
	if bindErr := ctx.BodyParser(&aiDescriptionBody); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read body"})
	}

	description, aiErr := SendUserQuestion(aiDescriptionBody.Title)
	if aiErr != nil {
		fmt.Println("this is the ai error ", aiErr)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "AI Description Generation Timed Out"})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": description,
		// "collaborate": canCollaborate,
	})
}

func GetVendorProductTagAI(ctx *fiber.Ctx) error {
	// Define the struct type
	type AiDescriptionBody struct {
		Title       string
		Description string
	}

	// Initialize an instance of AiDescriptionBody
	var aiDescriptionBody AiDescriptionBody
	if bindErr := ctx.BodyParser(&aiDescriptionBody); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read body"})
	}

	tag, aiErr := SendUserQuestionTag(aiDescriptionBody.Title, aiDescriptionBody.Description)
	if aiErr != nil {
		fmt.Println("this is the ai error ", aiErr)
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "AI Description Generation Timed Out"})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": tag,
		// "collaborate": canCollaborate,
	})
}

func SendUserQuestion(productTitle string) (response string, err error) {
	ctx := context.Background()
	envErr := godotenv.Load(".env")

	if envErr != nil {
		fmt.Printf("Failed to load .env file: %v\n", envErr)
	}

	apiKey := os.Getenv("AI_API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}

	defer client.Close()

	chatBoard := []Data.ChatMessages{
		{Role: "user", Text: "Using only the title provided, generate a product description with no other text or explanation. Provide only the description itself without any additional comments or instructions. The description should be detailed enough to attract customers and highlight key features and benefits of the product. Title: Samsong galaxy s23"},
		{Role: "model", Text: "The Samsung Galaxy S23 is a sleek and powerful smartphone designed for those who demand the best in performance and style. With a stunning 6.1-inch Dynamic AMOLED display, it brings vibrant colors and crisp details to everything from streaming to gaming. Powered by a high-performance processor, it handles multitasking and heavy apps with ease. Capture breathtaking photos and videos with its advanced triple-lens camera system, featuring enhanced low-light capabilities and ultra-high resolution. With 5G connectivity, a long-lasting battery, and fast charging, the Galaxy S23 keeps you connected and productive all day. Enjoy cutting-edge features, premium build quality, and Samsung’s renowned innovation with the Galaxy S23 – your ultimate companion for modern living."},
	}

	// Initialize chat history for the current session
	var chatHistory []*genai.Content
	for _, message := range chatBoard {
		chatHistory = append(chatHistory, &genai.Content{
			Parts: []genai.Part{
				genai.Text(message.Text),
			},
			Role: message.Role,
		})
	}

	model := client.GenerativeModel("gemini-2.0-flash")
	cs := model.StartChat()
	cs.History = chatHistory // Set chat history

	// Send user's question as the next message
	resp, err := cs.SendMessage(ctx, genai.Text(RegisteredTrainingData(productTitle)))
	if err != nil {
		return "", err
	}

	// Extract response from candidates and return it
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		var responseParts []string
		for _, part := range resp.Candidates[0].Content.Parts {
			if textPart, ok := part.(genai.Text); ok {
				responseParts = append(responseParts, string(textPart))
			}
		}
		response = strings.Join(responseParts, " ")
	}

	responseText := response
	modifiedResponse := strings.ReplaceAll(responseText, "Gemini", "Ngozi")
	modifiedResponse1 := strings.ReplaceAll(modifiedResponse, "Google", "BusinessConnect")

	return modifiedResponse1, nil
}

func RegisteredTrainingData(productTitle string) string {
	// Instruction for generating a pure description with no extra text or symbols
	RoleText := `Generate a high-quality product description using only the product title provided. 
	Do not add any additional text, explanations, symbols, instructions, or introductory phrases—only return the pure description. 
	The description should capture the product’s key features, benefits, and appeal to potential customers in a concise and compelling way.
	
	Title: ` + productTitle

	return RoleText
}

func RegisteredTrainingDataTag(productTitle, productDescription string) string {
	// Instruction for generating tags without any extra formatting or symbols
	RoleText := `Generate a list of product tags using only the product title and description provided. 
	Ensure the tags are relevant, concise, and separated by commas. 
	Do not add any additional text, explanations, symbols, or instructions—just return the raw tags.
	Title: ` + productTitle + ". " + " Description: " + productDescription

	return RoleText
}

func SendUserQuestionTag(productTitle, productDescription string) (response string, err error) {
	ctx := context.Background()
	envErr := godotenv.Load(".env")

	if envErr != nil {
		fmt.Printf("Failed to load .env file: %v\n", envErr)
	}

	apiKey := os.Getenv("AI_API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}

	defer client.Close()

	chatBoard := []Data.ChatMessages{
		{Role: "user", Text: "Generate a list of product tags using only the product title and description provided. The tags should be relevant, concise, and separated by commas. Do not include any additional text, explanations, or instructions—just the tags. Ensure that the tags reflect key attributes, features, and categories of the product. Title: Samsong galaxy s23. Description: The Samsung Galaxy S23 is a sleek and powerful smartphone designed for those who demand the best in performance and style. With a stunning 6.1-inch Dynamic AMOLED display, it brings vibrant colors and crisp details to everything from streaming to gaming. Powered by a high-performance processor, it handles multitasking and heavy apps with ease. Capture breathtaking photos and videos with its advanced triple-lens camera system, featuring enhanced low-light capabilities and ultra-high resolution. With 5G connectivity, a long-lasting battery, and fast charging, the Galaxy S23 keeps you connected and productive all day. Enjoy cutting-edge features, premium build quality, and Samsung’s renowned innovation with the Galaxy S23 – your ultimate companion for modern living."},
		{Role: "model", Text: "Samsung Galaxy S23, smartphone, 5G, Dynamic AMOLED display, triple-lens camera, high-performance processor, long-lasting battery, fast charging, 6.1-inch display, modern technology, mobile photography, Android, premium smartphone, low-light photography, gaming, multitasking, connectivity, Samsung, flagship phone."},
	}

	// Initialize chat history for the current session
	var chatHistory []*genai.Content
	for _, message := range chatBoard {
		chatHistory = append(chatHistory, &genai.Content{
			Parts: []genai.Part{
				genai.Text(message.Text),
			},
			Role: message.Role,
		})
	}

	model := client.GenerativeModel("gemini-2.0-flash")
	cs := model.StartChat()
	cs.History = chatHistory // Set chat history

	// Send user's question as the next message
	resp, err := cs.SendMessage(ctx, genai.Text(RegisteredTrainingDataTag(productTitle, productDescription)))
	if err != nil {
		return "", err
	}

	// Extract response from candidates and return it
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		var responseParts []string
		for _, part := range resp.Candidates[0].Content.Parts {
			if textPart, ok := part.(genai.Text); ok {
				responseParts = append(responseParts, string(textPart))
			}
		}
		response = strings.Join(responseParts, " ")
	}

	responseText := response
	modifiedResponse := strings.ReplaceAll(responseText, "Gemini", "Ngozi")
	modifiedResponse1 := strings.ReplaceAll(modifiedResponse, "Google", "BusinessConnect")

	return modifiedResponse1, nil
}
