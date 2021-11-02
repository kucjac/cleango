// Package xmeta contains helper functions and keys to use with grpc metadata.
package xmeta

import (
	"context"
	"net"

	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"google.golang.org/grpc/metadata"
)

// Constant metadata key definitions.
const (
	KeyAuthorization   = "authorization"
	KeyContentLanguage = "content_language"
	KeyAcceptLanguages = "accept_language"
	KeyCurrency        = "currency"
	KeyUserID          = "user_id"
	KeyRemoteIP        = "remote_ip"
	KeyRequestID       = "request_id"
)

// IncomingCtxUserID gets the UserID from the incoming context.
func IncomingCtxUserID(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	return getKey(md, KeyUserID)
}

// IncomingCtxSetUserID sets up UserID in the incoming context metadata.
func IncomingCtxSetUserID(ctx context.Context, userID string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	SetUserID(md, userID)
	return metadata.NewIncomingContext(ctx, md)
}

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

// IncomingCtxRemoteIP gets the RemoteIP from the incoming context.
func IncomingCtxRemoteIP(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	return getKey(md, KeyRemoteIP)
}

// IncomingCtxSetRemoteIP sets up RemoteIP in the incoming context metadata.
func IncomingCtxSetRemoteIP(ctx context.Context, remoteIP net.IP) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	SetRemoteIP(md, remoteIP)
	return metadata.NewIncomingContext(ctx, md)
}

// SetRemoteIP sets the remoteIP in the metadata.
func SetRemoteIP(md metadata.MD, remoteIP net.IP) {
	md.Set(KeyRemoteIP, remoteIP.String())
}

// IncomingCtxRequestID gets the RequestID from the incoming context.
func IncomingCtxRequestID(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	return getKey(md, KeyRequestID)
}

// IncomingCtxSetRequestID sets up RequestID in the incoming context metadata.
func IncomingCtxSetRequestID(ctx context.Context, requestID string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	SetRequestID(md, requestID)
	return metadata.NewIncomingContext(ctx, md)
}

// SetRequestID sets the request id in the metadata.
func SetRequestID(md metadata.MD, requestID string) {
	md.Set(KeyRequestID, requestID)
}

// IncomingCtxAuthorization gets the Authorization from the incoming context.
func IncomingCtxAuthorization(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	return getKey(md, KeyAuthorization)
}

// IncomingCtxSetAuthorization sets up Authorization in the incoming context metadata.
func IncomingCtxSetAuthorization(ctx context.Context, auth string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	SetAuthorization(md, auth)
	return metadata.NewIncomingContext(ctx, md)

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

// IncomingCtxAcceptLanguages gets the AcceptLanguages from the incoming context.
func IncomingCtxAcceptLanguages(ctx context.Context) ([]language.Tag, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, false
	}
	return AcceptLanguages(md)
}

// IncomingCtxSetAcceptLanguages sets up AcceptLanguages in the incoming context metadata.
func IncomingCtxSetAcceptLanguages(ctx context.Context, tags []language.Tag) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	SetAcceptLanguages(md, tags)
	return metadata.NewIncomingContext(ctx, md)
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

// IncomingCtxContentLanguage gets the ContentLanguage from the incoming context.
func IncomingCtxContentLanguage(ctx context.Context) (language.Tag, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return language.Tag{}, false
	}
	return ContentLanguage(md)
}

// IncomingCtxSetContentLanguage sets up ContentLanguage in the incoming context metadata.
func IncomingCtxSetContentLanguage(ctx context.Context, tag language.Tag) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	SetContentLanguage(md, tag)
	return metadata.NewIncomingContext(ctx, md)
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

// IncomingCtxCurrency gets the Currency from the incoming context.
func IncomingCtxCurrency(ctx context.Context) (currency.Unit, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return currency.Unit{}, false
	}
	return Currency(md)
}

// IncomingCtxSetCurrency sets up Currency in the incoming context metadata.
func IncomingCtxSetCurrency(ctx context.Context, c currency.Unit) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	SetCurrency(md, c)
	return metadata.NewIncomingContext(ctx, md)
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
