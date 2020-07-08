# battle ship game

Proof of concept battleship game.
Ui written with phaser 3 and typescript.
Backend written in golang with gorilla

## How to build

This should be made more simple but time and all that.

First build the ui by running `npm install` `npm run build` in the ui folder.

Next build the backend by running `go build -o bin/go-battleship ./cmd/go-battleship/main.go` in the root folder

Run the whole thing by executing `./bin/go-battleship`

## How to develop

To develop you need to run two servers. One runs and watches the frontend code, the other one runs the go server

Run the ui by running `npm run watch` in the ui folder. This will run a server on port 100001.
Run the backend by running `go run ./cmd/battleship/` in the root folder. This will run a server on port 100002.

## How to deploy

Make sure the frontend is built, commit to master, and run `git push heroku master`. You can then open the webpage by running `heroku open` and view the logs with `heroku logs --tail`.
