package ast

import (
	"fmt"
	"strings"
)

// ════════════════════════════════════════════════════════════════
//  SYMBOL TABLE — Table de hachage pour stockage des symboles
// ════════════════════════════════════════════════════════════════

// ─── Énumérations ───────────────────────────────────────────────

type SymbolKind string

const (
	KIND_VAR   SymbolKind = "variable"
	KIND_CONST SymbolKind = "constante"
	KIND_ARRAY SymbolKind = "tableau"
)

type DataType string

const (
	TYPE_INTEGER DataType = "integer"
	TYPE_FLOAT   DataType = "float"
	TYPE_UNKNOWN DataType = "unknown"
)

// ─── Symbol ─────────────────────────────────────────────────────

type Symbol struct {
	Name        string
	Kind        SymbolKind
	Type        DataType
	Value       string
	ArraySize   int
	Line        int
	Column      int
	Initialized bool
}

// ─── SymbolTable (Hash Table with Chaining) ─────────────────────

const TABLE_SIZE = 256

type entry struct {
	key  string
	sym  Symbol
	next *entry
}

type SymbolTable struct {
	buckets [TABLE_SIZE]*entry
	count   int
}

// NewSymbolTable crée une nouvelle table de symboles
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{}
}

// hashStr calcule le hash d'une chaîne
func hashStr(key string) int {
	h := 0
	for _, c := range key {
		h = (h*31 + int(c)) % TABLE_SIZE
	}
	return h
}

// Insert ajoute un symbole à la table
// Retourne une erreur si l'identificateur est déjà déclaré
func (st *SymbolTable) Insert(sym Symbol) error {
	if st.Lookup(sym.Name) != nil {
		return fmt.Errorf("identificateur '%s' déjà déclaré", sym.Name)
	}
	idx := hashStr(sym.Name)
	e := &entry{key: sym.Name, sym: sym, next: st.buckets[idx]}
	st.buckets[idx] = e
	st.count++
	return nil
}

// Lookup recherche un symbole par nom
// Retourne nil si non trouvé
func (st *SymbolTable) Lookup(name string) *Symbol {
	idx := hashStr(name)
	for e := st.buckets[idx]; e != nil; e = e.next {
		if e.key == name {
			return &e.sym
		}
	}
	return nil
}

// MarkInitialized marque un symbole comme initialisé
func (st *SymbolTable) MarkInitialized(name string) {
	idx := hashStr(name)
	for e := st.buckets[idx]; e != nil; e = e.next {
		if e.key == name {
			e.sym.Initialized = true
			return
		}
	}
}

// Update met à jour la valeur d'un symbole
func (st *SymbolTable) Update(name, value string) error {
	idx := hashStr(name)
	for e := st.buckets[idx]; e != nil; e = e.next {
		if e.key == name {
			e.sym.Value = value
			return nil
		}
	}
	return fmt.Errorf("symbole '%s' non trouvé", name)
}

// Print affiche la table de symboles au format tabulaire
func (st *SymbolTable) Print() {
	line := strings.Repeat("─", 72)
	fmt.Println("\n┌" + line + "┐")
	fmt.Printf("│%s│\n", center("TABLE DES SYMBOLES", 72))
	fmt.Println("├──────┬──────────────┬───────────┬─────────┬───────┬────────────┤")
	fmt.Printf("│ %-4s │ %-12s │ %-9s │ %-7s │ %-5s │ %-10s │\n",
		"Idx", "Nom", "Catégorie", "Type", "Ligne", "Valeur")
	fmt.Println("├──────┼──────────────┼───────────┼─────────┼───────┼────────────┤")

	for i, e := range st.buckets {
		for ; e != nil; e = e.next {
			v := e.sym.Value
			if e.sym.Kind == KIND_ARRAY {
				v = fmt.Sprintf("[%d]", e.sym.ArraySize)
			}
			fmt.Printf("│ %-4d │ %-12s │ %-9s │ %-7s │ %-5d │ %-10s │\n",
				i, e.sym.Name, e.sym.Kind, e.sym.Type, e.sym.Line, v)
		}
	}

	fmt.Println("└──────┴──────────────┴───────────┴─────────┴───────┴────────────┘")
	fmt.Printf("  Total : %d symboles\n\n", st.count)
}

// center centre une chaîne dans une largeur donnée
func center(s string, w int) string {
	pad := (w - len(s)) / 2
	return strings.Repeat(" ", pad) + s + strings.Repeat(" ", w-len(s)-pad)
}

// ════════════════════════════════════════════════════════════════
//  TYPE UTILITIES
// ════════════════════════════════════════════════════════════════

// MergeTypes détermine le type résultant de l'opération sur deux types
func MergeTypes(a, b DataType) DataType {
	if a == TYPE_FLOAT || b == TYPE_FLOAT {
		return TYPE_FLOAT
	}
	if a == TYPE_UNKNOWN || b == TYPE_UNKNOWN {
		return TYPE_UNKNOWN
	}
	return TYPE_INTEGER
}
