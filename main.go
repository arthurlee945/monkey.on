package main

//30
import (
	"fmt"
	"os"
	"os/user"

	"github.com/arthurlee945/monkey.on/repl"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}
	fmt.Printf("Monkey.on is read to Monkeying %s!\n", user.Username)
	fmt.Printf("Start type in commands:\n\n")

	// repl.Start(os.Stdin, os.Stdout)
	repl.StartParser(os.Stdin, os.Stdout)
}
