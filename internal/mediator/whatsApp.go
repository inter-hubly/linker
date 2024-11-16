package mediator

type WhatsApp interface {
	Persist()
}

type whatsApp struct {
}

func NewWhatsApp() WhatsApp {
	return &whatsApp{}

}
func (w *whatsApp) Persist() {

}
