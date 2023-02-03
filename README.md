<h1 align="center">go-slack-chat-gpt</h1>
<p align="center">
    <a href="https://github.com/drkennetz/go-slack-chat-gpt3/graphs/contributors" alt="Contributors">
        <img src="https://img.shields.io/github/contributors/drkennetz/go-slack-chat-gpt3.svg" /></a>
    <a href="https://github.com/drkennetz/go-slack-chat-gpt3/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/drkennetz/go-slack-chat-gpt3" /></a>
    <a href="https://github.com/drkennetz/go-slack-chat-gpt3">
        <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/drkennetz/go-slack-chat-gpt3/ci.yml">
    </a>
<!--    <a href="https://github.com/drkennetz/go-slack-chat-gpt3">
        <img src="https://img.shields.io/github/workflow/status/drkennetz/go-slack-chat-gpt3/Service%20Testing/main" alt="Service Testing Status">
    </a> -->
    <a href="#sponsors" alt="Sponsors on Open Collective">
        <img alt="GitHub Sponsors" src="https://img.shields.io/github/sponsors/drkennetz"></a>
    <a href="https://github.com/drkennetz/go-slack-chat-gpt3/issues">
        <img src="https://img.shields.io/github/issues/drkennetz/go-slack-chat-gpt3" alt="Issues">
   </a>
   <a href="#fork">
        <img src="https://img.shields.io/github/forks/drkennetz/go-slack-chat-gpt3?label=Fork" alt="Fork">
   </a>
   <a href='#LastCommit'>
       <img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/drkennetz/go-slack-chat-gpt3">
   </a>
   <a href='#GoVersion'>
      <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/drkennetz/go-slack-chat-gpt3">
   </a>
   <a href='#RepSize'>
      <img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/drkennetz/go-slack-chat-gpt3">
   </a>
</p>


![example](./example/slack-gpt-bot.gif)

go-slack-chat-gpt is a simple slack bot server which handles slack app mention events, sending the event to chat-gpt and responding to the channel with chat-gpt's response.

## Table of Contents
- [Quick Start](#Quick-Start)
- [Setup](#Setup)
- [Contributing](#Contributing)
- [Open an Issue](#Issues)



## Quick Start
Build the binary, add tokens to config, and run!

### Config
The config is parsed by [viper](https://github.com/spf13/viper) and must conform to the [viper supported config types](https://github.com/spf13/viper#reading-config-files).

A walk through of getting a chat-gpt token, setting up a slack bot and giving it proper permissions, and getting tokens will be discussed in detail further down the README. This section will just ensure that you build the tool and have it functional on the CLI.
```
CGPT_API_KEY=sk-...z7
SLACK_APP_TOKEN=xapp-1-...47
SLACK_BOT_TOKEN=xoxb-...S0
```
Note: All slack app tokens start with `xapp` and all slack bot tokens start with `xoxb`. 
```bash
# build
go build -o ./bin/slackgpt

# run help
./bin/slackgpt -h
This program is a slack bot that sends mentions to chat-gpt and responds with chat-gpt result

VERSION: development

Usage: slackgpt --config CONFIG

Options:
  --config CONFIG, -c CONFIG
                         config file with slack app+bot tokens, chat-gpt API token
  --help, -h             display this help and exit
  --version              display version and exit

# and run prog with proper config + tokens
./bin/slackgpt -c ./.env 
2023/02/01 14:53:19 Config values parsed
socketmode: 2023/02/01 14:53:19 socket_mode_managed_conn.go:258: Starting SocketMode
2023/02/01 14:53:19 Connecting to Slack with Socket Mode...
...
```

# Setup
In order to run the slack-chat-gpt app, you need a chat-gpt api key, a slack bot with correct OAuth permissions and a bot token, and a slack app token.
While the chat-gpt api key is fairly straightforward, the slack bot is a bit more challenging so
a walk-through will be included below.

The config is parsed by [viper](https://github.com/spf13/viper) and must conform to the [viper supported config types](https://github.com/spf13/viper#reading-config-files).

## Chat-GPT API Key
For the chat-gpt API key, sign up for an account and copy the key: https://platform.openai.com/account/api-keys.
Once you have this key, add it to your config file as:
```CGPT_API_KEY=<mykey>```

## Slack App Creation
1. Sign into the workspace for which you want to create the application for
2. Go to api.slack.com and click "Create an App"

![create an app](./example/create_an_app.png)
3. A pop-up will appear asking if you'd like to create from scratch, or from a manifest. Click "From scratch".
4. Give the app a name. Mine will be called "slackgpt".
5. Select the correct workspace. Mine is called "app-development".
6. Click Create App

## Setting up app
1. If you've done the above, you will be taken to the Basic Information page of your app.
2. From here, click the "Socket Mode" tab under "Settings" in the left panel.

![socket mode](./example/socket_mode.png)

3. Click the button to "Enable Socket Mode"

![enable socket mode](./example/enable_socket_mode.png)

4. Under "Connect using Socket Mode", click "App Level Token". 
5. Click the "Add features and functionality" dropdown and click "Event Subscriptions".

![event subscriptions](./example/event_subscriptions.png)

6. This will take you to a page that will allow you to "Enable Events". Click it to "On".
7. This will trigger some options. Click the "Subscribe to bot events" dropdown and click "Add Bot User Event".

![bot user event](./example/add_bot_user_event.png)

8. Select "app_mention". This will enable your bot to access mentions of your bot. Your bot events section should now look like this:

![app mention event](./example/app_mention_bot_event.png)

9. Then, we must save changes - click the "Save Changes" button on the bottom right.
10. On the left sidebar under "Features", navigate to "OAuth & Permissions".
11. Scroll down the page to "Scopes", where you should already have an "app_mentions:read" scope applied under "Bot Token Scopes".

![oauth pre](./example/oauth_scope_pre.png)

12. We need to add "chat:write", as a scope - so select that and your new Bot Token Scopes should look like this:

![oauth post](./example/oauth_scope_post.png)

13. Finally, we need to collect our 2 tokens. Scroll up the OAuth page, and click "Install to Workspace". If you are an admin, this will happen automatically. If not, it will send a request.

![admin view](./example/admin_view.png)

14. Once it is installed, you will be redirected back to OAuth, where you will see a bot token.
15. Copy that token starting with `xoxb-` to your config file under the key: ```SLACK_BOT_TOKEN=<bot token>```
16. Finally, go back to "Basic Information" on the left sidebar and scroll down past "App Credentials" to "App-Level Tokens".

![app token](./example/app_level_tokens.png)

17. Create one by clicking "Generate Token and Scopes". **Note**: Unlike your bot token, you will not be able to see this text again after you leave, so be sure to copy right away!
18. A pop-up will ask for a token name and a scope. Name your token, and add the scope "connections:write".

![app token fields](./example/app_tokens_values.png)

19. Click "Generate", and copy the token starting with `xapp-` to your config file under the key: ```SLACK_APP_TOKEN=```.

20. This is it! Once you've made it here, you should have a config file with the following key=value fields:
```
CGPT_API_KEY=<value>
SLACK_BOT_TOKEN=xoxb-...
SLACK_APP_TOKEN=xapp-...
```

21. Now go to your slack workspace, and either add the app to the existing channel, or create a new channel to add the app to.

![add channel](./example/slack-channel.png)

22. Add people, or you can skip, but we need to add the app. We can do this by directly mentioning the app's handle. The app's handle will be what you named the app in the beginning. Mine was `slackgpt`

![add bot to channel](./example/invite_bot_to_channel.png)

23. Now that they are added to the channel, and you have your config with your tokens, build and run the app as seen in the Quick Start section at the top:

```
go build -o ./bin/slackgpt

## config.txt is my config
## Note, the filename must end in .env to use the config format above
./bin/slackgpt -c ./config.env

socketmode: 2023/02/02 20:14:28 socket_mode_managed_conn.go:258: Starting SocketMode
2023/02/02 20:14:28 Connecting to Slack with Socket Mode...
api: 2023/02/02 20:14:28 socket_mode.go:30: ...
...
```

24. Try out in slack:

![example](./example/example_chat.png)

## Contributing
Please follow the [Contribution File](./Contribution.md) to contribute to this repo.

## Issues
To submit an issue, select the issue template that most closely corresponds with your issue type and submit. Someone will get to you soon!