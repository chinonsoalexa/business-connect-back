package profile

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func pick[T any](items []T) T {
	return items[rand.Intn(len(items))]
}

var greetings = []func(string) string{
	func(r string) string { return "Hi" + r + "!" },
	func(r string) string { return "Hello" + r + "!" },
	func(r string) string { return "Hi" + r + "," },
	func(r string) string { return "Hello" + r + "," },
	func(r string) string { return "Good day" + r + "!" },
	func(r string) string { return "Hey" + r + "!" },
	func(r string) string { return "Greetings" + r + "!" },
	func(r string) string { return "Hi there" + r + "!" },
	func(r string) string { return "Hello there" + r + "!" },
	func(r string) string { return "Hey there" + r + "!" },
}

var intros = []func(string) string{
	func(s string) string { return "I'm " + s + "." },
	func(s string) string { return "My name is " + s + "." },
	func(s string) string { return s + " here." },
	func(s string) string { return "This is " + s + "." },
	func(s string) string { return "Hope you're well — I'm " + s + "." },
	func(s string) string { return "Just reaching out — I'm " + s + "." },
	func(s string) string { return "Quick intro — " + s + " here." },
	func(s string) string { return "Allow me to introduce myself — " + s + "." },
	func(s string) string { return "Hope this message finds you well. I'm " + s + "." },
	func(s string) string { return "Glad to connect — I'm " + s + "." },
}

var discoveryVerbs = []func() string{
	func() string { return "I came across" },
	func() string { return "I found" },
	func() string { return "I discovered" },
	func() string { return "I noticed" },
	func() string { return "I saw" },
	func() string { return "I stumbled on" },
	func() string { return "I recently found" },
	func() string { return "I happened to come across" },
	func() string { return "I came upon" },
	func() string { return "I was browsing and found" },
}

var businessObjects = []func(string) string{
	func(b string) string { return "your business, " + b },
	func(b string) string { return b },
	func(b string) string { return "your listing for " + b },
	func(b string) string { return "your business profile for " + b },
	func(b string) string { return "your page for " + b },
	func(_ string) string { return "your business listing" },
	func(_ string) string { return "your company profile" },
	func(_ string) string { return "your business information" },
	func(_ string) string { return "your business page" },
	func(b string) string { return "your brand, " + b },
}

var profileObjects = []func() string{
	func() string { return "your profile" },
	func() string { return "your contact" },
	func() string { return "your details" },
	func() string { return "your personal profile" },
	func() string { return "your contact information" },
	func() string { return "your listing" },
	func() string { return "your page" },
	func() string { return "your information" },
	func() string { return "your account" },
	func() string { return "your profile details" },
}

var platforms = []func() string{
	func() string { return "on Business Connect." },
	func() string { return "via Business Connect." },
	func() string { return "through Business Connect." },
	func() string { return "while browsing Business Connect." },
	func() string { return "while exploring Business Connect." },
	func() string { return "on the Business Connect platform." },
	func() string { return "during my search on Business Connect." },
	func() string { return "while checking Business Connect." },
	func() string { return "on Business Connect recently." },
	func() string { return "using Business Connect." },
}

var closers = []func() string{
	func() string { return "" },
	func() string { return " Thought I'd reach out." },
	func() string { return " and wanted to connect." },
	func() string { return " and decided to say hello." },
	func() string { return " and thought it’d be great to connect." },
	func() string { return " and wanted to reach out directly." },
	func() string { return " so I decided to message you." },
	func() string { return " and felt free to reach out." },
	func() string { return " and wanted to get in touch." },
	func() string { return " and thought I'd introduce myself." },
}

func BuildBusinessMessage(senderName, receiverName, businessName string) string {
	receiverPart := ""
	if receiverName != "" {
		receiverPart = " " + receiverName
	}

	parts := []string{
		pick(greetings)(receiverPart),
		pick(intros)(senderName),
		pick(discoveryVerbs)(),
		pick(businessObjects)(businessName),
		pick(platforms)(),
		pick(closers)(),
	}

	return strings.Join(parts, " ")
}

func BuildProfileMessage(senderName, receiverName string) string {
	receiverPart := ""
	if receiverName != "" {
		receiverPart = " " + receiverName
	}

	parts := []string{
		pick(greetings)(receiverPart),
		pick(intros)(senderName),
		pick(discoveryVerbs)(),
		pick(profileObjects)(),
		pick(platforms)(),
		pick(closers)(),
	}

	return strings.Join(parts, " ")
}

// app, web := GenerateWhatsAppLinks(
// 	"2348012345678",
// 	"Chinonso",
// 	"John",
// 	"E3 Restaurant",
// )

// fmt.Println("App:", app)
// fmt.Println("Web:", web)
