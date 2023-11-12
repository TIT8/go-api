# go-api

![Github Actions](https://github.com/TIT8/go-api/actions/workflows/fly.yml/badge.svg)

Go api for Triennale-elettronica-polimi "e-mail us" [section](https://triennale-elettronica-polimi.netlify.app/contact/#e-mail-us) with [friendly captcha](https://friendlycaptcha.com/) protection.

## How it works

1. &nbsp; It receives the form data submitted via [javascript POST request](https://github.com/valerionew/triennale-elettronica-polimi/blob/master/layouts/shortcodes/contact.html#L58) from the website.
2. &nbsp; Check if the data body isn't too big and malformed.
3. &nbsp; Check if the captcha submission was successful. If not, send back to the client the captcha error and jump to step 6.
4. &nbsp; It sends message to a Telegram BOT which will write to a <ins>private</ins> channel ([here](https://stackoverflow.com/questions/33858927/how-to-obtain-the-chat-id-of-a-private-telegram-channel) how).
5. &nbsp; Check if the communication with the Telegram API was successful.
6. &nbsp; Write the result on the HTTP header (200 or 406) and on the responsse to the javascript client.
7. &nbsp; The Javascript client will [handle the response](https://github.com/valerionew/triennale-elettronica-polimi/blob/master/layouts/shortcodes/contact.html#L73) and inform the user about the operation (inserting text on HTML).

It's hosted on **[Flyio](https://fly.io/)** [^1], automatically deployed on master changes via [Github Actions](https://github.com/TIT8/go-api/actions/workflows/fly.yml).    

[^1]: All the secret variables are stored as enviroment variables on Fliyio.

**Try it sending a GET request**: `curl https://formapi.fly.dev` 🪃

## Thanks

I would like to express my gratitude to [Marius](https://medium.com/geekculture/how-to-use-go-to-send-telegram-messages-to-your-phone-a819bdf7f35c) for his well-explained blog post, which served as an inspiration for my API.
