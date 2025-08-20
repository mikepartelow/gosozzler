package sozzler

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Quantity struct {
	s string
	f float64
}

func (q Quantity) String() string {
	parts := strings.Split(q.s, "/1")
	if parts[0] == "0" {
		return ""
	}
	return parts[0]
}

func (q Quantity) Float() float64 {
	return q.f
}

func (q Quantity) Scale(factor int) Quantity {
	q.f *= float64(factor)
	q.s = stringer(q.f)
	return q
}

func stringer(f float64) string {
	if float64(f) == 0 {
		return ""
	}

	fractionMap := map[float64]string{
		0.5:   "1/2", // "½"
		0.25:  "1/4", // "¼"
		0.75:  "3/4", // "¾"
		0.125: "1/8",
	}

	intPart := int(f)
	fracPart := f - float64(intPart)

	if fancy, ok := fractionMap[fracPart]; ok {
		if intPart == 0 {
			return fancy
		}
		return fmt.Sprintf("%d %s", intPart, fancy)
	}

	return fmt.Sprint(f)
}

func (q Quantity) MarshalYAML() (interface{}, error) {
	if q.s == "" {
		panic("wtf")
	}
	return q.s, nil
}

func (q *Quantity) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("quantity: unsupported YAML type %q", value.Tag)
	}

	v, err := parseFraction(s)
	if err != nil {
		return fmt.Errorf("quantity: %w", err)
	}
	*q = Quantity{s: s, f: v}
	return nil
}

func ParseQuantity(s string) (*Quantity, error) {
	f, err := parseFraction(s)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse quantity %q: %w", s, err)
	}
	return &Quantity{s: s, f: f}, nil
}

func parseFraction(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	if !strings.Contains(s, "/") {
		s += "/1"
	}

	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid fraction %q", s)
	}
	num, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	den, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil || den == 0 {
		return 0, fmt.Errorf("invalid fraction %q", s)
	}
	return num / den, nil
}
