# Mattermost DiceRoll Plugin

This plugin enables users to roll dice using the /roll command.

## Usage

`/roll FORMULA...`

Roll at most 10 [dice algebra](https://en.wikipedia.org/wiki/Dice_notation) `FORMULA`(s).

A single `FORMULA` has the canonical form `[N]dT[EXPLODE][FILTER...][TOTAL][SUCCESS]` and is evaluated from left to right, where

-   `N` is the _optional_ number of dice to roll (default: 1)
-   `T` is the type of dice to roll:
    -   A number: Roll `T`-sided dice (`T >= 2`) and aggregate the total
    -   `%`: Roll d100 (_percentile_) dice and aggregate the total
    -   `F`: Roll [Fudge](https://en.wikipedia.org/wiki/Fudge_%28role-playing_game_system%29) dice (equiprobable die outcomes {`plus`, `minus`, `blank`}) and aggregate the total
    -   `AE`: Roll Aetherium dice (d12 with outcomes {`switch` on (1-5), `chip` on (6-9), `short` on (10-11), `crash` on (12)} × {`disruption` on (5,9,11,12), `blank` otherwise}) and aggregate the symbols
-   `EXPLODE` enables _optional_ die explosion:
    -   `e>=K`: Roll 1 additional die for each die outcome greater than or equal to `K`
    -   `e<=K`: Roll 1 additional die for each die outcome less than or equal to `K`
-   Each _optional_ `FILTER` (sub)selects the dice used for aggregation:
    -   `dlK`: Drops the `K` lowest dice
    -   `klK`: Keeps only the `K` lowest dice
    -   `dhK`: Drops the `K` highest dice
    -   `khK`: Keeps only the `K` highest dice
-   `TOTAL` _optionally_ applies a modifier on the total:
    -   `+K`: Adds `K` to the total
    -   `-K`: Subtracts `K` from the total
    -   `*K`: Multiplies the total by `K`
    -   `/K`: Divides the total by `K`
-   `SUCCESS` _optionally_ determines the number of dice meeting a target number (successes)
    -   `s>=K`: Die outcomes greater than or equal to `K` are successes
    -   `s<=K`: Die outcomes less than or equal to `K` are successes

## Getting Started

Build the plugin:

```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/sh.ucw.mattermost-diceroll-plugin-${version}.tar.gz
```
