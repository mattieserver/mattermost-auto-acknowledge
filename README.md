# mattermost-auto-acknowledge

Automatically react to all messages send by a person in a certain team and channel in mattermost. 

## Build

```bash
go build cmd/mattermost-auto-acknowledge/acknowledge.go
```

## Configure

Create a folder `configs` and create a file `settings.json`.  
It should look like this: 
```json
{
    "mattermost": {
        "base_url": "http://localhost:8000",
        "username": "test",
        "password": "test",
        "team": "test",
        "channel": "test",
        "user": "test",
        "sub_messages": false
    }
}
```

## Run

```bash
./acknowledge
```