// Package scopes provides utilities for manipulating and interpreting Taskcluster scopes.
package scopes

import (
	"strings"
)

type (
	// `Given` represents a set of scopes assigned to a client.  For example:
	//
	//  myScopes := scopes.Given{
	//  	"abc:*",
	//  	"123:4:56",
	//  	"xyz",
	//  	"AB:*",
	//  }
	//
	// In order for a given scope to satisfy a required scope, either the given
	// scope and required scope need to match as strings, or the given scope
	// needs to be a prefix of the required scope, plus the `*` character. For
	// example, the given scope `abc:*` satisfies the required scope `abc:def`.
	Given []string
	// `Required` represents (in disjunctive normal form) permutations of
	// scopes that are sufficient to authorise a client to perform a particular
	// action. For example:
	//
	//  requiredScopes := scopes.Required{
	//  	{"abc:def", "AB:CD:EF"},
	//  	{"123:4:5"},
	//  	{"abc:def", "123:4"},
	//  	{"Xxyz"},
	//  }
	//
	// represents the requirement that the following scopes are "satisfied":
	//
	//  ("abc:def" AND "AB:CD:EF") OR "123:4:5" OR ("abc:def" AND "123:4") OR "Xxyz"
	//
	// Internally, Required types are comprised of a list of "scope sets". Only
	// one of its scope sets needs to be "satisfied" for the Required type to be
	// satisfied (the `OR`s above). Each scope set is said to be satisfied if
	// all of its scopes are satisfied (the `AND`s above). With this construction,
	// arbitrary scope requirements can be defined.
	//
	// Please note Required scopes do _not_ contain wildcard characters; they are
	// literal strings. This differs from Given scopes.
	Required []scopeSet
	// A scopeSet is a list of required scopes.
	scopeSet []string
)

// Returns `true` if the `given` scopes satisfy the `required` scopes.
func (given *Given) Satisfies(required *Required) bool {
	// if any of the required scopeSets are satisfied by `given`, then
	// `required` is said to be satisfied by `given`
	for _, set := range *required {
		if given.satisfiesScopeSet(&set) {
			return true
		}
	}
	return false
}

// Returns `true` if the `given` scopes satisfy all of the scopes in the
// scopeSet `set`
func (given *Given) satisfiesScopeSet(set *scopeSet) bool {
	// all scopes have to be satisfied in order for scope set to be satisfied
	for _, scope := range *set {
		if !given.satisfiesScope(&scope) {
			return false
		}
	}
	return true
}

// Returns `true` if any of the `given` scopes satisfies `requiredScope`.
func (given *Given) satisfiesScope(requiredScope *string) bool {
	// requiredScope is satisfied if any of the given scopes satisfies it
	for _, givenScope := range *given {
		if scopeMatch(&givenScope, requiredScope) {
			return true
		}
	}
	return false
}

// Returns `true` if `givenScope`  satisfies `requiredScope`.
func scopeMatch(givenScope *string, requiredScope *string) bool {
	return *requiredScope == *givenScope || (strings.HasSuffix(*givenScope, "*") && strings.HasPrefix(*requiredScope, (*givenScope)[0:len(*givenScope)-1]))
}
