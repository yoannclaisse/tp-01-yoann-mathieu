package annuaire

// Structure de contact
type Contact struct {
	Nom       string
	Prenom    string
	Telephone string
}

// Annuaire qui gère les contacts
type Annuaire struct {
	contacts map[string]Contact
}

// Crée un nouvel annuaire vide
func NewAnnuaire() *Annuaire {
	return &Annuaire{
		contacts: make(map[string]Contact),
	}
}
