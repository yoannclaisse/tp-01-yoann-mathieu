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

// Recherche un contact par nom
func (a *Annuaire) RechercherContact(nom string) (Contact, bool) {
	contact, existe := a.contacts[nom]
	return contact, existe
}

// Retourne tous les contacts
func (a *Annuaire) ListerContacts() []Contact {
	contacts := make([]Contact, 0, len(a.contacts))
	for _, contact := range a.contacts {
		contacts = append(contacts, contact)
	}
	return contacts
}
