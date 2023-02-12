<h1 align="center">go-slack-chat-gpt</h1>
<p align="center">
   <a href='#GoVersion'>
      <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/drkennetz/go-slack-chat-gpt3">
   </a>
    <a href="https://github.com/drkennetz/go-slack-chat-gpt3">
        <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/drkennetz/go-slack-chat-gpt3/ci.yml">
    </a>
    <a href="https://codecov.io/github/drkennetz/go-slack-chat-gpt3" > 
        <img src="https://codecov.io/github/drkennetz/go-slack-chat-gpt3/branch/main/graph/badge.svg?token=8IHKB8J1AN"/> 
    </a>
    <a href="https://github.com/drkennetz/go-slack-chat-gpt3/issues">
        <img src="https://img.shields.io/github/issues/drkennetz/go-slack-chat-gpt3" alt="Issues">
   </a>
</p>
<p align="center">
   <a href='#RepSize'>
      <img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/drkennetz/go-slack-chat-gpt3">
   </a>
    <a href="https://github.com/drkennetz/go-slack-chat-gpt3/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/drkennetz/go-slack-chat-gpt3" /></a>
   <a href='#LastCommit'>
       <img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/drkennetz/go-slack-chat-gpt3">
   </a>
</p>

go-slack-chat-gpt is a simple slack bot server which handles slack app mention events, sending the event to chat-gpt and responding to the channel with chat-gpt's response.

## Table of Contents
- [Quick Start](#Quick-Start)
- [Bot Setup](./example/walkthrough.md)
- [Contributing](#Contributing)
- [Open an Issue](#Issues)
- [Code of Conduct](#Code-of-Conduct)


## Quick Start
Build the binary, add tokens to config, and run!

### Build
```bash
git clone https://github.com/drkennetz/go-slack-chat-gpt3.git
cd go-slack-chat-gpt3 && go build -o ./bin/slackgpt
```

### Config

For a more thorough walk-through of setting up the bot and getting tokens, visit [this detailed doc](./example/walkthrough.md).
```
CGPT_API_KEY=sk-...z7
SLACK_APP_TOKEN=xapp-1-...47
SLACK_BOT_TOKEN=xoxb-...S0
```

### Run

#### Help
```bash
./bin/slackgpt -h
This program is a slack bot that sends mentions to chat-gpt and responds with chat-gpt result

VERSION: development

Usage: slackgpt --config CONFIG

Options:
  --config CONFIG, -c CONFIG
                         config file with slack app+bot tokens, chat-gpt API token
  --help, -h             display this help and exit
  --version              display version and exit
```
#### Run
```
./bin/slackgpt -c ./config.env 
2023/02/01 14:53:19 Config values parsed
socketmode: 2023/02/01 14:53:19 socket_mode_managed_conn.go:258: Starting SocketMode
2023/02/01 14:53:19 Connecting to Slack with Socket Mode...
...
```

## Contributing
Please follow the [Contribution File](./Contribution.md) to contribute to this repo.

## Issues
To submit an issue, select the issue template that most closely 
corresponds with your issue type and submit. Someone will get to you soon!

## Code of Conduct
Please note that go-slack-chat-gpt3 has a [Code of Conduct](./CODE_OF_CONDUCT.md).
By participating in this community, you agree to abide by its rules. 
Failure to abide will result in warning and potentially expulsion from this community.
