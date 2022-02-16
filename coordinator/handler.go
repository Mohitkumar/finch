package coordinator

type MemberChangeType int

const (
	MemberJoin  MemberChangeType = 1
	MemberLeave MemberChangeType = 2
)

type MemberChangeMsg struct {
	Name             string
	Address          string
	MemberChangeType MemberChangeType
}
type handler struct {
	ch chan MemberChangeMsg
}

func (h *handler) Join(name, addr string) error {
	msg := MemberChangeMsg{
		Name:             name,
		Address:          addr,
		MemberChangeType: MemberJoin,
	}
	h.ch <- msg
	return nil
}
func (h *handler) Leave(name string) error {
	msg := MemberChangeMsg{
		Name:             name,
		MemberChangeType: MemberLeave,
	}
	h.ch <- msg
	return nil
}
