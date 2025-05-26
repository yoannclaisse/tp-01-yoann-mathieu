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

// Create empty annuaire
func NewAnnuaire() *Annuaire {
	return &Annuaire{
		contacts: make(map[string]Contact),
	}
}

// Add contact
func (a *Annuaire) AjouterContact(nom, prenom, telephone string) error {
	if nom == "" || prenom == "" || telephone == "" {
		return errors.New("tous les champs sont requis")
	}

	for _, c := range a.contacts {
		if c.Nom == nom && c.Telephone == telephone {
			return errors.New("un contact avec ce nom et ce téléphone existe déjà")
		}
	}

	a.contacts[nom] = Contact{
		Nom:       nom,
		Prenom:    prenom,
		Telephone: telephone,
	}

	return nil
}

// Search
func (a *Annuaire) RechercherContact(nom string) (Contact, bool) {
	contact, existe := a.contacts[nom]
	return contact, existe
}

// List contacts
func (a *Annuaire) ListerContacts() []Contact {
	contacts := make([]Contact, 0, len(a.contacts))
	for _, contact := range a.contacts {
		contacts = append(contacts, contact)
	}
	return contacts
}
