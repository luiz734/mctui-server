package subcommands

var Args CliArgs

type AddUserCmd struct {
	Username string `short:"u" name:"username" help:"Username for new user" required:""`
	Password string `short:"p" name:"password" help:"Password for new user" required:""`
}
type CliArgs struct {
	AddUser *AddUserCmd `cmd:"add-user" help:"Add a new user"`
	// Allows the program to run without args
	Dumb struct{} `cmd:"" default:""`
}
