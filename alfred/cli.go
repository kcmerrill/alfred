package alfred

// CLI will take in os.args and return the task and argument
func CLI(args []string) (string, []string) {
	// only the applicaton was passed in
	if len(args) <= 1 {
		return "", []string{}
	}
	return args[1], args[2:]
}
