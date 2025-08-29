package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

var homeTmpl *template.Template
var receiptTmpl *template.Template

const (
	SMTPUser = "dmitri@duskrose.co.za"
	SMTPPass = "DuskRose2025!"
	SMTPHost = "mail.duskrose.co.za"
	SMTPPort = "587"
)

func sendReceipt(invoice Invoice) error {
	subject := fmt.Sprintf("Payment Receipt: R%.2f", invoice.Amount)
	body := fmt.Sprintf(`From: %s
To: %s
Subject: %s

Hello %s,

This is your payment receipt.

Amount: R%.2f
Paid on: %s

Thank you for your payment.

Dusk Rose Pty (Ltd)
`, invoice.Business, invoice.Email, subject, invoice.Client, invoice.Amount, invoice.DueDate)

	return SendEmail(invoice.Email, subject, body)
}

func sendOverdueReminders() {
	today, _ := time.Parse("2006-01-02", CurrentDate())

	for i := range Invoices {
		inv := &Invoices[i]
		if inv.Paid || !inv.Sent {
			continue // Skip if paid or not yet sent (not overdue)
		}

		dueDate, err := time.Parse("2006-01-02", inv.DueDate)
		if err != nil {
			continue
		}

		// If due date was before today → overdue
		if dueDate.Before(today) {
			subject := fmt.Sprintf("Reminder: Payment Overdue – R%.2f", inv.Amount)
			body := fmt.Sprintf(
				"Hi %s,\n\nThis is a reminder that your payment of R%.2f is overdue.\n\nPlease settle as soon as possible.\n\nThank you!",
				inv.Client, inv.Amount)

			err := SendEmail(inv.Email, subject, body)
			if err == nil {
				fmt.Printf("✅ Sent overdue reminder to %s\n", inv.Email)
			} else {
				fmt.Printf("❌ Failed to send overdue reminder to %s: %v\n", inv.Email, err)
			}
		}
	}
}

func LoadTemplates() {
	homeTmpl = template.Must(template.New("home").Parse(HomeHTML))
	receiptTmpl = template.Must(template.New("receipt").Parse(ReceiptHTML))
}

func LoadInvoices() {
	data, err := os.ReadFile("invoices.json")
	if os.IsNotExist(err) {
		fmt.Println("✅ No invoices.json — starting fresh")
		return
	}
	if err != nil {
		fmt.Println("❌ Failed to read invoices.json:", err)
		return
	}
	json.Unmarshal(data, &Invoices)

	NextID = 1
	for _, inv := range Invoices {
		if inv.ID >= NextID {
			NextID = inv.ID + 1
		}
	}
}

func SaveInvoices() {
	data, err := json.MarshalIndent(Invoices, "", "  ")
	if err != nil {
		fmt.Println("❌ Failed to marshal invoices:", err)
		return
	}
	err = os.WriteFile("invoices.json", data, 0644)
	if err != nil {
		fmt.Println("❌ Failed to write invoices.json:", err)
	}
}

func SendEmail(to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", SMTPUser, to, subject, body)
	return smtp.SendMail(
		SMTPHost+":"+SMTPPort,
		smtp.PlainAuth("", SMTPUser, SMTPPass, SMTPHost),
		SMTPUser,
		[]string{to},
		[]byte(msg),
	)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	homeTmpl.Execute(w, struct{ Invoices []Invoice }{Invoices})
}

func AddInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	business := r.FormValue("business")
	client := r.FormValue("client")
	email := r.FormValue("email")
	amount := r.FormValue("amount")
	dueDate := r.FormValue("dueDate")

	var amt float64
	fmt.Sscanf(amount, "%f", &amt)

	invoice := Invoice{
		ID:        NextID,
		Business:  business,
		Client:    client,
		Email:     email,
		Amount:    amt,
		DueDate:   dueDate,
		Sent:      false,
		Paid:      false,
		CreatedAt: CurrentDate(),
	}
	Invoices = append(Invoices, invoice)
	NextID++

	SaveInvoices()

	threeDaysBefore := time.Now().AddDate(0, 0, 3).Format("2006-01-02")
	if dueDate == threeDaysBefore {
		subject := fmt.Sprintf("Upcoming Payment: R%.2f from %s", amt, business)
		body := fmt.Sprintf("Hi %s,\n\nThis is a reminder that your payment of R%.2f is due in 3 days.\n\nThank you!", client, amt)
		SendEmail(email, subject, body)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func MarkPaidHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	for i := range Invoices {
		if Invoices[i].ID == id {
			Invoices[i].Paid = true
			SaveInvoices()
			go sendReceipt(Invoices[i])
			break
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ReceiptHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	for _, inv := range Invoices {
		if inv.ID == id && inv.Paid {
			receiptTmpl.Execute(w, inv)
			return
		}
	}
	http.Error(w, "Receipt not found", http.StatusNotFound)
}

func StyleCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "style.css")
}