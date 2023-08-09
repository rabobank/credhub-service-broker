package security

import (
	"github.com/gorilla/mux"
)

// MiddlewareBuilder interface for security middleware builder functions
type MiddlewareBuilder interface {
	MatchPrefix(pathPrefix string) MatchBuilder
	Default(authenticator Authenticator) MiddlewareBuilder
	Build() mux.MiddlewareFunc
}

// MatchBuilder interface to build matching rules to apply a specific authentication method
type MatchBuilder interface {
	MatchPrefix(pathPrefix string) MatchBuilder
	AuthenticateWith(authenticator Authenticator) MiddlewareBuilder
}

type middlewareBuilder struct {
	matchingAuthenticators []authenticatorMatchers
	defaultAuthenticator   Authenticator
}

func (b *middlewareBuilder) addMatchingAuthenticatorBuilder(matchingAuthenticatorBuilder *matchingAuthenticatorBuilder) {
	b.matchingAuthenticators = append(b.matchingAuthenticators, authenticatorMatchers{
		matchers:     matchingAuthenticatorBuilder.matchers,
		authenticate: matchingAuthenticatorBuilder.authenticator,
	})
}

func (b *middlewareBuilder) Default(authenticator Authenticator) MiddlewareBuilder {
	if b.defaultAuthenticator != nil {
		panic("default authentication has already been set")
	}
	b.defaultAuthenticator = authenticator
	return b
}

func (b *middlewareBuilder) Build() mux.MiddlewareFunc {
	authenticate := b.defaultAuthenticator
	if authenticate == nil {
		authenticate = Anonymous
	}
	authenticators := b.matchingAuthenticators
	authenticators = append(authenticators, authenticatorMatchers{[]string{"/"}, authenticate})
	return Middleware(authenticators)
}

func (b *middlewareBuilder) MatchPrefix(prefix string) MatchBuilder {
	builder := &matchingAuthenticatorBuilder{middlewareBuilder: b}
	return builder.MatchPrefix(prefix)
}

type matchingAuthenticatorBuilder struct {
	middlewareBuilder *middlewareBuilder
	matchers          []string
	authenticator     Authenticator
}

func (b *matchingAuthenticatorBuilder) MatchPrefix(prefix string) MatchBuilder {
	b.matchers = append(b.matchers, prefix)
	return b
}

func (b *matchingAuthenticatorBuilder) AuthenticateWith(authenticator Authenticator) MiddlewareBuilder {
	b.authenticator = authenticator
	b.middlewareBuilder.addMatchingAuthenticatorBuilder(b)
	return b.middlewareBuilder
}
