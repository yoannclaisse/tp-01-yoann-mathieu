package annuaire

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

func (a *Annuaire) SupprimerContact(nom string) error {
	if _, existe := a.contacts[nom]; !existe {
		return errors.New("contact non trouvé")
	}

	delete(a.contacts, nom)
	return nil
}

func (a *Annuaire) ModifierContact(nom, nouveauPrenom, nouveauTelephone string) error {
	contact, existe := a.contacts[nom]
	if !existe {
		return errors.New("contact non trouvé")
	}

	if nouveauPrenom != "" {
		contact.Prenom = nouveauPrenom
	}
	if nouveauTelephone != "" {
		contact.Telephone = nouveauTelephone
	}

	a.contacts[nom] = contact
	return nil
}

func (a *Annuaire) NombreContacts() int {
	return len(a.contacts)
}

// SauvegarderEnJSON sauvegarde l'annuaire dans un fichier JSON
func (a *Annuaire) SauvegarderEnJSON(nomFichier string) error {
	dir := filepath.Dir(nomFichier)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(a.contacts, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(nomFichier, data, 0644)
}

// ChargerDepuisJSON charge l'annuaire depuis un fichier JSON
func (a *Annuaire) ChargerDepuisJSON(nomFichier string) error {
	if _, err := os.Stat(nomFichier); os.IsNotExist(err) {
		return nil
	}

	data, err := ioutil.ReadFile(nomFichier)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &a.contacts)
}
