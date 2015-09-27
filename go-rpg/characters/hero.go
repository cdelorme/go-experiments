package characters

type Hero struct {
	Health  uint32
	Mana    uint32
	Attack  int32
	Defense int32
	Name    string
}

// add maxHealth setting

func (h *Hero) Fight(m *Monster) {

	// roll for initiative?
	// ignore negatives to avoid adding health during attacks
	h.Health = h.Health - uint32(h.Defense-m.Attack) // bad logic~

	// counter
	m.Health = m.Health - uint32(m.Defense-h.Attack) // also bad

	h.Update()
}

func (h *Hero) Update() {
	h.regen()
}

func (h *Hero) regen() {
	// prevent regen beyond maxHealth
	h.Health += 2
}
