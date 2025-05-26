package main

import (
	"flag"
	"fmt"
	"os"
	"tp1/annuaire"
)

func main() {
	var action = flag.String("action", "", "Action to perform (add, list, search)")
	var nom = flag.String("nom", "", "Nom du contact")
	var prenom = flag.String("prenom", "", "Prénom du contact")
	var tel = flag.String("tel", "", "Numéro de téléphone")
	var fichier = flag.String("fichier", "data/contacts.json", "Fichier JSON pour import/export")

	flag.Parse()

	ann := annuaire.NewAnnuaire()
	ann.ChargerDepuisJSON(*fichier)

	switch *action {
	case "add":
		if *nom == "" || *prenom == "" || *tel == "" {
			fmt.Println("Erreur: nom, prénom et téléphone requis")
			os.Exit(1)
		}
		err := ann.AjouterContact(*nom, *prenom, *tel)
		if err != nil {
			fmt.Printf("Erreur: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Contact %s %s ajouté avec succès\n", *prenom, *nom)
	case "list":
		contacts := ann.ListerContacts()
		if len(contacts) == 0 {
			fmt.Println("Aucun contact trouvé")
		} else {
			fmt.Println("Liste des contacts:")
			for _, contact := range contacts {
				fmt.Printf("- %s %s: %s\n", contact.Prenom, contact.Nom, contact.Telephone)
			}
		}
	case "search":
		if *nom == "" {
			fmt.Println("Erreur: nom requis")
			os.Exit(1)
		}
		contact, existe := ann.RechercherContact(*nom)
		if existe {
			fmt.Printf("Contact trouvé: %s %s - %s\n", contact.Prenom, contact.Nom, contact.Telephone)
		} else {
			fmt.Printf("Aucun contact trouvé avec le nom: %s\n", *nom)
		}
	case "supprimer":
		if *nom == "" {
			fmt.Println("Erreur: nom requis")
			os.Exit(1)
		}
		err := ann.SupprimerContact(*nom)
		if err != nil {
			fmt.Printf("Erreur: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Contact %s supprimé avec succès\n", *nom)
	case "modifier":
		if *nom == "" {
			fmt.Println("Erreur: nom requis")
			os.Exit(1)
		}
		err := ann.ModifierContact(*nom, *prenom, *tel)
		if err != nil {
			fmt.Printf("Erreur: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Contact %s modifié avec succès\n", *nom)
	case "export":
		err := ann.SauvegarderEnJSON(*fichier)
		if err != nil {
			fmt.Printf("Erreur lors de l'export: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Contacts exportés vers %s\n", *fichier)
	case "import":
		err := ann.ChargerDepuisJSON(*fichier)
		if err != nil {
			fmt.Printf("Erreur lors de l'import: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Contacts importés depuis %s\n", *fichier)
		if *action != "" && *action != "list" && *action != "search" {
			ann.SauvegarderEnJSON(*fichier)
		}
	case "":
		flag.PrintDefaults()
	default:
		fmt.Printf("Action '%s' non implémentée\n", *action)
		os.Exit(1)
	}
}
