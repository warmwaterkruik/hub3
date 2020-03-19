package domain

import (
	"errors"
	"unicode"
)

// errors
var (
	ErrIDTooLong          = errors.New("identifier is too long")
	ErrIDNotLowercase     = errors.New("uppercase not allowed in identifier")
	ErrIDInvalidCharacter = errors.New("only letters are allowed in organization")
	ErrIDCannotBeEmpty    = errors.New("empty string is not a valid identifier")
	ErrIDExists           = errors.New("identifier already exists")
	ErrOrgNotFound        = errors.New("organization not found")
)

var (
	// MaxLengthID the maximum length of an identifier
	MaxLengthID = 10

	// protected organization names
	protected = []OrganizationID{
		OrganizationID("public"),
		OrganizationID("all"),
	}
)

// OrganizationID represents a short identifier for an Organization.
//
// The maximum length is MaxLengthID.
//
// In JSON the OrganizationID is represented as 'orgID'.
type OrganizationID string

// Organization is a basic building block for storing information.
// Everything that is stored by ikuzo must have an organization.ID as part of its metadata.
type Organization struct {
	ID          OrganizationID `json:"orgID"`
	Description string         `json:"description,omitempty"`
}

// NewOrganizationID returns an OrganizationID and an error if the supplied input is invalid.
func NewOrganizationID(input string) (OrganizationID, error) {
	id := OrganizationID(input)
	if err := id.Valid(); err != nil {
		return OrganizationID(""), err
	}

	return id, nil
}

// Valid validates the identifier.
//
// - ErrIDTooLong is returned when ID is too long
//
// - ErrIDNotLowercase is returned when ID contains uppercase characters
//
// - ErrIDInvalidCharacter is returned when ID contains non-letters
//
func (id OrganizationID) Valid() error {
	if id == "" {
		return ErrIDCannotBeEmpty
	}

	if len(id) > MaxLengthID {
		return ErrIDTooLong
	}

	for _, p := range protected {
		if id == p {
			return ErrIDExists
		}
	}

	for _, r := range id {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return ErrIDNotLowercase
		}

		if !unicode.IsLetter(r) {
			return ErrIDInvalidCharacter
		}
	}

	return nil
}