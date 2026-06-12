package health

// Check represents a single health validation.
type Check interface {
	Name() string
	Run() CheckResult
}
