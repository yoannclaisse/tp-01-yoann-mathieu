package main

import (
	"flag"
	"fmt"
	"os"
	"tp1/annuaire"
)

func main() {
	var action = flag.String("action", "", "Action à effectuer (ajouter)")
	var nom = flag.String("nom", "", "Nom du contact")
	var prenom = flag.String("prenom", "", "Prénom du contact")
	var tel = flag.String("tel", "", "Numéro de téléphone")

	flag.Parse()

	ann := annuaire.NewAnnuaire()

	switch *action {
	case "ajouter":
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

	case "":
		flag.PrintDefaults()
	default:
		fmt.Printf("Action '%s' non implémentée\n", *action)
		os.Exit(1)
	}
}
