// Package meta contains a metadata store used across multiple services. It is used to pass the headers from the request
// to other services.
package meta

import (
	"context"
	"strings"

	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type metadataKey struct{}

// Meta is the metadata container used in the flow of applications.
type Meta map[string]string

// Constant metadata key definitions.
const (
	KeyContentLanguage = "contentLanguage"
	KeyAcceptLanguages = "acceptLanguage"
	KeyCurrency        = "currency"
	KeyUserID          = "userId"
	KeyToken           = "token"
	KeyContentEncoding = "contentEncoding"
	KeyRemoteIP        = "remoteIp"
)

// Get gets the key from the given metadata.
func (m Meta) Get(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

// Set sets the key in the given metadata.
func (m Meta) Set(key, value string) {
	m[key] = value
}

// Delete clears the key in the given metadata.
func (m Meta) Delete(key string) {
	delete(m, key)
}

// SetContentEncoding sets the content encoding.
func (m Meta) SetContentEncoding(encoding string) {
	m[KeyContentEncoding] = encoding
}

// ContentEncoding sets the content encoding.
func (m Meta) ContentEncoding() (string, bool) {
	v, ok := m[KeyContentEncoding]
	return v, ok
}

// SetContentLanguage sets up language metadata.
func (m Meta) SetContentLanguage(tag language.Tag) {
	m[KeyContentLanguage] = tag.String()
}

// ContentLanguage gets the content language from the metadata.
func (m Meta) ContentLanguage() (language.Tag, bool) {
	t, ok := m[KeyContentLanguage]
	if !ok {
		return language.Tag{}, false
	}
	tag, err := language.Parse(t)
	if err != nil {
		return language.Tag{}, false
	}
	return tag, true
}

// SetAcceptLanguages sets up request accepted languages by the user with the quality order.
func (m Meta) SetAcceptLanguages(tags []language.Tag) {
	var sb strings.Builder
	for i, tag := range tags {
		sb.WriteString(tag.String())
		if i < len(tags)-1 {
			sb.WriteRune(',')
		}
	}
	m[KeyAcceptLanguages] = sb.String()
}

// AcceptLanguages gets accepted languages by the user stored by the quality.
func (m Meta) AcceptLanguages() ([]language.Tag, bool) {
	langs, ok := m[KeyAcceptLanguages]
	if !ok {
		return nil, false
	}
	var (
		tags []language.Tag
		sb   strings.Builder
	)
	for i, rn := range langs {
		if rn != ',' {
			sb.WriteRune(rn)
			if i != len(langs)-1 {
				continue
			}
		}

		tag, err := language.Parse(sb.String())
		if err != nil {
			return nil, false
		}
		tags = append(tags, tag)
		sb.Reset()
	}
	return tags, true
}

// AcceptLanguage gets accepted language with the highest quality.
func (m Meta) AcceptLanguage() (language.Tag, bool) {
	tags, ok := m.AcceptLanguages()
	if !ok {
		return language.Tag{}, false
	}
	if len(tags) == 0 {
		return language.Tag{}, false
	}
	return tags[0], true
}

// SetCurrency sets up currency metadata.
func (m Meta) SetCurrency(c currency.Unit) {
	m[KeyCurrency] = c.String()
}

// Currency gets the currency from the map.
func (m Meta) Currency() (currency.Unit, bool) {
	c, ok := m[KeyCurrency]
	if !ok {
		return currency.Unit{}, false
	}
	u, err := currency.ParseISO(c)
	if err != nil {
		return currency.Unit{}, false
	}
	return u, true
}

// SetUserID sets up user id
func (m Meta) SetUserID(userID string) {
	m[KeyUserID] = userID
}

// UserID gets the user id from the metadata.
func (m Meta) UserID() (string, bool) {
	u, ok := m[KeyUserID]
	return u, ok
}

// SetRemoteIP sets the remoteIP in the metadata.
func (m Meta) SetRemoteIP(remoteIP string) {
	m[KeyRemoteIP] = remoteIP
}

// RemoteIP gets the remote IP from the metadata.
func (m Meta) RemoteIP() (string, bool) {
	v, ok := m[KeyRemoteIP]
	return v, ok
}

// SetToken sets the token in the metadata context.
func (m Meta) SetToken(token string) {
	m[KeyToken] = token
}

// Token gets the token from the metadata.
func (m Meta) Token() (string, bool) {
	v, ok := m[KeyToken]
	return v, ok
}

// Get gets the data from the Meta map.
func Get(ctx context.Context, key string) (string, bool) {
	md, ok := FromContext(ctx)
	if !ok {
		return "", false
	}
	val, ok := md[key]
	return val, ok
}

// MustFromContext gets the metadata from the context.
func MustFromContext(ctx context.Context) Meta {
	m, ok := ctx.Value(metadataKey{}).(Meta)
	if !ok {
		return Meta{}
	}
	return m
}

// FromContext gets the metadata from the context.
func FromContext(ctx context.Context) (Meta, bool) {
	m, ok := ctx.Value(metadataKey{}).(Meta)
	if !ok {
		return nil, false
	}
	return m, true
}

// NewContext creates a new context with the given metadata
func NewContext(ctx context.Context, md Meta) context.Context {
	return context.WithValue(ctx, metadataKey{}, md)
}

// MergeContext merges metadata to existing metadata, overwriting if specified
func MergeContext(ctx context.Context, patchMd Meta, overwrite bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	md, _ := ctx.Value(metadataKey{}).(Meta)
	cmd := make(Meta, len(md))
	for k, v := range md {
		cmd[k] = v
	}
	for k, v := range patchMd {
		if _, ok := cmd[k]; ok && !overwrite {
			// skip
		} else if v != "" {
			cmd[k] = v
		} else {
			delete(cmd, k)
		}
	}
	return context.WithValue(ctx, metadataKey{}, cmd)
}

// ContentLanguage gets the content language from the context metadata.
func ContentLanguage(ctx context.Context) (language.Tag, bool) {
	md, ok := FromContext(ctx)
	if !ok {
		return language.Tag{}, false
	}
	return md.ContentLanguage()
}

// AcceptLanguages gets the accepted languages stored in the
func AcceptLanguages(ctx context.Context) ([]language.Tag, bool) {
	md, ok := FromContext(ctx)
	if !ok {
		return nil, false
	}
	return md.AcceptLanguages()
}

// UserID gets metadata userID.
func UserID(ctx context.Context) (string, bool) {
	md, ok := FromContext(ctx)
	if !ok {
		return "", false
	}
	return md.UserID()
}

// Token gets metadata token.
func Token(ctx context.Context) (string, bool) {
	md, ok := FromContext(ctx)
	if !ok {
		return "", false
	}
	return md.Token()
}

// RemoteIP gets metadata remote ip.
func RemoteIP(ctx context.Context) (string, bool) {
	md, ok := FromContext(ctx)
	if !ok {
		return "", false
	}
	return md.RemoteIP()
}
