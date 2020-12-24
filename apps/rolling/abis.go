package main

import "math/rand"

const (
	// BaseArmorClass is the base armor class for any character without armor.
	BaseArmorClass int = 10
	// BaseInitiative is the base initiative for any player.
	BaseInitiative int = 6
	// BaseSpeed is the base speed for any player.
	BaseSpeed int = 30
	// BaseHitDice is the base hit dice for any player
	BaseHitDice int = 10
)

func getModifierFromScore(score int) int {
	switch score {
	case 1:
		return -5
	case 2:
		return -4
	case 3:
		return -4
	case 4:
		return -3
	case 5:
		return -3
	case 6:
		return -2
	case 7:
		return -2
	case 8:
		return -1
	case 9:
		return -1
	case 10:
		return 0
	case 11:
		return 0
	case 12:
		return 1
	case 13:
		return 1
	case 14:
		return 2
	case 15:
		return 2
	case 16:
		return 3
	case 17:
		return 3
	case 18:
		return 4
	case 19:
		return 4
	case 20:
		return 5
	case 21:
		return 5
	case 22:
		return 6
	case 23:
		return 6
	case 24:
		return 7
	case 25:
		return 7
	case 26:
		return 8
	case 27:
		return 8
	case 28:
		return 9
	case 29:
		return 9
	case 30:
		return 10
	}
	return 0
}

// Ability represents abilities for every player.
type Ability struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

// CharacterSheet represents all stats for a player
type CharacterSheet struct {
	Score       *Ability
	Modifier    *Ability
	SavingThrow *Ability
	ArmorClass  int
	Initiative  int
	Speed       int
	HitDice     int
}

// NewAbility creates a new player ability  instance.
func NewAbility(str int, dex int, con int, intel int, wis int, cha int) *Ability {
	return &Ability{
		Strength:     str,
		Dexterity:    dex,
		Constitution: con,
		Intelligence: intel,
		Wisdom:       wis,
		Charisma:     cha,
	}
}

// NewCharacterSheet creates a new player character sheet
// instance.
func NewCharacterSheet(score *Ability) *CharacterSheet {
	modStr := getModifierFromScore(score.Strength)
	modDex := getModifierFromScore(score.Dexterity)
	modCon := getModifierFromScore(score.Constitution)
	modInt := getModifierFromScore(score.Intelligence)
	modWis := getModifierFromScore(score.Wisdom)
	modCha := getModifierFromScore(score.Charisma)
	cs := &CharacterSheet{
		Score:       score,
		Modifier:    NewAbility(modStr, modDex, modCon, modInt, modWis, modCha),
		SavingThrow: NewAbility(modStr, modDex, modCon, modInt, modWis, modCha),
		ArmorClass:  BaseArmorClass + modDex,
		Initiative:  BaseInitiative + modDex,
		Speed:       BaseSpeed,
		HitDice:     BaseHitDice,
	}
	return cs
}

// RollDice rolls a dice with the given number of faces.
func RollDice(faces int) int {
	return rand.Intn(faces + 1)
}
