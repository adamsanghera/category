package category

import (
	"errors"
	"time"

	"github.com/adamsanghera/badge"
	bus "github.com/adamsanghera/redisBus"
)

// Category is an interface for forming users into groups.
// Category is the nation, which users are citizens of.
// Membership in a category is maintained by a sorted list.
type Category interface {
	// Commands
	AddUser(uid interface{}) error
	RemoveUser(uid interface{}) error
	IssueBadge(uid interface{}) error
	RenewBadge(uid interface{}, badge badge.Badge) error
	RevokeBadge(uid interface{}) error

	// Queries
	IsGroupOf(uid interface{}) (bool, error)
	GetBadge(uid interface{}) (badge.Badge, time.Duration, error)
}

// UserGroup is an implementation of Category
type UserGroup struct {
	name string
}

// AddUser adds a user to Redis
func (g UserGroup) AddUser(uid interface{}) error {
	s, ok := uid.(string)
	if !ok {
		return errors.New("UserID provided was not a string")
	}
	_, err := bus.Client.Set(s, 1, 0).Result()
	return err
}

// RemoveUser deletes a user from Redis
func (g UserGroup) RemoveUser(uid interface{}) error {
	s, ok := uid.(string)
	if !ok {
		return errors.New("UserID provided was not a string")
	}
	_, err := bus.Client.Del(s).Result()
	return err
}

// IsGroupOf checks whether a user record exists in Redis
func (g UserGroup) IsGroupOf(uid interface{}) (bool, error) {
	s, ok := uid.(string)
	if !ok {
		return false, errors.New("UserID provided was not a string")
	}
	v, err := bus.Client.Get(s).Result()
	return v == "1", err
}
