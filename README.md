# battle ship game

Now with google open match.

# how to build

first build the frontend `npm run build`

There are 4 parts involved:

- director: open match director
- frontend: open match frontend
- matchmaking: macthmaking function
- game: runs the actual game logic

All these need their docker build and pushed:

```
  docker build -t triberraar/go-battleship-frontend -f Dockerfile-frontend .
  docker push triberraar/go-battleship-frontend

  docker build -t triberraar/go-battleship-game -f Dockerfile-game .
  docker push triberraar/go-battleship-game

  docker build -t triberraar/go-rps -f Dockerfile-rps .
  docker push triberraar/go-rps

  docker build -t triberraar/go-battleship-director -f Dockerfile-director .
  docker push triberraar/go-battleship-director

  docker build -t triberraar/go-battleship-matchmaking -f Dockerfile-matchmaking .
  docker push triberraar/go-battleship-matchmaking
```

make namespace `kubectl create namespace triberraar-mm`
apply the kubernetes file `kubectl apply -f kubernetes/matchmaking.yaml`

You also need a standard install of open match. Depending on the username length, a player is classified as `noob` or `master`.
noobs go to `localhost:10003`, masters go to `localhost:10003`.
to run in minikube, you need to port forward `kubectl port-forward --namespace triberraar-mm service/go-battleship-frontend 10002:10002` and also the 2 game instances `kubectl port-forward --namespace triberraar-mm service/go-battleship-game 10003:10003` adn `kubectl port-forward --namespace triberraar-mm service/go-battleship-game2 10004:10004`

# Run locally

Edit the following and uncomment the stuff:

- director
- frontend
- matchmaking

forward all internal matchmaking stuff:

```
kubectl port-forward --namespace open-match service/om-frontend 50504:50504
kubectl port-forward --namespace open-match service/om-backend 50505:50505
kubectl port-forward --namespace open-match service/om-query 50503:50503
```

kubectl port-forward --namespace triberraar-mm service/go-battleship-frontend 10002:10002
kubectl port-forward --namespace triberraar-mm service/go-battleship-game 10003:10003
kubectl port-forward --namespace triberraar-mm service/go-battleship-game2 10004:10004
kubectl port-forward --namespace triberraar-mm service/go-rps 10012:10012

run all components:

```
go run ./cmd/go-battleship-director/
go run ./cmd/go-battleship-frontend/
go run ./cmd/go-battleship-game/
go run ./cmd/go-battleship-matchmaking/
```

go to localhost:10002 and play

to remove
`kubectl delete namespace triberraar-mm`

## read logs

kubectl logs -n triberraar-mm --follow pod/go-battleship-director
