package templates

var ForgotPasswordTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Reset</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            padding: 20px;
        }

        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #fff;
            padding: 30px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }

        h2 {
            color: #333;
        }

        p {
            color: #555;
            line-height: 1.6;
        }

        ol {
            color: #555;
            padding-left: 20px;
        }

        a {
            color: #007bff;
            text-decoration: none;
        }

        a:hover {
            text-decoration: underline;
        }

    </style>
</head>
<body>

    <div class="container">

        <h2>Reset Your Password</h2>

        <p>Dear %s,</p>

        <p>It seems you've forgotten your password! Not to worry, we're here to help you regain access to your account. Follow the simple steps below to reset your password:</p>

        <ol>
            <li>Click on the following link to reset your password: <a href="[Reset Password Link]">Reset Password</a></li>
            <li>You'll be directed to a page where you can create a new password. Please choose a password that is secure but memorable.</li>
            <li>Once you've set your new password, you'll be able to log back into your account as usual.</li>
        </ol>

        <p>If you didn't request this password reset, please disregard this email. Your account is still secure, and no changes have been made.</p>

        <p>If you continue to experience any issues or have any questions, feel free to reach out to our support team at <a href="mailto:[Support Email]">[Support Email]</a>.</p>

    </div>

</body>
</html>
`
