# Dicebot

Dicebot is a simple chatbot. It has a kind of framework for adding new
verbs, and supports stateful verbs (via GCP Cloud Datastore, or local
files.)

There is a simple CLI, for easy tweaking and development:
```
go run ./cmd/db/*go
```

There is also a Google Appengine app that implements a LINE webhook:
```
gcloud app deploy ./app/dicebot --project=${YOUR_GCP_PROJECT}
```

To get the bot running inside your LINE groups, you'll need to:
* setup an official "LINE Business ID" account (instructions
[here](https://respond.io/blog/the-ultimate-guide-to-line-for-business/#8crh6)
* get the channel secret for that account
* get a long-lived channel token for that account
* `mv ./pkg/config/myconfig.sample ./pkg/config/myconfig.go`
* edit the new file, putting in your secret & token
* deploy the app
* update your LINE biz acount config, point it to your new webhook (`https://host/line` - the `host` is your new appengine service URL)

A few of the verbs work together to implement a subset of the the 5E
DnD rules, to automate skill checks, attacks, etc.

## Verbs

The main verb is `roll`, for rolling dice. There are some other verbs
that help out with running an RPG in a chat group:
```
db help

## Setup things
#
db bot claim strider       # First: claim a character name

### char

The `db char` command controls your basic character sheet - class,
name, attributes, weapons, etc. So yeah, this is all a bit clunky, but
you only do it once !

```
db char                       # review your character
db char list                  # see what you can set
db char set str 12            # set an attribute (can remove some too)

db char set armor plate       # specify what armor you're wearing
db char set shield 1          # assert you're carrying a shield (set to 0 to unassert)
db char set weapon longsword  # can have multiple; last added is default attack

db inv                        # review your inventory
db inv stash some rope        # add to your inventory
```

Once it knows your class and level, it will list what buffs (class
features) you have. For ones where you get to choose, you can add the
buff by hand.

No support for magic weapons/armor. Yet.

### roll

Need to roll dice ? Oh boy, that's why we're here ! The attribute checks are
all made against those of your character.
```
db roll 4d6+3                      # just roll 'em
db roll vs int                     # save vs int, default difficulty class
db roll vs str >=15 withadvantage  # "hard" (DC=15), but with advantage
db roll vs wis for NOTICING BEAR   # explain what the roll is (gets printed in the result)
```

### spells

The `db spells` command lets you set up your list of prepared spells,
and keeps track of your slots as you cast them.

```
db spells -init wizard int 4 2  # specify your spell slots (and proficiency attribute)

db spells                       # review your prepared spells, and uncast slots
db spells add cure-wounds       # prepare a spell, so you can cast it later
db spells remove web            # remove a spell, then replace it

db spells cast cure-wounds [3]  # cast the spell [at level 3] - consumes a slot
db spells refresh               # have a long rest, and regain all your slots
```

### attack

This verb automates a fair bit of the attack rules. It works on
*copies* of the characters and of the monster definitions. The idea is
to take care of most bookkeeping, but let you do things by hand too.

For a player, the default weapon is the most recent one added via `db
char set weapon`. They can also specify an attack with any of their
other listed weapons. And if they're a spellcaster, they can attack
with the special weapon `magic`, which will do a magic attack roll.

Monsters can use any of their named actions as attacks.

```
## Things the DM will do
#
db attack -reset                   # clear the slate
db attack add orc.4 wolf.4         # add in some monsters (see `db rules monster`)
db attack strider by wolf.2 bite   # second wolf uses `bite` to attack player `strider`

## Things a player can do
#
db attack                          # review the state of things
db attack join                     # will assign an init score

## Weapon attacks
#
db attack wolf.1                   # attack with default weapon
db attack wolf.2 with shortbow     # use a different weapon
db attacl wolf.3 advantage         # you have advantage !
db attack wolf.2 wolf.3 4d6        # do 4d6 damage to two wolves

## Spellcasters can also attack
#
db attack wolf.4 magic-missile     # spellcaster ? do it !

# Cast at a higher level, with damage shared over multiple targets
db attack wolf3 wolf.4 magic missile level 4

## Finally, you can tweak the HP and AC values as you see fit
#
db attack player2 tweak hp +5     # player2 gets healed !
db attack wolf.4  tweak hp -9     # the wolf takes some damage
db attack player3 tweak ac -4     # player3's armor breaks :(
```

One gotcha about magic - it will consume a spellslot from the main
character data object, not the copy in the encounter.

### rules

The bot knows about some spells, equipment, character buffs, and
monsters, courtesy of https://github.com/bagelbits/5e-database. It'll
search for matches if you type a substring.
```
db rules equip sword        # list all the things that sound like a sword
db rules spell cure-wounds  # get details on a particular spell
db rules monster goblin
db rules buff fighter       # list all the fighter buffs
```

## Misc notes

On LINE, it needs to maintain a mapping between the LINE user IDs and
the names that identify characters. The easiest way is to get users to
claim their name once they're logged into LINE: `db bot claim NAME`.
Note that until the user has agreed to the bot's Terms of Use, the
user IDs won't be visible to the bot; the user needs to add the bot to
their friend list and agree to the ToU.

If you are an admin user (see `./config/myconfig.go`), you will need
to claim your name too :) Once you have it, you can then masquerade as
other users, by prefixing all your bot commands with `as USER`, e.g.
`db as foobar roll vs int`.
