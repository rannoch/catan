package domain

import "errors"

type DiceRoller interface {
	Roll() Roll
}

type NumberToken int64

const NumberTokenEmpty = NumberToken(0)

var InvalidNumberToken = errors.New("invalid number token")
var InvalidRoll = errors.New("invalid roll")

func NewNumberToken(n int64) (NumberToken, error) {
	if n < 2 || n > 12 {
		return 0, InvalidNumberToken
	}

	return NumberToken(n), nil
}

func MustGetNumberToken(n int64) NumberToken {
	token, err := NewNumberToken(n)
	if err != nil {
		panic(err)
	}

	return token
}

type D6Roll int64

const (
	D6Roll1 = D6Roll(1)
	D6Roll2 = D6Roll(2)
	D6Roll3 = D6Roll(3)
	D6Roll4 = D6Roll(4)
	D6Roll5 = D6Roll(5)
	D6Roll6 = D6Roll(6)
)

func NewD6Roll(n int64) (D6Roll, error) {
	if n < 1 || n > 6 {
		return 0, InvalidRoll
	}

	return D6Roll(n), nil
}

func MustGetD6Roll(n int64) D6Roll {
	roll, err := NewD6Roll(n)
	if err != nil {
		panic(err)
	}
	return roll
}

type Roll struct {
	d6Roll1 D6Roll
	d6Roll2 D6Roll
}

func NewRoll(d6Roll1 D6Roll, d6Roll2 D6Roll) Roll {
	return Roll{d6Roll1: d6Roll1, d6Roll2: d6Roll2}
}

func (r Roll) IsRobber() bool {
	return r.Value() == 7
}

func (r Roll) Value() int64 {
	return int64(r.d6Roll1 + r.d6Roll2)
}

func (r Roll) NumberToken() NumberToken {
	return MustGetNumberToken(r.Value())
}
