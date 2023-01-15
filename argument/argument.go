package argument

import "os"

const (
	PLAIN     = "plain"     // Switch off the authentication and encryption for SDS Service
	BROADCAST = "broadcast" // runs only broadcaster
	REPLY     = "reply"     // runs only request-reply server
)

// any command line data that comes after the files are .env file paths
// Any argument for application without '--' prefix is considered to be path to the
// environment file.
func GetEnvPaths() ([]string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, nil
	}

	paths := make([]string, 0)

	for _, arg := range args {
		if arg[:2] != "--" {
			paths = append(paths, arg)
		}
	}

	return paths, nil
}

// Load arguments, not the environment variable paths.
// Arguments are with --prefix
func GetArguments() ([]string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, nil
	}

	parameters := make([]string, 0)

	for _, arg := range args {
		if arg[:2] == "--" {
			parameters = append(parameters, arg[2:])
		}
	}

	return parameters, nil
}

// This function is same as `env.HasArgument`,
// except `env.ArgumentExist()` loads arguments automatically.
func Exist(argument string) (bool, error) {
	arguments, err := GetArguments()
	if err != nil {
		return false, err
	}

	return Has(arguments, argument), nil
}

// Whehter the given argument exists or not.
func Has(arguments []string, required string) bool {
	for _, argument := range arguments {
		if argument == required {
			return true
		}
	}

	return false
}
