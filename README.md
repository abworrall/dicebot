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
that help out with running an RPG in a chat group:
```
db help

# Setup things
db bot claim lleldron      # First: claim a character name
db char                    # review your character 'sheet'
db inv list                # review your inventory
db inv stash some rope     # populate your inventory

# Things while playing
db roll 5d6                # roll some dice !
db save vs str             # test your strength !
db inv use 2               # assert you will use item #2
db vow to eat less         # make a vow to do something
db vow list                # look at all the vows made by the party
```
