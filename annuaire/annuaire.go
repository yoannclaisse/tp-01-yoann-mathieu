package annuaire

import "errors"

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

// Ajoute un nouveau contact à l'annuaire
func (a *Annuaire) AjouterContact(nom, prenom, telephone string) error {
	if nom == "" || prenom == "" || telephone == "" {
		return errors.New("tous les champs sont requis")
	}

	if _, existe := a.contacts[nom]; existe {
		return errors.New("un contact avec ce nom existe déjà")
	}

	a.contacts[nom] = Contact{
		Nom:       nom,
		Prenom:    prenom,
		Telephone: telephone,
	}

	return nil
}
