package hypr


type Workspace struct {
	ID int;
	Name string;
}
type Client struct{
	Workspace Workspace;
	Monitor int;
	Title string;
	Class string;
	PID int;
}