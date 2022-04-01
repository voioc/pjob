package define

type ServerList struct {
	GroupID   int
	GroupName string
	Servers   map[int]string
}
