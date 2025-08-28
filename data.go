package main

import "time"

type Invoice struct {
	ID         int    `json:"id"`
	Business   string `json:"business"`
	Client     string `json:"client"`
	Email      string `json:"email"`
	Amount     float64 `json:"amount"`
	DueDate    string `json:"dueDate"`
	Sent       bool   `json:"sent"`
	Paid       bool   `json:"paid"`
	CreatedAt  string `json:"createdAt"`
}

var Invoices []Invoice
var NextID = 1

func CurrentDate() string {
	return time.Now().Format("2006-01-02")
}