# Secure Discord Webhook

[![Go Version](https://img.shields.io/github/go-mod/go-version/andrinoff/secure-discord-webhook)](https://golang.org)


A secure, serverless Discord webhook proxy built with Go and deployed on Vercel. This project provides a safe method for sending messages to a Discord webhook from a client-side application without exposing the webhook URL.

---

## 🚀 About The Project

This project acts as a secure proxy, forwarding requests from your application to your Discord webhook URL. The primary goal is to prevent the exposure of your sensitive Discord webhook URL on the client-side.

It includes two distinct serverless function endpoints:
* `/api/webhook`
* `/api/contact`

The key security feature is Origin-based restriction. The serverless function will only accept requests from a whitelisted origin (e.g., your website), which you must configure in your environment variables. This effectively blocks any unauthorized domains from using your webhook.

---

## 🛠️ Getting Started

Follow these steps to get a local copy up and running.

### Prerequisites

* Go (version 1.18 or higher is recommended)
* A [Vercel](https://vercel.com/) account for deployment
* At least one Discord webhook URL

### Installation

1.  **Clone the repository:**
    ```sh
    git clone [https://github.com/andrinoff/secure-discord-webhook.git](https://github.com/andrinoff/secure-discord-webhook.git)
    ```
2.  **Navigate to the project directory:**
    ```sh
    cd secure-discord-webhook
    ```
3.  **Install Go dependencies:**
    ```sh
    go mod tidy
    ```
4.  **Configure your environment:**
    Create a `.env` file in the root of the project and add your Discord webhook URLs. You will also need to set the allowed origin.
    ```env
    # The primary webhook URL
    DISCORD_WEBHOOK_URL="YOUR_DISCORD_WEBHOOK_URL_1"

    # A secondary webhook URL for the contact form
    DISCORD_WEBHOOK_URL_2="YOUR_DISCORD_WEBHOOK_URL_2"
    ```
    The allowed origin is hardcoded as `https://tbilisi.hackclub.com` in the Go files. For production, you should modify the `allowedOrigin` variable in `api/webhook/index.go` and `api/contact/index.go` to match your website's URL.

---

## 💡 Usage

Once deployed to Vercel, your application will expose two API endpoints:

* `https://your-vercel-app-url.vercel.app/api/webhook`
* `https://your-vercel-app-url.vercel.app/api/contact`

You can trigger these webhooks by sending a `POST` request from your allowed origin with a JSON payload.

**Example Payload:**
```json
{
  "content": "Hello from your secure webhook!",
  "username": "My Awesome App"
}
```

* `content`: The message text that will be sent to your Discord channel.
* `username` (optional): This will override the default username of the webhook for this specific message.

---

## 🤝 Contributing

Contributions make the open-source community an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

---
