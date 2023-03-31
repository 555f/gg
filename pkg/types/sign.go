package types

type Sign struct {
	Params     Vars
	Results    Vars
	IsVariadic bool
	IsNamed    bool
	Recv       any
}
