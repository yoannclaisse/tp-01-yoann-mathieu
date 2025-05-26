package main

import (
	"flag"
	"fmt"
	"os"
	"tp1/annuaire"
)

func main() {
	var action = flag.String("action", "", "Action à effectuer")
	var nom = flag.String("nom", "", "Nom du contact")
	var prenom = flag.String("prenom", "", "Prénom du contact")
	var tel = flag.String("tel", "", "Numéro de téléphone")

	flag.Parse()

	ann := annuaire.NewAnnuaire()

	fmt.Printf("nom: %s, prenom: %s, tel: %s, ann: %v\n", *nom, *prenom, *tel, ann)
	switch *action {
	case "":
		flag.PrintDefaults()
	default:
		fmt.Printf("Action '%s' non implémentée\n", *action)
		os.Exit(1)
	}
}
