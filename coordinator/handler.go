package coordinator

type handler struct{}

func (h *handler) Join(name, addr string) error {
	return nil
}
func (h *handler) Leave(name string) error {
	return nil
}
