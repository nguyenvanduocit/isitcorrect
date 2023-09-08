# Is it correct?

Yesterday, I chated with a friend in English, and ... you know, it's hard.

Writing in code is easier than Englisn. So, I wrote it in an hour and a half, so it's not perfect. If you find any bugs, please report them in the issues section.

## Structure

Currently, using OpenAI is the fastest way to get results, but I don't want to pay. Meanwhile, ChatGPT is completely free. So, I've developed a Chrome extension to inject custom scripts into ChatGPT and then connect to a WebSocket server, essentially acting as an API proxy.

On the other side, I use Golang to create an application that provides a UI, clipboard monitor, and WebSocket server.

## Usage

First, you should build the extension using the command: "pnpm run build." After that, load the extension into Chrome.

Secondly, construct the app using the following command: "wails build." It is essential to install Wails beforehand.

Then, log in to ChatGPT and keep the tab open. There's no need to make it active.

Finally, execute the application and press the button to establish a connection with the WebSocket server.

The app have only two options:

Let's Hoang do it for you

1. Enable the 'Monitor clipboard' option, and the app will continuously monitor your clipboard, sending its content to the WebSocket server. It will check the clipboard every 1 second.
2. Always on top: Enable the 'Always on top' option to ensure that the app remains above other windows at all times.

p/s: The English on this Readme is not good, because it was created before the app was done. :)))

