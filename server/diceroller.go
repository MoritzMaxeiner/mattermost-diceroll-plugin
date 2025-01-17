package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/model"
	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"
)

type DiceRoller struct {
	source rand.Source

	distributions map[int]*distuv.Categorical

	formulaPattern *regexp.Regexp

	filterPattern *regexp.Regexp

	imageLinks map[DiceSystem]map[int]string
}

func NewDiceRoller() *DiceRoller {
	r := &DiceRoller{}

	r.source = NewCryptoSeededMT64()

	r.distributions = make(map[int]*distuv.Categorical)

	r.formulaPattern = regexp.MustCompile(`^` +
		`(?P<num_dice>[0-9]+)?d(?P<type>%|F|AE|[0-9]+)` +
		`(e(?P<explode_op>>=|<=)(?P<explode_tn>-?[0-9]+))?` +
		`((?P<filters>(?:(?:kh|kl|dh|dl)[1-9][0-9]*)+))?` +
		`((?P<total_modifier_op>[+\-*/])(?P<total_modifier_val>[0-9]+))?` +
		`((?P<success_op>>=|<=)(?P<success_tn>-?[0-9]+))?` +
		`$`)

	r.filterPattern = regexp.MustCompile(`(k[hl]|d[hl])([1-9][0-9]*)`)

	r.imageLinks = make(map[DiceSystem]map[int]string)
	r.imageLinks[dsAetherium] = make(map[int]string)

	return r
}

func (r *DiceRoller) RollSingle(sides int) int {
	if dist, ok := r.distributions[sides]; ok {
		return 1 + int(dist.Rand())
	} else {
		weights := make([]float64, sides)
		for i := range weights {
			weights[i] = 1
		}

		newDist := distuv.NewCategorical(weights, r.source)
		dist = &newDist
		r.distributions[sides] = dist
		return 1 + int(dist.Rand())
	}
}

type DiceSystem int

const (
	dsStandard DiceSystem = iota
	dsFudge
	dsAetherium
)

type AetheriumSymbol int

const (
	aesDisruption AetheriumSymbol = iota
	aesSwitch
	aesChip
	aesShort
	aesCrash
)

func ParseDiceType(type_ string) (DiceSystem, int) {
	if type_ == "%" {
		return dsStandard, 100
	} else if type_ == "F" {
		return dsFudge, 3
	} else if type_ == "AE" {
		return dsAetherium, 12
	} else {
		num_sides, _ := strconv.Atoi(type_)
		return dsStandard, num_sides
	}
}

func (r *DiceRoller) FormatDiceResult(system DiceSystem, roll int) string {
	if system == dsFudge {
		switch roll {
		case -1:
			return "**-**"
		case 0:
			return "&nbsp;&nbsp;"
		case 1:
			return "**+**"
		}
	} else if system == dsAetherium {
		if link, ok := r.imageLinks[dsAetherium][roll]; ok {
			return link
		}
	}

	// NOTE: system == dsStandard or any other non-special
	return fmt.Sprintf("**%v**", roll)
}

func LowestDieOutcome(system DiceSystem, numSides int) int {
	if system == dsFudge {
		return -1
	}

	return 1
}

func HighestDieOutcome(system DiceSystem, numSides int) int {
	if system == dsFudge {
		return 1
	}

	return numSides
}

func (r *DiceRoller) RollNotation(notation string) *model.SlackAttachment {
	warnings := ""

	fields := []*model.SlackAttachmentField{}

	if r.formulaPattern.MatchString(notation) {
		matches := make(map[string]string)

		values := r.formulaPattern.FindStringSubmatch(notation)
		groups := r.formulaPattern.SubexpNames()

		for idx := 1; idx < len(groups); idx++ {
			matches[groups[idx]] = values[idx]
		}

		numDice, err := strconv.Atoi(matches["num_dice"])
		if err != nil {
			numDice = 1
		}

		if numDice < 1 {
			warnings += fmt.Sprintf("⚠️ %v is too few dice, rolling 1.\n", numDice)
			numDice = 1
		}

		if numDice > 100 {
			warnings += fmt.Sprintf("⚠️ %v is too many dice, rolling 100.\n", numDice)
			numDice = 100
		}

		system, numSides := ParseDiceType(matches["type"])

		if numSides < 2 {
			warnings += fmt.Sprintf("⚠️ %v is too few sides, rolling %dd2.\n", numSides, numDice)
			numSides = 2
		}

		var rolls []int
		for idx := 0; idx < numDice; idx++ {
			roll := r.RollSingle(numSides)

			if system == dsFudge {
				roll -= 2
			}

			rolls = append(rolls, roll)
		}

		// Optionally explode dice
		explodeDice := len(matches["explode_op"]) > 0
		numExplosions := 0
		if explodeDice {
			op := matches["explode_op"]
			targetNumber, _ := strconv.Atoi(matches["explode_tn"])

			lowest := LowestDieOutcome(system, numSides)
			highest := HighestDieOutcome(system, numSides)
			if op == ">=" && targetNumber <= lowest {
				warnings += fmt.Sprintf("⚠️ %v is too low an explosion target number, rolling with e>=%v instead.\n", targetNumber, lowest+1)
				targetNumber = lowest + 1
			} else if op == "<=" && targetNumber >= highest {
				warnings += fmt.Sprintf("⚠️ %v is too high an explosion target number, rolling with e<=%v instead.\n", targetNumber, highest-1)
				targetNumber = highest - 1
			}

			shouldExplode := func(roll int) bool {
				switch op {
				case ">=":
					return roll >= targetNumber
				case "<=":
					return roll <= targetNumber
				}

				return false
			}

			for idx := range rolls {
				if shouldExplode(rolls[idx]) {
					numExplosions += 1
				}
			}

			for idx := 0; idx < numExplosions; idx += 1 {
				roll := r.RollSingle(numSides)
				rolls = append(rolls, roll)
				if shouldExplode(roll) {
					numExplosions += 1
				}
			}
		}

		// Optionally filter highest/lowest dice
		useRoll := make([]bool, len(rolls))
		filtering := len(matches["filters"]) > 0
		if filtering {
			filtered := ArgSort(rolls)

			match := r.filterPattern.FindAllStringSubmatch(matches["filters"], -1)
			for idx := 0; idx < len(match); idx++ {
				kind := match[idx][1]
				numDice, _ := strconv.Atoi(match[idx][2])
				numDice = min(max(0, numDice), len(filtered))

				switch kind {
				case "dl":
					filtered = filtered[numDice:]
				case "kl":
					filtered = filtered[:numDice]
				case "dh":
					filtered = filtered[:len(filtered)-numDice]
				case "kh":
					filtered = filtered[len(filtered)-numDice:]
				}
			}

			for idx := 0; idx < len(filtered); idx++ {
				useRoll[filtered[idx]] = true
			}
		} else {
			for idx := 0; idx < len(rolls); idx++ {
				useRoll[idx] = true
			}
		}

		// Format rolled dice
		formatSystem := system
		if system == dsFudge && filtering {
			formatSystem = dsStandard
		}
		formatDie := func(idx int) string {
			str := r.FormatDiceResult(formatSystem, rolls[idx])
			if filtering {
				if useRoll[idx] {
					return str
				} else {
					return fmt.Sprintf("~~%v~~", str)
				}
			} else {
				return str
			}
		}

		rollField := ""
		multiLine := false
		for idx := 0; idx < len(rolls); idx++ {
			if idx > 0 && (idx%10) == 0 {
				rollField += "|\n"
				if !multiLine {
					multiLine = true
					rollField += "|-|\n"
				}
			}
			rollField += "|" + formatDie(idx)
		}
		rollField += "|\n"
		if !multiLine {
			rollField += "|-|\n||"
		}
		fields = append(fields, &model.SlackAttachmentField{
			Title: notation,
			Value: rollField,
		})

		if explodeDice {
			fields = append(fields, &model.SlackAttachmentField{
				Title: "Explosions",
				Value: fmt.Sprintf("%v", numExplosions),
			})
		}

		// Optionally aggregate the total
		aggregateTotal := system == dsStandard || system == dsFudge
		if aggregateTotal {
			modifyTotal := len(matches["total_modifier_op"]) > 0

			total := 0
			for idx := 0; idx < len(rolls); idx++ {
				if useRoll[idx] {
					total += rolls[idx]
				}
			}

			if modifyTotal {
				modifiedTotal := float64(total)
				modifier, _ := strconv.Atoi(matches["total_modifier_val"])

				formatString := "%v"
				switch matches["total_modifier_op"] {
				case "+":
					modifiedTotal += float64(modifier)
				case "-":
					modifiedTotal -= float64(modifier)
				case "*":
					modifiedTotal *= float64(modifier)
				case "/":
					modifiedTotal /= float64(modifier)
					formatString = "%.2f"
				}

				fields = append(fields, &model.SlackAttachmentField{
					Title: "Total",
					Value: fmt.Sprintf(formatString, modifiedTotal),
				})
			} else {
				fields = append(fields, &model.SlackAttachmentField{
					Title: "Total",
					Value: fmt.Sprintf("%v", total),
				})
			}
		}

		// Optionally aggregate the successes
		aggregateSuccesses := len(matches["success_op"]) > 0
		if aggregateSuccesses {
			op := matches["success_op"]
			targetNumber, _ := strconv.Atoi(matches["success_tn"])

			numSuccesses := 0
			for idx := range rolls {
				if !useRoll[idx] {
					continue
				}

				switch op {
				case ">=":
					if rolls[idx] >= targetNumber {
						numSuccesses += 1
					}
				case "<=":
					if rolls[idx] <= targetNumber {
						numSuccesses += 1
					}
				}
			}

			fields = append(fields, &model.SlackAttachmentField{
				Title: "Successes",
				Value: fmt.Sprintf("%v", numSuccesses),
			})
		}

		if system == dsAetherium {
			numSymbols := map[AetheriumSymbol]int{
				aesDisruption: 0,
				aesSwitch:     0,
				aesChip:       0,
				aesShort:      0,
				aesCrash:      0,
			}

			for idx := range rolls {
				if useRoll[idx] {
					roll := rolls[idx]

					if roll == 5 || roll == 9 || roll == 11 || roll == 12 {
						numSymbols[aesDisruption] += 1
					}

					switch {
					case 1 <= roll && roll <= 5:
						numSymbols[aesSwitch] += 1
					case 6 <= roll && roll <= 9:
						numSymbols[aesChip] += 1
					case 10 <= roll && roll <= 11:
						numSymbols[aesShort] += 1
					case roll == 12:
						numSymbols[aesCrash] += 1
					}
				}
			}

			symbolField := "|Disruption|Switch|Chip|Short|Crash|\n|-|\n"
			symbolField += fmt.Sprintf("|%v", numSymbols[aesDisruption])
			symbolField += fmt.Sprintf("|%v", numSymbols[aesSwitch])
			symbolField += fmt.Sprintf("|%v", numSymbols[aesChip])
			symbolField += fmt.Sprintf("|%v", numSymbols[aesShort])
			symbolField += fmt.Sprintf("|%v", numSymbols[aesCrash])

			fields = append(fields, &model.SlackAttachmentField{
				Title: "Symbols",
				Value: symbolField,
			})
		}
	} else {
		warnings += "⚠️ I have no idea what to do."
	}

	return &model.SlackAttachment{
		Fields: fields,
		Footer: warnings,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
