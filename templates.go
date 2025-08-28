package main

var HomeHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Dusk Rose Invoice Reminder</title>
  <link rel="stylesheet" href="/style.css">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
  <div class="container">
    <h1>📄 Dusk Rose Invoice Reminder</h1>
    <p class="tagline">Automated payment reminders for South African small businesses</p>

    <form method="post" action="/add" class="invoice-form">
      <h2>Create New Invoice</h2>

      <label>Your Business Name
        <input type="text" name="business" placeholder="e.g. Thando Tutoring" required>
      </label>

      <label>Client Name
        <input type="text" name="client" placeholder="e.g. Sipho Mokoena" required>
      </label>

      <label>Client Email
        <input type="email" name="email" placeholder="e.g. sipho@email.com" required>
      </label>

      <label>Amount (ZAR)
        <input type="number" name="amount" placeholder="400" min="1" required>
      </label>

      <label>Due Date
        <input type="date" name="dueDate" required>
      </label>

      <button type="submit">Create Invoice</button>
    </form>

    <h2>📬 Active Invoices</h2>
    <table>
      <tr>
        <th>ID</th>
        <th>Client</th>
        <th>Email</th>
        <th>Amount</th>
        <th>Due</th>
        <th>Status</th>
        <th>Action</th>
      </tr>
      {{range .Invoices}}
      {{if not .Paid}}
      <tr>
        <td>{{.ID}}</td>
        <td>{{.Client}}</td>
        <td>{{.Email}}</td>
        <td>R{{printf "%.2f" .Amount}}</td>
        <td>{{.DueDate}}</td>
        <td>{{if .Sent}}🔴 Sent{{else}}🟡 Pending{{end}}</td>
        <td>
          <form method="post" action="/mark-paid" style="display:inline">
            <input type="hidden" name="id" value="{{.ID}}">
            <button type="submit">Mark Paid</button>
          </form>
        </td>
      </tr>
      {{end}}
      {{end}}
    </table>

    <h2>✅ Paid Invoices</h2>
    <table>
      <tr>
        <th>ID</th>
        <th>Client</th>
        <th>Email</th>
        <th>Amount</th>
        <th>Paid On</th>
        <th>Receipt</th>
      </tr>
      {{range .Invoices}}
      {{if .Paid}}
      <tr>
        <td>{{.ID}}</td>
        <td>{{.Client}}</td>
        <td>{{.Email}}</td>
        <td>R{{printf "%.2f" .Amount}}</td>
        <td>{{.DueDate}}</td>
        <td><a href="/receipt?id={{.ID}}" target="_blank">📄 View</a></td>
      </tr>
      {{end}}
      {{end}}
    </table>
  </div>
</body>
</html>`

var ReceiptHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Receipt</title>
  <link rel="stylesheet" href="/style.css">
</head>
<body>
  <div class="receipt">
    <h2>🧾 Payment Receipt</h2>
    <p><strong>From:</strong> {{.Business}}</p>
    <p><strong>To:</strong> {{.Client}}</p>
    <p><strong>Email:</strong> {{.Email}}</p>
    <div class="amount">Amount: <strong>R{{printf "%.2f" .Amount}}</strong></div>
    <p><strong>Status:</strong> Paid on {{.DueDate}}</p>
    <p class="thank-you">Thank you for your payment!</p>
    <p class="footer">Dusk Rose Invoice Reminder • Generated on {{.CreatedAt}}</p>
  </div>
</body>
</html>`