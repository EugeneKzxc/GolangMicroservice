package main

var htmlTamplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Order Details</title>
    <style>
        body { font-family: Arial, sans-serif; }
        ul { list-style-type: none; }
        li { margin-bottom: 5px; }
        .container { margin-bottom: 20px; }
    </style>
</head>
<body>
	<form action="/" method="post">
		<input type="text" name="id" placeholder="Enter Order ID" required>
		<button type="submit">Get Order Details</button>
	</form>
    {{ if .OrderUID }}
    <div class="container">
        <h1>Order Details for ID: {{ .OrderUID }}</h1>
        <p>Track Number: {{ .TrackNumber }}</p>
        <p>Entry: {{ .Entry }}</p>
        
        <h2>Delivery Information</h2>
        <ul>
            <li>Name: {{ .Delivery.Name }}</li>
            <li>Phone: {{ .Delivery.Phone }}</li>
            <li>Zip: {{ .Delivery.Zip }}</li>
            <li>City: {{ .Delivery.City }}</li>
            <li>Address: {{ .Delivery.Address }}</li>
            <li>Region: {{ .Delivery.Region }}</li>
            <li>Email: {{ .Delivery.Email }}</li>
        </ul>
        
        <h2>Payment Information</h2>
        <ul>
            <li>Transaction: {{ .Payment.Transaction }}</li>
            <li>Request ID: {{ .Payment.RequestID }}</li>
            <li>Currency: {{ .Payment.Currency }}</li>
            <li>Provider: {{ .Payment.Provider }}</li>
            <li>Amount: {{ .Payment.Amount }}</li>
            <li>Payment Date: {{ .Payment.PaymentDt }}</li>
            <li>Bank: {{ .Payment.Bank }}</li>
            <li>Delivery Cost: {{ .Payment.DeliveryCost }}</li>
            <li>Goods Total: {{ .Payment.GoodsTotal }}</li>
            <li>Custom Fee: {{ .Payment.CustomFee }}</li>
        </ul>

        <h2>Items</h2>
        {{ range .Items }}
        <div>
            <h3>Item: {{ .Name }}</h3>
            <ul>
                <li>Chart ID: {{ .ChrtID }}</li>
                <li>Track Number: {{ .TrackNumber }}</li>
                <li>Price: {{ .Price }}</li>
                <li>RID: {{ .RID }}</li>
                <li>Sale: {{ .Sale }}</li>
                <li>Size: {{ .Size }}</li>
                <li>Total Price: {{ .TotalPrice }}</li>
                <li>NmID: {{ .NmID }}</li>
                <li>Brand: {{ .Brand }}</li>
                <li>Status: {{ .Status }}</li>
            </ul>
        </div>
        {{ end }}

        <h2>Additional Information</h2>
        <ul>
            <li>Locale: {{ .Locale }}</li>
            <li>Internal Signature: {{ .InternalSignature }}</li>
            <li>Customer ID: {{ .CustomerID }}</li>
            <li>Delivery Service: {{ .DeliveryService }}</li>
            <li>Shard Key: {{ .ShardKey }}</li>
            <li>SM ID: {{ .SmID }}</li>
            <li>Date Created: {{ .DateCreated }}</li>
            <li>OOF Shard: {{ .OofShard }}</li>
        </ul>
    </div>
    {{ else }}
    <p>No order found with the given ID.</p>
    {{ end }}
</body>
</html>
`
