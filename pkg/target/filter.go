package target

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/anarcher/kroller/pkg/resource"
)

// The filter and format is inspired by https://github.com/grafana/tanka/tree/master/pkg/process

func Filter(list resource.RolloutList, exprs Matchers) resource.RolloutList {
	out := make(resource.RolloutList, 0, len(list))

	for _, r := range list {
		target := fmt.Sprintf("%s/%s/%s", r.Namespace(), r.Kind(), r.Name())
		if !exprs.MatchString(target) {
			continue
		}
		if exprs.IgnoreString(target) {
			continue
		}
		out = append(out, r)
	}

	return out
}

type Matcher interface {
	MatchString(string) bool
}

type Ignorer interface {
	IgnoreString(string) bool
}

type Matchers []Matcher

func (e Matchers) MatchString(s string) bool {
	b := false
	i := 0
	for _, exp := range e {
		_, ok := exp.(Ignorer)
		if ok {
			continue
		}
		b = b || exp.MatchString(s)
		i++
	}

	// when no matches, return true
	if i == 0 {
		return true
	}

	return b
}

func (e Matchers) IgnoreString(s string) bool {
	b := false
	for _, exp := range e {
		i, ok := exp.(Ignorer)
		if !ok {
			continue
		}
		b = b || i.IgnoreString(s)
	}

	return b
}

func RegExps(rs []*regexp.Regexp) Matchers {
	xprs := make(Matchers, 0, len(rs))
	for _, r := range rs {
		xprs = append(xprs, r)
	}

	return xprs
}

func StrExps(strs ...string) (Matchers, error) {
	exps := make(Matchers, 0, len(strs))
	for _, raw := range strs {
		s := fmt.Sprintf(`(?i)^%s$`, strings.TrimPrefix(raw, "!"))

		var exp Matcher
		exp, err := regexp.Compile(s)
		if err != nil {
			return nil, ErrBadExpr{err}
		}

		if strings.HasPrefix(raw, "!") {
			exp = NegMatcher{exp: exp}
		}
		exps = append(exps, exp)
	}

	return exps, nil
}

func MustStrExps(strs ...string) Matchers {
	exps, err := StrExps(strs...)
	if err != nil {
		panic(err)
	}

	return exps
}

type ErrBadExpr struct {
	inner error
}

func (e ErrBadExpr) Error() string {
	return fmt.Sprintf("%s.\n See https://tanka.dev/output-filtering/#regular-expressions for details on regular expressions.", strings.Title(e.inner.Error()))
}

type NegMatcher struct {
	exp Matcher
}

func (n NegMatcher) MatchString(s string) bool {
	return true
}

func (n NegMatcher) IgnoreString(s string) bool {
	return n.exp.MatchString(s)
}
