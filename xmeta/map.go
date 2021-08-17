// Package xmeta contains helper functions and keys to use with grpc metadata.
package xmeta

import (
	"net"

	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"google.golang.org/grpc/metadata"
)

// Constant metadata key definitions.
const (
	KeyAuthorization   = "__cg_authorization__"
	KeyContentLanguage = "__cg_content_language__"
	KeyAcceptLanguages = "__cg_accept_language__"
	KeyCurrency        = "__cg_currency__"
	KeyUserID          = "__cg_user_id__"
	KeyRemoteIP        = "__cg_remote_ip__"
	KeyRequestID       = "__cg_request_id__"
)

// SetUserID sets up user id in the metadata.
func SetUserID(md metadata.MD, userID string) {
	md.Set(KeyUserID, userID)
}

// UserID gets the user id from the metadata.
func UserID(md metadata.MD) (string, bool) {
	u := md.Get(KeyUserID)
	if len(u) == 0 {
		return "", false
	}
	return u[0], true
}

// SetRemoteIP sets the remoteIP in the metadata.
func SetRemoteIP(md metadata.MD, remoteIP net.IP) {
	md.Set(KeyRemoteIP, remoteIP.String())
}

// SetRequestID sets the request id in the metadata.
func SetRequestID(md metadata.MD, requestID string) {
	md.Set(KeyRequestID, requestID)
}

// SetAuthorization sets the token in the metadata context.
func SetAuthorization(md metadata.MD, auth string) {
	md.Set(KeyAuthorization, auth)
}

// AcceptLanguages gets the accepted languages stored in the context metadata.
func AcceptLanguages(md metadata.MD) ([]language.Tag, bool) {
	langs, ok := md[KeyAcceptLanguages]
	if !ok {
		return nil, false
	}
	var tags []language.Tag
	for i, lang := range langs {
		if i != len(langs)-1 {
			continue
		}
		tag, err := language.Parse(lang)
		if err != nil {
			return nil, false
		}
		tags = append(tags, tag)
	}
	return tags, true
}

// SetAcceptLanguages sets up accept languages in the metadata.
func SetAcceptLanguages(md metadata.MD, tags []language.Tag) {
	values := make([]string, len(tags))
	for i, tag := range tags {
		values[i] = tag.String()
	}
	md.Set(KeyAcceptLanguages, values...)
}

// ContentLanguage gets the content language from the context metadata.
func ContentLanguage(md metadata.MD) (language.Tag, bool) {
	if md == nil {
		return language.Tag{}, false
	}
	t, ok := md[KeyContentLanguage]
	if !ok {
		return language.Tag{}, false
	}
	if len(t) == 0 {
		return language.Tag{}, false
	}
	tag, err := language.Parse(t[0])
	if err != nil {
		return language.Tag{}, false
	}
	return tag, true
}

// SetContentLanguage sets the content language in the context.
func SetContentLanguage(md metadata.MD, tag language.Tag) {
	md.Set(KeyContentLanguage, tag.String())
}

// Currency gets the metadata currency.
func Currency(md metadata.MD) (currency.Unit, bool) {
	c, ok := md[KeyCurrency]
	if !ok {
		return currency.Unit{}, false
	}
	if len(c) == 0 {
		return currency.Unit{}, false
	}
	u, err := currency.ParseISO(c[0])
	if err != nil {
		return currency.Unit{}, false
	}
	return u, true
}

// SetCurrency sets up given currency in the metadata.
func SetCurrency(md metadata.MD, c currency.Unit) {
	md.Set(KeyCurrency, c.String())
}

// AcceptLanguage gets the first accepted language stored in the context metadata.
func AcceptLanguage(md metadata.MD) (language.Tag, bool) {
	tags, ok := md[KeyAcceptLanguages]
	if !ok {
		return language.Tag{}, false
	}
	if len(tags) == 0 {
		return language.Tag{}, false
	}
	tag, err := language.Parse(tags[0])
	if err != nil {
		return language.Tag{}, false
	}
	return tag, true
}

// Authorization gets metadata token.
func Authorization(md metadata.MD) (string, bool) {
	return getKey(md, KeyAuthorization)
}

// RemoteIP gets metadata remote ip.
func RemoteIP(md metadata.MD) (string, bool) {
	return getKey(md, KeyRemoteIP)
}

// RequestID gets the request identifier.
func RequestID(md metadata.MD) (string, bool) {
	return getKey(md, KeyRequestID)
}

func getKey(md metadata.MD, key string) (string, bool) {
	values := md.Get(key)
	if len(values) == 0 {
		return "", false
	}
	return values[0], true
}
