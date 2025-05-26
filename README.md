# 📞 Annuaire Go

Un annuaire de contacts simple et élégant développé en Go, avec interface en ligne de commande et serveur web intégré.

## 📋 Table des matières

- [Fonctionnalités](#-fonctionnalités)
- [Installation](#-installation)
- [Utilisation](#-utilisation)
  - [Interface en ligne de commande](#interface-en-ligne-de-commande)
  - [Interface web](#interface-web)
- [Exemples](#-exemples)
- [Structure du projet](#-structure-du-projet)
- [Tests](#-tests)
- [API/Package](#-apipackage)
- [Développement](#-développement)
- [Contributeurs](#-contributeurs)

## ✨ Fonctionnalités

### Interface CLI

- ✅ **Ajouter** des contacts (nom, prénom, téléphone)
- 🔍 **Rechercher** un contact par nom
- 📝 **Lister** tous les contacts
- ✏️ **Modifier** les informations d'un contact
- 🗑️ **Supprimer** un contact
- 💾 **Persistance** automatique en JSON
- 📤 **Import/Export** de données JSON

### Interface Web

- 🌐 **Interface web moderne** accessible via navigateur
- 📱 **Design responsive** et intuitif
- ⚡ **Actions en temps réel** (ajout, suppression, recherche)
- 📊 **Statistiques** du nombre de contacts
- 💬 **Messages de confirmation** et d'erreur
- 🔒 **Confirmations** avant suppression

### Fonctionnalités techniques

- 🧪 **Tests unitaires** complets
- 📦 **Architecture modulaire** (package séparé)
- 🛡️ **Gestion d'erreurs** robuste
- 💽 **Sauvegarde automatique** des modifications

## 🚀 Installation

### Prérequis

- Go 1.21 ou plus récent
- Git

### Cloner le projet

```bash
git clone <url-du-repo>
cd annuaire
```

### Compiler le projet

```bash
go mod tidy
go build -o annuaire main.go
```

## 💻 Utilisation

### Interface en ligne de commande

#### Afficher l'aide

```bash
./annuaire
```

#### Ajouter un contact

```bash
./annuaire -action=ajouter -nom="Dupont" -prenom="Jean" -tel="0123456789"
```

#### Lister tous les contacts

```bash
./annuaire -action=lister
```

#### Rechercher un contact

```bash
./annuaire -action=rechercher -nom="Dupont"
```

#### Modifier un contact

```bash
./annuaire -action=modifier -nom="Dupont" -prenom="Pierre" -tel="0987654321"
```

#### Supprimer un contact

```bash
./annuaire -action=supprimer -nom="Dupont"
```

#### Utiliser un fichier JSON personnalisé

```bash
./annuaire -action=lister -fichier="mes_contacts.json"
```

### Interface web

#### Lancer le serveur web

```bash
./annuaire -serveur
```

Puis ouvrir votre navigateur sur : **<http://localhost:8080>**

L'interface web permet de :

- Visualiser tous les contacts en temps réel
- Ajouter de nouveaux contacts via un formulaire
- Rechercher des contacts spécifiques
- Supprimer des contacts avec confirmation
- Exporter/vider l'annuaire
- Voir les statistiques du nombre de contacts

## 📚 Exemples

### Scénario d'utilisation complet

```bash
# 1. Ajouter quelques contacts
./annuaire -action=ajouter -nom="Martin" -prenom="Alice" -tel="0123456789"
./annuaire -action=ajouter -nom="Durand" -prenom="Bob" -tel="0987654321"
./annuaire -action=ajouter -nom="Moreau" -prenom="Claire" -tel="0555123456"

# 2. Lister tous les contacts
./annuaire -action=lister

# 3. Rechercher un contact spécifique
./annuaire -action=rechercher -nom="Martin"

# 4. Modifier un contact
./annuaire -action=modifier -nom="Martin" -tel="0111222333"

# 5. Supprimer un contact
./annuaire -action=supprimer -nom="Durand"

# 6. Lancer l'interface web
./annuaire -serveur
```

### Utilisation avec fichiers personnalisés

```bash
# Créer un annuaire de travail séparé
./annuaire -action=ajouter -nom="Boss" -prenom="Manager" -tel="0100000000" -fichier="travail.json"

# Créer un annuaire personnel
./annuaire -action=ajouter -nom="Maman" -prenom="Chérie" -tel="0200000000" -fichier="famille.json"
```

## 🏗️ Structure du projet

```
annuaire/
├── main.go                    # Point d'entrée principal
├── go.mod                     # Dépendances Go
├── README.md                  # Ce fichier
├── annuaire/                  # Package principal
│   ├── annuaire.go           # Logique métier et structures
│   └── annuaire_test.go      # Tests unitaires
├── server/                    # Serveur web
│   └── server.go             # Interface web et routes HTTP
└── data/                      # Données persistantes
    └── contacts.json         # Fichier de sauvegarde par défaut
```

### Description des composants

#### `main.go`

Point d'entrée qui gère :

- Parsing des arguments de ligne de commande
- Routage vers les actions CLI ou serveur web
- Gestion des erreurs globales

#### `annuaire/annuaire.go`

Package principal contenant :

- Structures `Contact` et `Annuaire`
- Opérations CRUD (Create, Read, Update, Delete)
- Fonctions de persistance JSON
- Gestion des erreurs métier

#### `server/server.go`

Serveur web HTTP avec :

- Templates HTML intégrés
- Routes pour toutes les opérations
- Interface utilisateur moderne
- Gestion des formulaires et réponses

## 🧪 Tests

### Lancer tous les tests

```bash
go test ./...
```

### Tests avec couverture

```bash
go test -cover ./annuaire
```

### Tests détaillés

```bash
go test -v ./annuaire
```

### Tests couverts

- ✅ Ajout de contacts (valides et invalides)
- ✅ Recherche de contacts (existants et inexistants)
- ✅ Suppression de contacts
- ✅ Modification de contacts
- ✅ Persistance JSON (sauvegarde et chargement)
- ✅ Gestion des erreurs métier

## 📖 API/Package

### Structure `Contact`

```go
type Contact struct {
    Nom       string
    Prenom    string
    Telephone string
}
```

### Structure `Annuaire`

```go
type Annuaire struct {
    contacts map[string]Contact
}
```

### Méthodes principales

```go
// Création
func NewAnnuaire() *Annuaire

// Opérations CRUD
func (a *Annuaire) AjouterContact(nom, prenom, telephone string) error
func (a *Annuaire) RechercherContact(nom string) (Contact, bool)
func (a *Annuaire) ListerContacts() []Contact
func (a *Annuaire) ModifierContact(nom, nouveauPrenom, nouveauTelephone string) error
func (a *Annuaire) SupprimerContact(nom string) error

// Persistance
func (a *Annuaire) SauvegarderEnJSON(nomFichier string) error
func (a *Annuaire) ChargerDepuisJSON(nomFichier string) error

// Utilitaires
func (a *Annuaire) NombreContacts() int
```

## 🛠️ Développement

### Architecture

- **Séparation des préoccupations** : logique métier, interface CLI, interface web
- **Tests unitaires** pour valider la logique métier
- **Gestion d'erreurs** cohérente dans tout le projet
- **Code Go idiomatique** avec conventions standards

### Standards de code

- Variables et fonctions en français pour la cohérence
- Gestion explicite des erreurs
- Documentation des fonctions publiques
- Tests pour toutes les fonctions métier

### Ajout de fonctionnalités

1. Ajouter la logique métier dans `annuaire/annuaire.go`
2. Ajouter les tests correspondants dans `annuaire/annuaire_test.go`
3. Mettre à jour l'interface CLI dans `main.go`
4. Mettre à jour l'interface web dans `server/server.go`

### Extension possible

- 🔍 Recherche par prénom ou téléphone
- 📱 Validation des numéros de téléphone
- 📧 Ajout d'email et adresse
- 🏷️ Système de tags/catégories
- 🔐 Authentification pour l'interface web
- 📊 Export en CSV/Excel
- 🌍 API REST pour intégrations externes

## 👥 Contributeurs

Ce projet a été développé collaborativement par :

- **Mathieu** - Développement principal, interface CLI, intégration
- **Yoann** - Structures de données, tests, interface web

### Historique des contributions

- Setup initial et architecture du projet
- Développement des fonctionnalités de base (CRUD)
- Implémentation des tests unitaires
- Ajout de la persistance JSON
- Création de l'interface web moderne
- Documentation et finalisation

---

## 📄 Licence

Ce projet est développé dans le cadre d'un TP académique.

## 🚀 Version

**Version actuelle :** 1.0.0

**Fonctionnalités :** CLI complète + Interface web + Tests + Persistance JSON

---
