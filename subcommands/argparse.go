package subcommands

var Args CliArgs

type AddUserCmd struct {
	Username string `short:"u" name:"username" help:"Username for new user" required:""`
	Password string `short:"p" name:"password" help:"Password for new user" required:""`
}

type ListCmd struct{}

type CliArgs struct {
	AddUser *AddUserCmd `cmd:"add-user" help:"Add a new user"`
	List    *ListCmd    `cmd:"list" help:"List users"`

	// Allows the program to run without args
	Dumb struct{} `cmd:"" default:""`
}
