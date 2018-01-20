package alfred

import "testing"

func TestEvaluate(t *testing.T) {
	result := evaluate("whoami", ".")
	if result == "whoami" {
		t.Fatalf("evaluate(whoami) should return the command, not the text")
	}

	result = evaluate("commanddoesnotexist", ".")
	if result != "commanddoesnotexist" {
		t.Fatalf("testCommand(commanddoesnotexist) should return the text, as the command does not exist")
	}
}

func TestTestCommand(t *testing.T) {
	result := testCommand("whoami", ".")
	if result == false {
		t.Fatalf("testCommand(whoami) should evaluate to true")
	}
}

func TestExecute(t *testing.T) {
	result, ok := execute("whoami", ".")
	if result == "whoami" {
		t.Fatalf("Expected whoami to return a non string")
	}

	if !ok {
		t.Fatalf("Should not have resulted in an error.")
	}
}
