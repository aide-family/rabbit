// Package enum is the enum package for the Rabbit service.
package enum

func (b Boolean) IsTrue() bool {
	return b == Boolean_true
}

func (b Boolean) IsFalse() bool {
	return b == Boolean_false
}

func (e Environment) IsUnknown() bool {
	return e == Environment_UNKNOWN
}

func (e Environment) IsDev() bool {
	return e == Environment_DEV
}

func (e Environment) IsTest() bool {
	return e == Environment_TEST
}

func (e Environment) IsPreview() bool {
	return e == Environment_PREVIEW
}

func (e Environment) IsProd() bool {
	return e == Environment_PROD
}
