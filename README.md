# Noty
This is a basic Discord bot that uses channels to handle todo lists.

## Basic Commands
- `/start` delete all messages in the channel and shows todo list. If list is empty it shows a little tip.
- `/todo add <text>` create a new todo
- `/todo done <id>` mark a todo as done
- `/todo notdone <id>` mark a todo as not done
- `/todo delete <id>` delete a todo
- `/todo update <id> <text>` update a todo
- `/todo reset` delete all todos

## Configuration
Noty need a configuration file in order to work properly: `config.yml`
The structure is similar the one below.
```YAML
token: "YourMagic.Discord.Token"
channels:
    - "ChannelIDs"
    - "WhereNoty"
    - "WillWork"
```

## Docker
Noty works well with docker too! Just a few steps and it will run on a docker container:
1. Create an image called "noty" `docker build -t noty .`
2. Run the image `docker run -d -v /your/absolute/path:/app/data noty`
3. Enjoy!

### Tips
- [Find Discord IDs](https://support.discordapp.com/hc/it/articles/206346498-Come-posso-trovare-l-ID-del-mio-server)
- Type `docker ps` to see all containers running