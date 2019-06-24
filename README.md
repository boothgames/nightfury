# [nightfury](https://github.com/boothgames/nightfury)

Puzzle based game server

## Development

```bash
$ git clone https://github.com/boothgames/nightfury.git
$ cd nightfury
$ make build
$ ./out/nightfury server
```

> specify --log-level `debug` for priting more detailed logging

## Setup

### Games

Create a json file with the game info as save it as `games.json`

```json
[
  {
    "name": "smile",
    "title": "Why so serious?",
    "instruction": "Let's put a smile on that face!",
    "type": "web",
    "mode": "embedded"
  },
  {
    "name": "snakes",
    "title": "Snakes",
    "instruction": "Help the snake to eat the right food",
    "type": "web",
    "mode": "embedded"
  },
  {
    "name": "seeker",
    "title": "Seek the diamonds",
    "instruction": "Ask volunter for a tablet/phone, collect all the diamonds within 60 seconds. Tilt the device in appropriate direction for movement. Beware of the consequences :)",
    "type": "mobile",
    "mode": "external",
    "metadata": {
      "codes": ["1234", "5678", "0987"]
    }
  }
]
```

Upload the games to nightfury

```bash
$ curl -H "Content-Type: application/json" --data @games.json http://localhost:5624/v1/bulk/games

```

### Hints

Create a json file with the hints and save it as `hints.json`

```json
[
  {
    "title": "hint title",
    "tag": ["tag-1", "tag-2"],
    "content": "hint content",
    "takeaway": "hint takeaway"
  }
]
```

Upload the hints to nightfury

```bash
$ curl -H "Content-Type: application/json" --data @hints.json http://localhost:5624/v1/bulk/hints

```

### Contributions

[nightfury](https://github.com/boothgames/nightfury) is an open source project under the Apache 2.0 license, and contributions are gladly welcomed! To submit your changes please open a pull request.