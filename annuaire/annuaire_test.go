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

func TestSupprimerContact(t *testing.T) {
	ann := NewAnnuaire()
	ann.AjouterContact("Test", "Contact", "0000000000")

	err := ann.SupprimerContact("Test")
	if err != nil {
		t.Errorf("Erreur lors de la suppression: %v", err)
	}

	if ann.NombreContacts() != 0 {
		t.Error("Le contact n'a pas été supprimé")
	}

	err = ann.SupprimerContact("Inexistant")
	if err == nil {
		t.Error("Attendu une erreur pour un contact inexistant")
	}
}

func TestModifierContact(t *testing.T) {
	ann := NewAnnuaire()
	ann.AjouterContact("Modif", "Test", "0000000000")

	err := ann.ModifierContact("Modif", "NouveauPrenom", "1111111111")
	if err != nil {
		t.Errorf("Erreur lors de la modification: %v", err)
	}

	contact, _ := ann.RechercherContact("Modif")
	if contact.Prenom != "NouveauPrenom" || contact.Telephone != "1111111111" {
		t.Errorf("Modification échouée: %+v", contact)
	}
}
