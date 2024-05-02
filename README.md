# Description

![Github Actions](https://github.com/TIT8/go-api/actions/workflows/fly.yml/badge.svg)

HTTP API that handles, under the hood, the [_e-mail us_](https://triennale-elettronica-polimi.netlify.app/contact/#e-mail-us) section of _Triennale-elettronica-polimi_ website.   

☁️ It's hosted on **[Flyio](https://fly.io/)**. [^1]  
🛡️ Protected by bot via [friendly captcha](https://friendlycaptcha.com/).  
👻 Hidden secrets via enviroment variable.  
⚡ Fast thanks to [Go](https://go.dev/).  
📬 Communicate with the Telegram API. [^2]

[^1]: Automatically deployed on `master` changes via [Github Actions](https://github.com/TIT8/go-api/actions/workflows/fly.yml).  
[^2]: The [Telegram API](https://core.telegram.org/)

## How it works

1. &nbsp; It receives the form data submitted via [javascript POST request](https://github.com/valerionew/triennale-elettronica-polimi/blob/master/layouts/shortcodes/contact.html#L58) from the website.
2. &nbsp; Check if the data body isn't too big and malformed.
3. &nbsp; Check if the captcha submission was successful. If not, send back to the client the captcha error and jump to step 6.
4. &nbsp; It sends message to a Telegram BOT which will write to a <ins>private</ins> channel ([here](https://stackoverflow.com/questions/33858927/how-to-obtain-the-chat-id-of-a-private-telegram-channel) how).
5. &nbsp; Check if the communication with the Telegram API was successful.
6. &nbsp; Write the result on the HTTP header (200 or 406) and on the responsse to the javascript client.
7. &nbsp; The Javascript client will handle the response and inform the user about the operation (inserting text on HTML).

## Why redirect the HTML form to this API if JavaScript can handle HTTP requests with ease?

1. Because the [Polimi website](https://github.com/valerionew/triennale-elettronica-polimi) is open source, it's important for future generations of students to have access to it and the ability to make changes. However, in JavaScript, how can we hide secret variables while maintaining an open-source project? These secret variables are crucial for accessing friendly captcha protection and the Telegram chat where messages are delivered. If these secrets are exploited, it could lead to misuse of the Telegram chat and render it unusable.

   You can refer to discussions [here](https://stackoverflow.com/questions/28890783/how-do-i-hide-a-variable-value-in-javascript) and [here](https://stackoverflow.com/questions/8520626/how-it-is-possible-to-not-expose-you-secret-key-with-a-javascript-oauth-library), which focus on the server-side approach. 

   But it's important to remember that Telegram is just an API endpoint, and it's entirely possible to send messages to chats from the browser as well (see [here](https://stackoverflow.com/questions/73084236/send-message-to-telegram-through-html-form-using-javascript) for an example). The same applies to the verification API of friendly captcha, which you can find [here](https://docs.friendlycaptcha.com/#/verification_api).

2. This API can be extended with additional features, such as the ability to directly receive files and create pull requests to the GitHub repository, something that cannot be achieved in the browser.

## Testing

❗ **Try it sending a GET request** [^3]
```
curl https://formapi.fly.dev
```

Or if you are on a web browser, just [click me](https://formapi.fly.dev).

[^3]: You cannot send POST request as you won't send the friendly captcha verification code. 


## Thanks

I would like to express my gratitude to [Marius](https://medium.com/geekculture/how-to-use-go-to-send-telegram-messages-to-your-phone-a819bdf7f35c) for his well-explained blog post, which served as an inspiration for my API.
