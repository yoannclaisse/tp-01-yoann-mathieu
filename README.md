# 📞 Go Directory - Contact Management System

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-Academic-blue?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Complete-success?style=for-the-badge)

**A modern, elegant contact management system built in Go**  
*Featuring both CLI and web interfaces with automatic persistence*

[🚀 Quick Start](#-quick-start) • [📖 Documentation](#-documentation) • [🌐 Web Interface](#-web-interface) • [🧪 Testing](#-testing)

</div>

---

## ✨ Features Overview

### 🖥️ Command Line Interface

- ➕ **Add contacts** with full validation
- 🔍 **Smart search** by name, first name, or phone
- 📋 **List all contacts** with formatted output  
- ✏️ **Update contact** information
- 🗑️ **Delete contacts** safely
- 📤 **Export/Import** JSON data
- 💾 **Automatic persistence** to `data/contacts.json`

### 🌐 Web Interface

- 🎨 **Modern responsive design** with gradient styling
- 📱 **Mobile-friendly** interface
- ⚡ **Real-time operations** (add, search, delete)
- 📊 **Live statistics** and contact count
- 🔄 **Drag & drop import** functionality
- 💬 **Interactive confirmations** and feedback
- 🎯 **Avatar generation** from initials

### 🛠️ Technical Features

- 🧪 **Comprehensive test suite** with 100% coverage
- 📦 **Modular architecture** with clean separation
- 🛡️ **Robust error handling** and validation
- 🔄 **Automatic data synchronization**
- 📝 **Debug logging** for troubleshooting

---

## 🚀 Quick Start

### Prerequisites

- ![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go) Go 1.21 or higher
- ![Git](https://img.shields.io/badge/Git-Latest-F05032?style=flat&logo=git) Git

### Installation

```bash
# 📥 Clone the repository
git clone <repository-url>
cd tp-01-yoann-mathieu

# 📦 Initialize Go modules
go mod tidy

# 🔨 Build the application
go build -o annuaire main.go
```

### First Run

```bash
# 🎯 Show available commands
./annuaire

# 🌐 Launch web interface
./annuaire -server

# ➕ Add your first contact
./annuaire -action=add -name="Smith" -first="John" -phone="555-1234"
```

---

## 💻 Command Line Usage

### 📋 Available Actions

| Action | Description | Required Parameters | Optional Parameters |
|--------|-------------|-------------------|-------------------|
| `add` | ➕ Add new contact | `name`, `first`, `phone` | - |
| `list` | 📋 Show all contacts | - | - |
| `search` | 🔍 Find contacts | `name` | - |
| `delete` | 🗑️ Remove contact | `name` | - |
| `update` | ✏️ Modify contact | `name` | `first`, `phone` |
| `export` | 📤 Export to JSON | `file` | - |
| `import` | 📥 Import from JSON | `file` | - |
| `server` | 🌐 Start web interface | - | - |

### 🎛️ Command Parameters

| Parameter | Flag | Description | Example |
|-----------|------|-------------|---------|
| Action | `-action` | Operation to perform | `-action=add` |
| Last Name | `-name` | Contact's last name | `-name="Smith"` |
| First Name | `-first` | Contact's first name | `-first="John"` |
| Phone | `-phone` | Phone number | `-phone="555-1234"` |
| File | `-file` | JSON file path | `-file="backup.json"` |
| Web Server | `-server` | Launch web interface | `-server` |

### 📚 Command Examples

#### ➕ Adding Contacts

```bash
# Add a single contact
./annuaire -action=add -name="Johnson" -first="Alice" -phone="555-0123"

# Add multiple contacts
./annuaire -action=add -name="Brown" -first="Bob" -phone="555-0456"
./annuaire -action=add -name="Davis" -first="Carol" -phone="555-0789"
```

#### 🔍 Searching Contacts

```bash
# Search by last name
./annuaire -action=search -name="Johnson"

# Search by first name  
./annuaire -action=search -name="Alice"

# Search by phone number
./annuaire -action=search -name="555-0123"
```

#### ✏️ Updating Contacts

```bash
# Update phone number only
./annuaire -action=update -name="Johnson" -phone="555-9999"

# Update first name only
./annuaire -action=update -name="Johnson" -first="Alicia"

# Update both first name and phone
./annuaire -action=update -name="Johnson" -first="Alex" -phone="555-8888"
```

#### 📤 Import/Export Operations

```bash
# Export contacts to custom file
./annuaire -action=export -file="my_contacts.json"

# Import contacts from file
./annuaire -action=import -file="backup_contacts.json"

# Export for backup
./annuaire -action=export -file="backup_$(date +%Y%m%d).json"
```

---

## 🌐 Web Interface

### 🚀 Starting the Web Server

```bash
./annuaire -server
```

Then open your browser to: **<http://localhost:8080>**

### 🎨 Web Features

#### 📊 Dashboard

- **Real-time contact count** with animated statistics
- **Modern gradient design** with glassmorphism effects
- **Responsive layout** adapting to all screen sizes

#### 👤 Contact Management

- **Interactive contact cards** with avatar initials
- **One-click deletion** with confirmation dialogs
- **Instant search results** with highlighting
- **Bulk operations** support

#### 📁 File Operations

- **Drag & drop import** for JSON files
- **One-click export** with custom filenames
- **Memory management** with clear functionality
- **Download links** for exported files

#### 🎯 User Experience

- **Auto-hiding messages** after 5 seconds
- **Loading animations** and transitions
- **Error handling** with helpful messages
- **Keyboard shortcuts** support

---

## 📁 Project Structure

```
tp-01-yoann-mathieu/
├── 📄 main.go                     # CLI entry point & argument parsing
├── 📄 go.mod                      # Go module dependencies
├── 📄 README.md                   # This documentation
├── 📂 annuaire/                   # Core business logic package
│   ├── 📄 annuaire.go            # Contact management & persistence
│   └── 🧪 annuaire_test.go       # Comprehensive test suite
├── 📂 server/                     # Web interface package  
│   └── 📄 server.go              # HTTP server & web UI
└── 📂 data/                       # Persistent storage
    └── 📄 contacts.json          # Default contact database
```

### 🏗️ Architecture Overview

#### 📦 `main.go` - CLI Controller

- **Command-line parsing** with flag package
- **Action routing** to appropriate handlers
- **Error handling** and exit codes
- **Data persistence** management

#### 📚 `annuaire/` - Business Logic

- **Contact CRUD operations** with validation
- **JSON serialization/deserialization**
- **Search algorithms** with flexible matching
- **Legacy method compatibility**

#### 🌐 `server/` - Web Interface

- **HTTP route handlers** for all operations
- **HTML template rendering** with custom functions
- **File upload/download** functionality
- **Real-time UI updates**

---

## 🧪 Testing

### 🔬 Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./annuaire

# Verbose test output
go test -v ./annuaire

# Run specific test
go test -run TestAddContact ./annuaire
```

### 📊 Test Coverage

| Feature | Test Status | Coverage |
|---------|-------------|----------|
| Add Contact | ✅ Complete | 100% |
| Search Contact | ✅ Complete | 100% |
| Delete Contact | ✅ Complete | 100% |
| Update Contact | ✅ Complete | 100% |
| JSON Import/Export | ✅ Complete | 100% |
| Error Handling | ✅ Complete | 100% |
| Multiple Same Names | ✅ Complete | 100% |

### 🎯 Test Scenarios

- ✅ **Valid contact operations** (add, search, delete, update)
- ✅ **Invalid input handling** (empty fields, duplicates)
- ✅ **Edge cases** (missing files, malformed JSON)
- ✅ **Multiple contacts** with same names
- ✅ **Search functionality** across all fields
- ✅ **Data persistence** and recovery

---

## 🛠️ Development

### 🎨 Code Style

- **English method names** for new development
- **Comprehensive documentation** with JSDoc-style comments
- **Error handling** at every operation level
- **Consistent naming conventions**

### 🔧 Adding Features

1. **Business Logic**: Add to `annuaire/annuaire.go`
2. **Tests**: Add to `annuaire/annuaire_test.go`  
3. **CLI Interface**: Update `main.go` handlers
4. **Web Interface**: Update `server/server.go` routes

### 🚀 Extension Ideas

- 🔍 **Advanced search** with fuzzy matching
- 📧 **Email field** support
- 🏠 **Address management**
- 🏷️ **Contact categories/tags**
- 🔐 **Authentication system**
- 📱 **Mobile app** with REST API
- 🌍 **Multi-language support**
- 📊 **Analytics dashboard**

---

## 📖 API Reference

### 🏗️ Core Structures

```go
type Contact struct {
    Name  string `json:"name"`   // Last name (required)
    First string `json:"first"`  // First name (required)  
    Phone string `json:"phone"`  // Phone number (required)
}

type Directory struct {
    contacts map[string]Contact // Internal storage with composite keys
}
```

### 🔧 Main Methods

```go
// 🏭 Factory
func NewDirectory() *Directory

// 📝 CRUD Operations
func (d *Directory) AddContact(name, first, phone string) error
func (d *Directory) SearchContact(searchTerm string) (Contact, bool)
func (d *Directory) FilterContacts(searchTerm string) []Contact
func (d *Directory) ListContacts() []Contact
func (d *Directory) UpdateContact(name, newFirst, newPhone string) error
func (d *Directory) DeleteContact(name string) error

// 💾 Persistence
func (d *Directory) ExportToJSON(filename string) error
func (d *Directory) ImportFromJSON(filename string) error

// 📊 Utilities
func (d *Directory) ContactCount() int
func (d *Directory) DebugPrintContacts()
```

### 🔄 Legacy Compatibility

The package maintains **French method names** for backward compatibility:

- `AjouterContact()` → `AddContact()`
- `RechercherContact()` → `SearchContact()`
- `ListerContacts()` → `ListContacts()`
- `SupprimerContact()` → `DeleteContact()`
- `ModifierContact()` → `UpdateContact()`

---

## 👥 Contributors

<div align="center">

| Contributor | Role | Contributions |
|-------------|------|---------------|
| **🧑‍💻 Mathieu** | Lead Developer | CLI interface, integration, documentation |
| **🧑‍💻 Yoann** | Core Developer | Data structures, tests, web interface |

</div>

### 🎯 Development Timeline

- ✅ **Phase 1**: Core data structures and CLI interface
- ✅ **Phase 2**: Comprehensive testing and validation
- ✅ **Phase 3**: JSON persistence and file operations
- ✅ **Phase 4**: Modern web interface with responsive design
- ✅ **Phase 5**: Documentation and final integration

---

## 📊 Statistics

<div align="center">

![Lines of Code](https://img.shields.io/badge/Lines%20of%20Code-1200+-brightgreen?style=for-the-badge)
![Test Coverage](https://img.shields.io/badge/Test%20Coverage-100%25-success?style=for-the-badge)
![Features](https://img.shields.io/badge/Features-15+-blue?style=for-the-badge)

</div>

---

## 📄 License

This project is developed as part of an academic assignment.

**Version**: 1.0.0 🚀  
**Status**: Production Ready ✅  
**Maintained**: Yes 🔄

---

<div align="center">

**⭐ If you found this project helpful, please give it a star!**

[🔝 Back to Top](#-go-directory---contact-management-system)

</div>
