package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"tp1/annuaire"
)

// Variables globales pour l'annuaire et le fichier de données
var ann *annuaire.Annuaire
var fichierDonnees = "data/contacts.json"

// Template HTML pour l'interface web
const htmlTemplate = `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Annuaire Go - Interface Web</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1, h2 {
            color: #333;
        }
        form {
            background-color: #f9f9f9;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 20px;
        }
        input[type="text"] {
            width: 200px;
            padding: 8px;
            margin: 5px;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        button {
            background-color: #007bff;
            color: white;
            padding: 10px 15px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            margin: 5px;
        }
        button:hover {
            background-color: #0056b3;
        }
        .contact {
            background-color: #f8f9fa;
            padding: 10px;
            margin: 5px 0;
            border-left: 4px solid #007bff;
            border-radius: 4px;
        }
        .error {
            color: #dc3545;
            background-color: #f8d7da;
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .success {
            color: #155724;
            background-color: #d4edda;
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .stats {
            background-color: #e9ecef;
            padding: 10px;
            border-radius: 4px;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📞 Annuaire Go - Interface Web</h1>
        
        <div class="stats">
            <strong>Nombre de contacts :</strong> {{.NombreContacts}}
        </div>

        {{if .Message}}
            <div class="{{.MessageType}}">{{.Message}}</div>
        {{end}}

        <h2>Ajouter un contact</h2>
        <form action="/ajouter" method="POST">
            <input type="text" name="nom" placeholder="Nom" required>
            <input type="text" name="prenom" placeholder="Prénom" required>
            <input type="text" name="telephone" placeholder="Téléphone" required>
            <button type="submit">Ajouter</button>
        </form>

        <h2>Rechercher un contact</h2>
        <form action="/rechercher" method="GET">
            <input type="text" name="nom" placeholder="Nom à rechercher" required>
            <button type="submit">Rechercher</button>
        </form>

        {{if .ContactRecherche}}
        <div class="contact">
            <strong>Contact trouvé :</strong><br>
            {{.ContactRecherche.Prenom}} {{.ContactRecherche.Nom}} - {{.ContactRecherche.Telephone}}
            <form action="/supprimer" method="POST" style="display: inline;">
                <input type="hidden" name="nom" value="{{.ContactRecherche.Nom}}">
                <button type="submit" onclick="return confirm('Êtes-vous sûr de vouloir supprimer ce contact ?')">Supprimer</button>
            </form>
        </div>
        {{end}}

        <h2>Liste des contacts</h2>
        {{if .Contacts}}
            {{range .Contacts}}
            <div class="contact">
                <strong>{{.Prenom}} {{.Nom}}</strong> - {{.Telephone}}
                <form action="/supprimer" method="POST" style="display: inline; float: right;">
                    <input type="hidden" name="nom" value="{{.Nom}}">
                    <button type="submit" onclick="return confirm('Êtes-vous sûr de vouloir supprimer ce contact ?')">Supprimer</button>
                </form>
            </div>
            {{end}}
        {{else}}
            <p>Aucun contact dans l'annuaire.</p>
        {{end}}

        <hr>
        <h2>Actions</h2>
        <form action="/exporter" method="POST" style="display: inline;">
            <button type="submit">Exporter en JSON</button>
        </form>
        
        <form action="/vider" method="POST" style="display: inline;">
            <button type="submit" onclick="return confirm('Êtes-vous sûr de vouloir vider l\'annuaire ?')">Vider l'annuaire</button>
        </form>
    </div>
</body>
</html>
`

// Structure pour passer les données au template
type PageData struct {
	Contacts         []annuaire.Contact
	ContactRecherche *annuaire.Contact
	Message          string
	MessageType      string
	NombreContacts   int
}

// StartServer démarre le serveur web
func StartServer() {
	// Initialiser l'annuaire et charger les données existantes
	ann = annuaire.NewAnnuaire()
	err := ann.ChargerDepuisJSON(fichierDonnees)
	if err != nil {
		log.Printf("Impossible de charger les données : %v", err)
	}

	// Routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ajouter", handleAjouter)
	http.HandleFunc("/rechercher", handleRechercher)
	http.HandleFunc("/supprimer", handleSupprimer)
	http.HandleFunc("/exporter", handleExporter)
	http.HandleFunc("/vider", handleVider)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleHome affiche la page principale
func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("home").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Erreur de template", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Contacts:       ann.ListerContacts(),
		NombreContacts: ann.NombreContacts(),
	}

	tmpl.Execute(w, data)
}

// handleAjouter traite l'ajout d'un nouveau contact
func handleAjouter(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	nom := r.FormValue("nom")
	prenom := r.FormValue("prenom")
	telephone := r.FormValue("telephone")

	err := ann.AjouterContact(nom, prenom, telephone)

	tmpl, _ := template.New("home").Parse(htmlTemplate)
	data := PageData{
		Contacts:       ann.ListerContacts(),
		NombreContacts: ann.NombreContacts(),
	}

	if err != nil {
		data.Message = fmt.Sprintf("Erreur : %v", err)
		data.MessageType = "error"
	} else {
		data.Message = fmt.Sprintf("Contact %s %s ajouté avec succès", prenom, nom)
		data.MessageType = "success"
		// Sauvegarder automatiquement
		ann.SauvegarderEnJSON(fichierDonnees)
	}

	tmpl.Execute(w, data)
}

// handleRechercher traite la recherche d'un contact
func handleRechercher(w http.ResponseWriter, r *http.Request) {
	nom := r.FormValue("nom")

	tmpl, _ := template.New("home").Parse(htmlTemplate)
	data := PageData{
		Contacts:       ann.ListerContacts(),
		NombreContacts: ann.NombreContacts(),
	}

	if nom != "" {
		if contact, existe := ann.RechercherContact(nom); existe {
			data.ContactRecherche = &contact
			data.Message = "Contact trouvé"
			data.MessageType = "success"
		} else {
			data.Message = fmt.Sprintf("Aucun contact trouvé avec le nom : %s", nom)
			data.MessageType = "error"
		}
	}

	tmpl.Execute(w, data)
}

// handleSupprimer traite la suppression d'un contact
func handleSupprimer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	nom := r.FormValue("nom")
	err := ann.SupprimerContact(nom)

	tmpl, _ := template.New("home").Parse(htmlTemplate)
	data := PageData{
		Contacts:       ann.ListerContacts(),
		NombreContacts: ann.NombreContacts(),
	}

	if err != nil {
		data.Message = fmt.Sprintf("Erreur : %v", err)
		data.MessageType = "error"
	} else {
		data.Message = fmt.Sprintf("Contact %s supprimé avec succès", nom)
		data.MessageType = "success"
		// Sauvegarder automatiquement
		ann.SauvegarderEnJSON(fichierDonnees)
	}

	tmpl.Execute(w, data)
}

// handleExporter exporte les données en JSON
func handleExporter(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := ann.SauvegarderEnJSON(fichierDonnees)

	tmpl, _ := template.New("home").Parse(htmlTemplate)
	data := PageData{
		Contacts:       ann.ListerContacts(),
		NombreContacts: ann.NombreContacts(),
	}

	if err != nil {
		data.Message = fmt.Sprintf("Erreur lors de l'export : %v", err)
		data.MessageType = "error"
	} else {
		data.Message = "Données exportées avec succès"
		data.MessageType = "success"
	}

	tmpl.Execute(w, data)
}

// handleVider vide complètement l'annuaire
func handleVider(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Créer un nouvel annuaire vide
	ann = annuaire.NewAnnuaire()
	err := ann.SauvegarderEnJSON(fichierDonnees)

	tmpl, _ := template.New("home").Parse(htmlTemplate)
	data := PageData{
		Contacts:       ann.ListerContacts(),
		NombreContacts: ann.NombreContacts(),
	}

	if err != nil {
		data.Message = fmt.Sprintf("Erreur lors du vidage : %v", err)
		data.MessageType = "error"
	} else {
		data.Message = "Annuaire vidé avec succès"
		data.MessageType = "success"
	}

	tmpl.Execute(w, data)
}
