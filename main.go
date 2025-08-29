package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	LoadTemplates()
	LoadInvoices()

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/add", AddInvoiceHandler)
	http.HandleFunc("/mark-paid", MarkPaidHandler)
	http.HandleFunc("/receipt", ReceiptHandler)
	http.HandleFunc("/style.css", StyleCSS)

	go func() {
		for {
			now := time.Now()
			next := now.Add(24 * time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), 8, 0, 0, 0, next.Location())
			time.Sleep(next.Sub(now))

			today := CurrentDate()

			// Send "due today" reminders
			for i := range Invoices {
				inv := &Invoices[i]
				if inv.Paid || inv.Sent || inv.DueDate != today {
					continue
				}
				subject := fmt.Sprintf("Reminder: R%.2f Due Today", inv.Amount)
				body := fmt.Sprintf("Hi %s,\n\nThis is a reminder that your payment of R%.2f is due today.\n\nThank you!", inv.Client, inv.Amount)
				SendEmail(inv.Email, subject, body)
				inv.Sent = true
				SaveInvoices()
			}

			// Send overdue reminders (uses the function)
			sendOverdueReminders()
		}
	}()

	fmt.Println("🚀 Dusk Rose Pty (Ltd) InvoiceBot running on http://localhost:8080")
	fmt.Println("📄 Receipts available under 'Paid Invoices'")
	http.ListenAndServe(":8080", nil)
}