package xauth

import (
	"sort"
	"strings"
)

// Scopes is a slice of auth scopes that allows basic operation on checking and comparing multiple scopes.
type Scopes []Scope

// Sort sorts the scopes at first by the common domain, than by each sub scope.
func (s Scopes) Sort() {
	sort.Slice(s, func(i, j int) bool {
		s1, s2 := s[i], s[j]
		if s1.Domain != s2.Domain {
			return s1.Domain < s2.Domain
		}
		minLength := minScopeLength(s1, s2)
		for k := 0; k < minLength; k++ {
			if s1.SubScopes[k] != s2.SubScopes[k] {
				return s1.SubScopes[k] < s2.SubScopes[k]
			}
		}
		return len(s1.SubScopes) < len(s2.SubScopes)
	})
}

// MatchScopes checks if the 'request' scopes matches 'origin' scopes in a hierarchical way.
func MatchScopes(originScopes, request Scopes) bool {
originLoop:
	for i := range originScopes {
		for j := range request {
			if originScopes[i].Domain != request[j].Domain {
				continue
			}
			if len(originScopes[i].SubScopes) < len(request[j].SubScopes) {
				// No matter what are the sub scope values stored in the 'toCheck' scope,
				// when the origin is less specific than the 'toCheck' - 'toCheck' wouldn't match.
				continue
			}

			// Check if all the 'toCheck' sub scopes matches with the origin.
			for k := range request[j].SubScopes {
				if request[j].SubScopes[k] != originScopes[i].SubScopes[k] {
					// If any of the subsequent sub scopes doesn't match to the origin it should return false.
					continue
				}
			}
			// Passes, continue with another scope from origin scopes.
			continue originLoop
		}
		return false
	}
	return true
}

// Scope is the authorization scope unmarshaled into computable form.
type Scope struct {
	Domain    string
	SubScopes []string
}

// String formats the scope into a string format.
func (s *Scope) String() string {
	sb := strings.Builder{}
	if len(s.Domain) != 0 {
		sb.WriteString(s.Domain)
		sb.WriteRune('/')
	}
	for i, subScope := range s.SubScopes {
		sb.WriteString(subScope)
		if i != len(s.SubScopes)-1 {
			sb.WriteRune('.')
		}
	}
	return sb.String()
}

// Match determines if a scope 'toCheck' matches the 'origin' scope.
// The origin might be more specific than the 'toCheck' scope.
func Match(origin, toCheck Scope) bool {
	if origin.Domain != toCheck.Domain {
		return false
	}
	if len(origin.SubScopes) < len(toCheck.SubScopes) {
		// No matter what are the sub scope values stored in the 'toCheck' scope,
		// when the origin is less specific than the 'toCheck' - 'toCheck' wouldn't match.
		return false
	}

	// Check if all the 'toCheck' sub scopes matches with the origin.
	for i := range toCheck.SubScopes {
		if toCheck.SubScopes[i] != origin.SubScopes[i] {
			// If any of the subsequent sub scopes doesn't match to the origin it should return false.
			return false
		}
	}
	return true
}

func minScopeLength(s1, s2 Scope) int {
	if len(s1.SubScopes) < len(s2.SubScopes) {
		return len(s1.SubScopes)
	}
	return len(s2.SubScopes)
}

// SplitScopes splits the scope stored in the oauth query or token.
func SplitScopes(scope string) []string {
	return strings.Split(scope, " ")
}

// ParseScopes parses the scope into hierarchical structure.
func ParseScopes(inputScopes []string) Scopes {
	if len(inputScopes) == 0 {
		return nil
	}
	// Split the scope by the whitespace.
	scopes := make([]Scope, len(inputScopes))
	for i := range inputScopes {
		scopes[i] = parseScope(inputScopes[i])
	}
	return scopes
}

// ParseScope parses the string scope into a authorization scope.
func ParseScope(scope string) Scope {
	return parseScope(scope)
}

func parseScope(scope string) Scope {
	var s Scope
	if domainIndex := strings.LastIndexByte(scope, '/'); domainIndex != -1 {
		s.Domain, scope = scope[:domainIndex], scope[domainIndex+1:]
	}
	s.SubScopes = strings.Split(scope, ".")
	return s
}
