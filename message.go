package main

// Message represents a structured data type for storing and exchanging messages.
// It contains three fields: Sender, Recipient, and Content.
// Sender represents the sender of the message.
// Recipient represents the recipient of the message.
// Content represents the content or body of the message.
// The fields Sender and Recipient are identified by tags in JSON format.
// These tags are used when serializing or deserializing the Message struct to JSON.
// The "json" tag specifies the corresponding JSON key for each field.
// The ",omitempty" option indicates that if the field's value is empty, it will be omitted in the JSON representation.
type Message struct {
	Sender    string `json:"sender,omitempty"`    // The sender of the message
	Recipient string `json:"recipient,omitempty"` // The recipient of the message
	Content   string `json:"content,omitempty"`   // The content of the message
}
