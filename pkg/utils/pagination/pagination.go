package pagination

import (
	"encoding/base64"
	"errors"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	DefaultPageSize int32 = 100
	MaximumPageSize int32 = 5000
	CursorDuration        = 24 * time.Hour
)

var (
	ErrMissingExpiration  = errors.New("cursor does not have an expires timestamp")
	ErrCursorExpired      = errors.New("cursor has expired and is no longer useable")
	ErrUnparsableToken    = errors.New("could not parse the next page token")
	ErrTokenQueryMismatch = errors.New("cannot change query parameters during pagination")
)

func New(startIndex, endIndex string, pageSize int32) *Cursor {
	if pageSize == 0 {
		pageSize = DefaultPageSize
	}

	return &Cursor{
		StartIndex: startIndex,
		EndIndex:   endIndex,
		PageSize:   pageSize,
		Expires:    timestamppb.New(time.Now().Add(CursorDuration)),
	}
}

func Parse(token string) (cursor *Cursor, err error) {
	var data []byte
	if data, err = base64.RawURLEncoding.DecodeString(token); err != nil || len(data) == 0 {
		return nil, ErrUnparsableToken
	}

	cursor = &Cursor{}
	if err = proto.Unmarshal(data, cursor); err != nil {
		return nil, ErrUnparsableToken
	}

	// TODO: ensure the request matches the original query and return ErrTokenQueryMismatch

	// Ensure the cursor has not expired
	var expired bool
	if expired, err = cursor.HasExpired(); err != nil {
		return nil, err
	}

	if expired {
		return nil, ErrCursorExpired
	}

	return cursor, nil
}

func (c *Cursor) NextPageToken() (token string, err error) {
	var expired bool
	if expired, err = c.HasExpired(); err != nil {
		return "", err
	}

	if expired {
		return "", ErrCursorExpired
	}

	var data []byte
	if data, err = proto.Marshal(c); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func (c *Cursor) HasExpired() (bool, error) {
	if c.Expires == nil {
		return false, ErrMissingExpiration
	}
	return time.Now().After(c.Expires.AsTime()), nil
}

func (c *Cursor) IsZero() bool {
	return c.StartIndex == "" && c.EndIndex == "" && c.PageSize == 0 && c.Expires == nil
}
