package annuaire

import (
	"testing"
)

func TestAjouterContact(t *testing.T) {
	ann := NewAnnuaire()

	// ajout
	err := ann.AjouterContact("Dupont", "Jean", "0123456789")
	if err != nil {
		t.Errorf("Erreur lors de l'ajout: %v", err)
	}

	// ajout contact
	err = ann.AjouterContact("Dupont", "Pierre", "0987654321")
	if err != nil {
		t.Errorf("Erreur inattendue lors de l'ajout d'un contact avec même nom mais téléphone différent: %v", err)
	}

	// ajout avec champs vides
	err = ann.AjouterContact("", "Test", "0000000000")
	if err == nil {
		t.Error("Attendu une erreur pour un nom vide")
	}
}

func TestRechercherContact(t *testing.T) {
	ann := NewAnnuaire()
	ann.AjouterContact("Martin", "Alice", "0123456789")

	// recherche existant
	contact, existe := ann.RechercherContact("Martin")
	if !existe {
		t.Error("Contact non trouvé alors qu'il devrait exister")
	}

	if contact.Prenom != "Alice" || contact.Telephone != "0123456789" {
		t.Errorf("Données incorrectes: %+v", contact)
	}

	// recherche inexistant
	_, existe = ann.RechercherContact("Inexistant")
	if existe {
		t.Error("Contact trouvé alors qu'il ne devrait pas exister")
	}
}
