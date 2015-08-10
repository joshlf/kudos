// Package purple_unicorn provides types which represent the
// configurations of various entities in kudos such as courses
// and assignments. These types are able to verify that a given
// configuration satisfies all of the constraints associated
// with that type of configuration (for example, that an
// assignment has at least one problem).
//
// Each type which represents a configuration implements the
// Validator interface, which means that it is able to validate
// its configuration using either the Validate or MustValidate
// methods. Additionally, while it is not enforced by the type
// system, each property which could, when mutated, violate the
// constraints of the configuration has three setters with the
// following behavior:
//
//  (c *Config) SetX(x X) error
// SetX sets property X, and validates to make sure that the
// change has not violated any of c's constraints. If it has,
// the change is rolled back so that c is unchanged, and an
// error is returned. Otherwise, the change stays, and the return
// value is nil.
//
//  (c *Config) MustSetX(x X)
// MustSetX is identical to SetX, except that if an error is
// encountered, it panics instead of returning the error.
//
//  (c *Config) SetXNoValidate(x X)
// SetXNoValidate sets property x, but does not perform any
// subsequent validation. If the mutation puts c in an invalid
// state, it is left there. It is up to the caller to later
// call Validate or MustValidate in order to assure the validity
// of c.
//
// If the use of a NoValidate setter is desired (usually for
// performance reasons or because the operation will intentionally
// place the configuration into an invalid state, but a future
// operation will bring it back into a valid state), and the
// caller wishes to be able to recover the previous good state
// in case validation fails, the caller should use the configuration's
// Copy method to deep copy the configuration so that it can be
// restored from the backup in case validation fails.
package purple_unicorn
