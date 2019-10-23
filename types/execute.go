package types

// execute result
type ExecuteResult struct {
	Error  error  // error of execution
	StdOut string // standard output
	StdErr string // error output
}
