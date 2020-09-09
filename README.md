# battle ship game

Now with google open match and agones.

# how to build

Start a minikube:
`minikube start --kubernetes-version v1.16.0`

Start open match:
`kubectl apply --namespace open-match -f kubernetes/open-match-core.yaml`

Start agones:
`kubectl create namespace agones-system`
`kubectl apply -f kubernetes/agones.yaml`

Install custom components for open match (see below to build and push):
`kubectl apply -f kubernetes/matchmaking.yaml`

Start up a game server fleet with autoscaler:
`kubectl apply -f kubernetes/fleet.yaml`

first build the frontend `npm run build`

There are 4 parts involved:

- director: open match director
- frontend: open match frontend
- matchmaking: macthmaking function
- game: runs the actual game logic

All these need their docker build and pushed:

```
  docker build -t triberraar/frontend -f Dockerfile-frontend .
  docker push triberraar/frontend

  docker build -t triberraar/battleship -f Dockerfile-battleship .
  docker push triberraar/battleship

  docker build -t triberraar/rps -f Dockerfile-rps .
  docker push triberraar/rps

  docker build -t triberraar/director -f Dockerfile-director .
  docker push triberraar/director

  docker build -t triberraar/matchmaking -f Dockerfile-matchmaking .
  docker push triberraar/matchmaking
```
