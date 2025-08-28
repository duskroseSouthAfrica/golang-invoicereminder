package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
)

var homeTmpl *template.Template
var receiptTmpl *template.Template

const (
	SMTPUser = "yourbot@gmail.com"
	SMTPPass = "your-app-password"
	SMTPHost = "smtp.gmail.com"
	SMTPPort = "587"
)

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

	if dueDate == CurrentDate() {
		subject := fmt.Sprintf("Payment Due: R%.2f from %s", amt, business)
		body := fmt.Sprintf("Hi %s,\n\nYour payment of R%.2f is due today.\n\nThank you!", client, amt)
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
			break
		}
	}
	SaveInvoices()
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