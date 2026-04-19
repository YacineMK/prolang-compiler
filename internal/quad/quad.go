package quad

import "fmt"

type Quad struct {
	Op     string
	Arg1   string
	Arg2   string
	Result string
}

func (q Quad) String() string {
	return fmt.Sprintf("( %-6s , %-10s , %-10s , %-10s )",
		q.Op, q.Arg1, q.Arg2, q.Result)
}

type Manager struct {
	Quads      []Quad
	tempCount  int
	labelCount int
}

func NewManager() *Manager {
	return &Manager{
		Quads: make([]Quad, 0),
	}
}

func (m *Manager) Emit(op, arg1, arg2, result string) int {
	m.Quads = append(m.Quads, Quad{
		Op:     op,
		Arg1:   arg1,
		Arg2:   arg2,
		Result: result,
	})
	return len(m.Quads) - 1
}

func (m *Manager) NewTemp() string {
	name := fmt.Sprintf("T%d", m.tempCount)
	m.tempCount++
	return name
}

func (m *Manager) NewLabel() string {
	name := fmt.Sprintf("L%d", m.labelCount)
	m.labelCount++
	return name
}

func (m *Manager) EmitLabel(label string) {
	m.Emit("LABEL", label, "", "")
}

func (m *Manager) Backpatch(index int, label string) {
	if index >= 0 && index < len(m.Quads) {
		m.Quads[index].Result = label
	}
}

func Print(title string, quads []Quad) {
	fmt.Println("\n══════════════════════════════════════════")
	fmt.Println(title)
	fmt.Println("══════════════════════════════════════════")

	for i, q := range quads {
		fmt.Printf("%03d  %s\n", i, q.String())
	}

	fmt.Printf("\nTotal quads: %d\n", len(quads))
}