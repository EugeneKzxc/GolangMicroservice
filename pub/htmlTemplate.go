package main

var htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Order Publisher</title>
</head>
<body>
    <form action="/" method="post">
        <label for="id">Order UID:</label><br>
        <input type="text" id="id" name="id" required><br><br>
        <input type="submit" value="Publish">
    </form>
</body>
</html>
`
