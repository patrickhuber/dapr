package plugin

type Launcher interface {
	Launch() (Client, error)
}
