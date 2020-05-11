# Dicebot

Dicebot is a simple chatbot. It has a kind of framework for adding new
verbs, and supports stateful verbs (via GCP Cloud Datastore, or local
files.)

There is a Google Appengine app that implements a LINE webhook:
```
gcloud app deploy ./app/dicebot --project=${YOUR_GCP_PROJECT}
```

There is also a simple CLI, for easier development of new verbs:
```
go run ./cmd/db/*go
```

The main verb is `roll`, for rolling dice. There are some other verbs
that help out with running an RPG in a chat group.
