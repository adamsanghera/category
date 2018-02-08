package category

import (
	"errors"
	"reflect"
	"time"

	"github.com/adamsanghera/badge"

	"github.com/go-redis/redis"
)

// Category is an interface for forming users into groups.
// Category is the nation, which users are citizens of.
// Membership in a category is maintained by a sorted list.
type Category interface {
	// Commands
	AddUser(uid interface{}) error
	RemoveUser(uid interface{}) error
	IssueBadge(uid interface{}) error
	RenewBadge(uid interface{}, bid interface{}) error
	RevokeBadge(uid interface{}) error

	// Queries
	IsGroupOf(uid interface{}) (bool, error)
	GetBadge(uid interface{}) (interface{}, reflect.Kind, time.Duration, error)
}

// UserGroup is an implementation of Category
type UserGroup struct {
	name         string
	minter       badge.RandomTokenMinter
	timeToLive   time.Duration
	badgeDB      redis.Client
	membershipDB redis.Client
}

// AddUser adds a user to Redis
func (g UserGroup) AddUser(uid interface{}) error {
	s, ok := uid.(string)
	if !ok {
		return errors.New("UserID provided was not a string")
	}
	_, err := g.membershipDB.Set(s, 1, 0).Result()
	return err
}

// RemoveUser deletes a user from Redis
func (g UserGroup) RemoveUser(uid interface{}) error {
	s, ok := uid.(string)
	if !ok {
		return errors.New("UserID provided was not a string")
	}
	_, err := g.membershipDB.Del(s).Result()
	return err
}

// IssueBadge creates a new badge, and associates it with the given user.
func (g UserGroup) IssueBadge(uid interface{}) error {
	s, ok := uid.(string)
	if !ok {
		return errors.New("UserID provided was not a string")
	}

	// Create the token
	token, _, err := g.minter.Mint()
	if err != nil {
		return err
	}

	// Set the token
	_, err = g.badgeDB.Set(s, token, g.timeToLive).Result()
	return err
}

// RenewBadge renews a given badge for the given user, after checking that the association is valid.
func (g UserGroup) RenewBadge(uid interface{}, bid interface{}) error {
	s, ok := uid.(string)
	if !ok {
		return errors.New("UserID provided was not a string")
	}
	b, ok := bid.(string)
	if !ok {
		return errors.New("Badge provided was not a string")
	}
	badge, err := g.badgeDB.Get(s).Result()
	if err != nil {
		return errors.New("Badge expired")
	}
	if badge != b {
		return errors.New("Incorrect badge")
	}
	return g.IssueBadge(uid)
}

// RevokeBadge removes a token for the given user.
func (g UserGroup) RevokeBadge(uid interface{}) error {
	s, ok := uid.(string)
	if !ok {
		return errors.New("UserID provided was not a string")
	}

	// Delete the badge
	_, err := g.badgeDB.Del(s).Result()
	return err
}

// IsGroupOf checks whether a user record exists in Redis
func (g UserGroup) IsGroupOf(uid interface{}) (bool, error) {
	s, ok := uid.(string)
	if !ok {
		return false, errors.New("UserID provided was not a string")
	}
	v, err := g.membershipDB.Get(s).Result()
	return v == "1", err
}

func (g UserGroup) GetBadge(uid interface{}) (interface{}, reflect.Kind, time.Duration, error) {
	s, ok := uid.(string)
	if !ok {
		return nil, reflect.TypeOf(nil).Kind(), 0, errors.New("UserID provided was not a string")
	}
	badge, err := g.badgeDB.Get(s).Result()
	if err != nil {
		return nil, reflect.TypeOf(nil).Kind(), 0, errors.New("Badge expired")
	}
	return badge, reflect.TypeOf("").Kind(), g.badgeDB.TTL(s).Val(), nil
}
