# Description

![Github Actions](https://github.com/TIT8/go-api/actions/workflows/fly.yml/badge.svg)

HTTP API that handles, under the hood, the [_e-mail us_](https://triennale-elettronica-polimi.netlify.app/contact/#e-mail-us) section of _Triennale-elettronica-polimi_ website.   

☁️ It's hosted on **[Flyio](https://fly.io/)**. [^1]  
🛡️ Protected by bot via [friendly captcha](https://friendlycaptcha.com/).  
👻 Hidden secrets via enviroment variable.  
⚡ Fast thanks to [Go](https://go.dev/).  

[^1]: Automatically deployed on `master` changes via [Github Actions](https://github.com/TIT8/go-api/actions/workflows/fly.yml).  

## How it works

1. &nbsp; It receives the form data submitted via [javascript POST request](https://github.com/valerionew/triennale-elettronica-polimi/blob/master/layouts/shortcodes/contact.html#L58) from the website.
2. &nbsp; Check if the data body isn't too big and malformed.
3. &nbsp; Check if the captcha submission was successful. If not, send back to the client the captcha error and jump to step 6.
4. &nbsp; It sends message to a Telegram BOT which will write to a <ins>private</ins> channel ([here](https://stackoverflow.com/questions/33858927/how-to-obtain-the-chat-id-of-a-private-telegram-channel) how).
5. &nbsp; Check if the communication with the Telegram API was successful.
6. &nbsp; Write the result on the HTTP header (200 or 406) and on the responsse to the javascript client.
7. &nbsp; The Javascript client will handle the response and inform the user about the operation (inserting text on HTML).

## Testing

❗ **Try it sending a GET request** [^2]
```
curl https://formapi.fly.dev
```

[^2]: You cannot send POST request as you won't send the friendly captcha verification code. 

## Thanks

I would like to express my gratitude to [Marius](https://medium.com/geekculture/how-to-use-go-to-send-telegram-messages-to-your-phone-a819bdf7f35c) for his well-explained blog post, which served as an inspiration for my API.
