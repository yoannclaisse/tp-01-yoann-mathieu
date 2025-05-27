package server

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"tp1/annuaire"
)

// Global directory instance for managing contacts across all HTTP handlers
// This singleton pattern allows all web requests to operate on the same contact data
var dir *annuaire.Directory

// Custom template functions for HTML rendering and data manipulation
// These functions extend the default Go template functionality for better UI presentation
var templateFuncs = template.FuncMap{
	// substr extracts a substring and converts it to uppercase for avatar initials
	"substr": func(s string, start, length int) string {
		if start >= len(s) {
			return ""
		}
		end := start + length
		if end > len(s) {
			end = len(s)
		}
		return strings.ToUpper(s[start:end])
	},
	// eq provides equality comparison for template conditionals
	"eq": func(a, b interface{}) bool {
		return a == b
	},
}

// HTML template for the web interface
const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Directory - Web Interface</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            backdrop-filter: blur(10px);
            overflow: hidden;
        }

        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
            position: relative;
        }

        .header::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 0;
            right: 0;
            height: 4px;
            background: linear-gradient(90deg, #ff6b6b, #4ecdc4, #45b7d1);
        }

        .header h1 {
            font-size: 2.5rem;
            font-weight: 300;
            margin-bottom: 10px;
        }

        .header .subtitle {
            font-size: 1.1rem;
            opacity: 0.9;
        }

        .stats-card {
            background: linear-gradient(135deg, #ff6b6b 0%, #ee5a52 100%);
            color: white;
            margin: 20px;
            padding: 20px;
            border-radius: 15px;
            text-align: center;
            box-shadow: 0 10px 30px rgba(255, 107, 107, 0.3);
        }

        .stats-card i {
            font-size: 2rem;
            margin-bottom: 10px;
        }

        .stats-number {
            font-size: 2.5rem;
            font-weight: bold;
            margin: 10px 0;
        }

        .main-content {
            padding: 30px;
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 30px;
        }

        .section-card {
            background: white;
            border-radius: 15px;
            padding: 25px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
            border: 1px solid rgba(0, 0, 0, 0.05);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }

        .section-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.12);
        }

        .section-title {
            display: flex;
            align-items: center;
            font-size: 1.4rem;
            font-weight: 600;
            color: #333;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #f0f0f0;
        }

        .section-title i {
            margin-right: 10px;
            color: #667eea;
        }

        .form-group {
            margin-bottom: 20px;
        }

        .input-group {
            position: relative;
            margin-bottom: 15px;
        }

        .input-group i {
            position: absolute;
            left: 15px;
            top: 50%;
            transform: translateY(-50%);
            color: #999;
        }

        input[type="text"], input[type="file"] {
            width: 100%;
            padding: 15px 15px 15px 45px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 1rem;
            transition: border-color 0.3s ease, box-shadow 0.3s ease;
        }

        input[type="text"]:focus, input[type="file"]:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }

        .btn {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 12px 25px;
            border-radius: 10px;
            font-size: 1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            display: inline-flex;
            align-items: center;
            gap: 8px;
            text-decoration: none;
        }

        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(102, 126, 234, 0.3);
        }

        .btn-success {
            background: linear-gradient(135deg, #4ecdc4 0%, #44a08d 100%);
        }

        .btn-success:hover {
            box-shadow: 0 10px 25px rgba(78, 205, 196, 0.3);
        }

        .btn-danger {
            background: linear-gradient(135deg, #ff6b6b 0%, #ee5a52 100%);
        }

        .btn-danger:hover {
            box-shadow: 0 10px 25px rgba(255, 107, 107, 0.3);
        }

        .btn-small {
            padding: 8px 15px;
            font-size: 0.9rem;
        }

        .message {
            padding: 15px 20px;
            border-radius: 10px;
            margin: 20px;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .message.success {
            background: linear-gradient(135deg, #d4edda 0%, #c3e6cb 100%);
            color: #155724;
            border-left: 4px solid #28a745;
        }

        .message.error {
            background: linear-gradient(135deg, #f8d7da 0%, #f5c6cb 100%);
            color: #721c24;
            border-left: 4px solid #dc3545;
        }

        .contacts-grid {
            grid-column: 1 / -1;
            margin-top: 20px;
        }

        .contact-card {
            background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
            border-radius: 10px;
            padding: 20px;
            margin-bottom: 15px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            transition: all 0.3s ease;
            border-left: 4px solid #667eea;
        }

        .contact-card:hover {
            transform: translateX(5px);
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);
        }

        .contact-info {
            display: flex;
            align-items: center;
            gap: 15px;
        }

        .contact-avatar {
            width: 50px;
            height: 50px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            font-size: 1.2rem;
        }

        .contact-details h3 {
            color: #333;
            margin-bottom: 5px;
        }

        .contact-details p {
            color: #666;
            display: flex;
            align-items: center;
            gap: 5px;
        }

        .search-result {
            background: linear-gradient(135deg, #fff3cd 0%, #ffeaa7 100%);
            border: 2px solid #ffc107;
            border-radius: 10px;
            padding: 20px;
            margin: 20px 0;
        }

        .search-results {
            background: linear-gradient(135deg, #fff3cd 0%, #ffeaa7 100%);
            border: 2px solid #ffc107;
            border-radius: 10px;
            padding: 20px;
            margin: 20px 0;
        }

        .search-results h3 {
            margin-bottom: 15px;
            color: #856404;
        }

        .file-management {
            grid-column: 1 / -1;
            background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
            border-radius: 15px;
            padding: 25px;
            margin-top: 20px;
        }

        .file-actions {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }

        .file-card {
            background: white;
            border-radius: 10px;
            padding: 20px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.08);
        }

        .no-contacts {
            text-align: center;
            padding: 40px;
            color: #666;
            font-size: 1.1rem;
        }

        .no-contacts i {
            font-size: 4rem;
            color: #ddd;
            margin-bottom: 20px;
        }

        @media (max-width: 768px) {
            .main-content {
                grid-template-columns: 1fr;
                gap: 20px;
                padding: 20px;
            }
            
            .header h1 {
                font-size: 2rem;
            }
            
            .contact-card {
                flex-direction: column;
                align-items: flex-start;
                gap: 15px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1><i class="fas fa-address-book"></i> Go Directory</h1>
            <p class="subtitle">Modern Web Interface - Local Memory Management</p>
        </div>
        
        <div class="stats-card">
            <i class="fas fa-users"></i>
            <div class="stats-number">{{.ContactCount}}</div>
            <div>Contacts in memory</div>
        </div>

        {{if .Message}}
            <div class="message {{.MessageType}}">
                {{if eq .MessageType "success"}}
                    <i class="fas fa-check-circle"></i>
                {{else}}
                    <i class="fas fa-exclamation-triangle"></i>
                {{end}}
                <span>{{.Message}}</span>
            </div>
        {{end}}

        <div class="main-content">
            <div class="section-card">
                <h2 class="section-title">
                    <i class="fas fa-user-plus"></i>
                    Add Contact
                </h2>
                <form action="/add" method="POST">
                    <div class="input-group">
                        <i class="fas fa-user"></i>
                        <input type="text" name="name" placeholder="Last Name" required>
                    </div>
                    <div class="input-group">
                        <i class="fas fa-user"></i>
                        <input type="text" name="first" placeholder="First Name" required>
                    </div>
                    <div class="input-group">
                        <i class="fas fa-phone"></i>
                        <input type="text" name="phone" placeholder="Phone Number" required>
                    </div>
                    <button type="submit" class="btn">
                        <i class="fas fa-plus"></i>
                        Add Contact
                    </button>
                </form>
            </div>

            <div class="section-card">
                <h2 class="section-title">
                    <i class="fas fa-search"></i>
                    Search Contact
                </h2>
                <form action="/search" method="GET">
                    <div class="input-group">
                        <i class="fas fa-search"></i>
                        <input type="text" name="name" placeholder="Search by name, first name, or phone number" required>
                    </div>
                    <button type="submit" class="btn">
                        <i class="fas fa-search"></i>
                        Search
                    </button>
                </form>
            </div>
        </div>

        {{if .SearchResults}}
        <div class="search-results">
            <h3><i class="fas fa-user-check"></i> Search Results ({{len .SearchResults}} found)</h3>
            {{range .SearchResults}}
            <div class="contact-card" style="margin-top: 15px;">
                <div class="contact-info">
                    <div class="contact-avatar">
                        {{substr .First 0 1}}{{substr .Name 0 1}}
                    </div>
                    <div class="contact-details">
                        <h3>{{.First}} {{.Name}}</h3>
                        <p><i class="fas fa-phone"></i> {{.Phone}}</p>
                    </div>
                </div>
                <form action="/delete" method="POST">
                    <input type="hidden" name="name" value="{{.Name}}">
                    <button type="submit" class="btn btn-danger btn-small" onclick="return confirm('Are you sure you want to delete this contact?')">
                        <i class="fas fa-trash"></i>
                        Delete
                    </button>
                </form>
            </div>
            {{end}}
        </div>
        {{end}}

        <div class="contacts-grid">
            <div class="section-card">
                <h2 class="section-title">
                    <i class="fas fa-list"></i>
                    Contact List
                </h2>
                {{if .Contacts}}
                    {{range .Contacts}}
                    <div class="contact-card">
                        <div class="contact-info">
                            <div class="contact-avatar">
                                {{substr .First 0 1}}{{substr .Name 0 1}}
                            </div>
                            <div class="contact-details">
                                <h3>{{.First}} {{.Name}}</h3>
                                <p><i class="fas fa-phone"></i> {{.Phone}}</p>
                            </div>
                        </div>
                        <form action="/delete" method="POST">
                            <input type="hidden" name="name" value="{{.Name}}">
                            <button type="submit" class="btn btn-danger btn-small" onclick="return confirm('Are you sure you want to delete this contact?')">
                                <i class="fas fa-trash"></i>
                                Delete
                            </button>
                        </form>
                    </div>
                    {{end}}
                {{else}}
                    <div class="no-contacts">
                        <i class="fas fa-address-book"></i>
                        <p>No contacts in directory</p>
                        <p style="font-size: 0.9rem; margin-top: 10px;">Start by adding your first contact!</p>
                    </div>
                {{end}}
            </div>
        </div>

        <div class="file-management">
            <h2 class="section-title">
                <i class="fas fa-file-archive"></i>
                File Management
            </h2>
            
            <div class="file-actions">
                <div class="file-card">
                    <h3><i class="fas fa-download"></i> Export Contacts</h3>
                    <form action="/export" method="POST" style="margin-top: 15px;">
                        <div class="input-group">
                            <i class="fas fa-file-export"></i>
                            <input type="text" name="filename" placeholder="File name" value="contacts_export.json" required>
                        </div>
                        <button type="submit" class="btn btn-success">
                            <i class="fas fa-download"></i>
                            Prepare Download
                        </button>
                    </form>
                </div>
                
                <div class="file-card">
                    <h3><i class="fas fa-upload"></i> Import Contacts</h3>
                    <form action="/import" method="POST" enctype="multipart/form-data" style="margin-top: 15px;">
                        <div class="input-group">
                            <input type="file" name="file" accept=".json" required style="padding-left: 15px;">
                        </div>
                        <button type="submit" class="btn btn-success">
                            <i class="fas fa-upload"></i>
                            Import File
                        </button>
                    </form>
                </div>
                
                <div class="file-card">
                    <h3><i class="fas fa-broom"></i> Clear Memory</h3>
                    <p style="color: #666; margin: 15px 0;">Delete all contacts from local memory</p>
                    <form action="/clear" method="POST">
                        <button type="submit" class="btn btn-danger" onclick="return confirm('Are you sure you want to clear local memory?')">
                            <i class="fas fa-trash-alt"></i>
                            Clear Memory
                        </button>
                    </form>
                </div>
            </div>
        </div>
    </div>

    <script>
        // Add some basic interactivity
        document.addEventListener('DOMContentLoaded', function() {
            // Auto-hide messages after 5 seconds
            const messages = document.querySelectorAll('.message');
            messages.forEach(message => {
                setTimeout(() => {
                    message.style.opacity = '0';
                    message.style.transform = 'translateY(-20px)';
                    setTimeout(() => {
                        message.style.display = 'none';
                    }, 300);
                }, 5000);
            });
        });
    </script>
</body>
</html>
`

/**
 * PageData represents the data structure passed to HTML templates
 *
 * This structure encapsulates all data needed by the web interface templates
 * including contact lists, search results, messages, and statistics
 */
type PageData struct {
	Contacts      []annuaire.Contact // Complete list of all contacts for main display
	SearchResult  *annuaire.Contact  // Single search result (maintained for backward compatibility)
	SearchResults []annuaire.Contact // Multiple search results for enhanced search functionality
	Message       string             // Status message to display to user (success/error/info)
	MessageType   string             // CSS class type for message styling (success/error)
	ContactCount  int                // Total number of contacts for statistics display
}

/**
 * createTemplate creates an HTML template with custom functions
 *
 * @return {*template.Template} Parsed template ready for execution
 * @return {error} Error if template parsing fails
 *
 * This function combines the HTML template string with custom template functions
 * to create a fully functional template for web page rendering
 */
func createTemplate() (*template.Template, error) {
	return template.New("home").Funcs(templateFuncs).Parse(htmlTemplate)
}

/**
 * StartServer initializes and starts the HTTP web server on port 8080
 *
 * This function sets up the web application by:
 * - Initializing an empty contact directory (no automatic file loading)
 * - Registering all HTTP route handlers for web interface functionality
 * - Starting the HTTP server and listening for incoming connections
 *
 * The server will panic if it fails to bind to port 8080 or encounters
 * other critical startup errors
 */
func StartServer() {
	// Initialize empty directory (no automatic loading for web interface)
	// This gives users a clean slate and explicit control over data loading
	dir = annuaire.NewDirectory()

	// Register HTTP route handlers for all web interface functionality
	http.HandleFunc("/", handleHome)              // Main page with contact list and forms
	http.HandleFunc("/add", handleAdd)            // POST: Add new contact
	http.HandleFunc("/search", handleSearch)      // GET: Search for contacts
	http.HandleFunc("/delete", handleDelete)      // POST: Delete contact
	http.HandleFunc("/export", handleExport)      // POST: Export contacts to JSON
	http.HandleFunc("/import", handleImport)      // POST: Import contacts from JSON
	http.HandleFunc("/clear", handleClear)        // POST: Clear all contacts from memory
	http.HandleFunc("/download/", handleDownload) // GET: Download exported files

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/**
 * handleHome renders and serves the main page with contact list
 *
 * @param {http.ResponseWriter} w - HTTP response writer for sending HTML content
 * @param {*http.Request} r - HTTP request containing URL parameters and form data
 *
 * This handler processes the root route ("/") and displays:
 * - Complete list of all contacts in the directory
 * - Contact statistics (total count)
 * - Success/error messages from redirected operations
 * - All interactive forms for contact management
 */
func handleHome(w http.ResponseWriter, r *http.Request) {
	// Create template instance with custom functions
	tmpl, err := createTemplate()
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Prepare data structure for template rendering
	data := PageData{
		Contacts:     dir.ListContacts(), // Get all contacts for main display
		ContactCount: dir.ContactCount(), // Get statistics for header display
	}

	// Check for messages in URL parameters (from redirected operations)
	if msg := r.URL.Query().Get("message"); msg != "" {
		data.Message = msg
		data.MessageType = r.URL.Query().Get("type")
		// Default to success message type if not specified
		if data.MessageType == "" {
			data.MessageType = "success"
		}
	}

	// Execute template with prepared data and send to client
	tmpl.Execute(w, data)
}

/**
 * handleAdd processes POST requests to add new contacts
 *
 * @param {http.ResponseWriter} w - HTTP response writer for redirect responses
 * @param {*http.Request} r - HTTP request containing form data
 *
 * This handler:
 * - Validates HTTP method (POST only)
 * - Extracts contact data from form fields
 * - Attempts to add contact to directory
 * - Redirects back to home page with success/error message
 */
func handleAdd(w http.ResponseWriter, r *http.Request) {
	// Enforce POST method for data modification operations
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Extract contact information from form data
	name := r.FormValue("name")   // Last name from form
	first := r.FormValue("first") // First name from form
	phone := r.FormValue("phone") // Phone number from form

	// Attempt to add contact to directory with validation
	err := dir.AddContact(name, first, phone)

	// Prepare redirect URL with appropriate success/error message
	redirectURL := "/"
	if err != nil {
		// Format error message for user display
		message := fmt.Sprintf("Error: %v", err)
		redirectURL = fmt.Sprintf("/?message=%s&type=error", url.QueryEscape(message))
	} else {
		// Format success message with contact details
		message := fmt.Sprintf("Contact %s %s added successfully to local memory", first, name)
		redirectURL = fmt.Sprintf("/?message=%s&type=success", url.QueryEscape(message))
	}

	// Redirect back to home page to display result
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

/**
 * handleSearch processes search requests and displays results
 *
 * @param {http.ResponseWriter} w - HTTP response writer for HTML content
 * @param {*http.Request} r - HTTP request containing search parameters
 *
 * This handler provides comprehensive search functionality:
 * - Accepts search terms from query parameters
 * - Uses FilterContacts to find all matching contacts
 * - Displays search results alongside the main contact list
 * - Provides detailed debug output for troubleshooting search issues
 */
func handleSearch(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.FormValue("name")

	// DEBUG: Print comprehensive search debugging information
	// This debug block helps developers troubleshoot search functionality issues
	fmt.Printf("=== SEARCH DEBUG START ===\n")
	fmt.Printf("Search term received: '%s'\n", searchTerm)
	fmt.Printf("Total contacts in directory: %d\n", dir.ContactCount())

	// DEBUG: Display all contacts currently in the directory for verification
	// This helps identify data issues or contact storage problems
	contacts := dir.ListContacts()
	fmt.Printf("--- All Contacts in Directory ---\n")
	for i, contact := range contacts {
		fmt.Printf("Contact %d: '%s' '%s' - '%s'\n", i+1, contact.First, contact.Name, contact.Phone)
	}
	fmt.Printf("--- End Contact List ---\n")

	// Create template for rendering search results
	tmpl, _ := createTemplate()
	data := PageData{
		Contacts:     contacts,           // Show all contacts alongside search results
		ContactCount: dir.ContactCount(), // Display current statistics
	}

	// Process search request if search term is provided
	if searchTerm != "" {
		// DEBUG: Log the start of search processing
		fmt.Printf("Processing search for term: '%s'\n", searchTerm)

		// Use FilterContacts to get all matching contacts (not just first match)
		searchResults := dir.FilterContacts(searchTerm)

		// DEBUG: Report search results for verification
		fmt.Printf("Search completed. Found %d results:\n", len(searchResults))
		for i, result := range searchResults {
			fmt.Printf("  Result %d: %s %s - %s\n", i+1, result.First, result.Name, result.Phone)
		}

		if len(searchResults) > 0 {
			// Store search results for template display
			data.SearchResults = searchResults
			// Maintain backward compatibility by setting first result as SearchResult
			data.SearchResult = &searchResults[0]

			// Set appropriate success message based on result count
			if len(searchResults) == 1 {
				data.Message = "Contact found"
			} else {
				data.Message = fmt.Sprintf("%d contacts found", len(searchResults))
			}
			data.MessageType = "success"

			// DEBUG: Log template data being prepared
			fmt.Printf("Template data prepared:\n")
			fmt.Printf("  SearchResults count: %d\n", len(searchResults))
			fmt.Printf("  Message: '%s'\n", data.Message)
			fmt.Printf("  MessageType: '%s'\n", data.MessageType)
		} else {
			// No results found - prepare error message
			data.Message = fmt.Sprintf("No contact found matching: %s", searchTerm)
			data.MessageType = "error"

			// DEBUG: Log no-match scenario for troubleshooting
			fmt.Printf("No matches found for search term: '%s'\n", searchTerm)
			fmt.Printf("This could indicate:\n")
			fmt.Printf("  - Search term doesn't match any contact exactly\n")
			fmt.Printf("  - Case sensitivity issues\n")
			fmt.Printf("  - Contact data structure problems\n")
		}
	}

	// DEBUG: Final debug output before template execution
	fmt.Printf("=== SEARCH DEBUG END ===\n\n")

	// Execute template with search results and contact data
	if err := tmpl.Execute(w, data); err != nil {
		// DEBUG: Log template execution errors for debugging
		fmt.Printf("TEMPLATE EXECUTION ERROR: %v\n", err)
		fmt.Printf("Data structure passed to template: %+v\n", data)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

/**
 * handleDelete processes POST requests to delete contacts
 *
 * @param {http.ResponseWriter} w - HTTP response writer for redirect responses
 * @param {*http.Request} r - HTTP request containing contact name to delete
 *
 * This handler:
 * - Validates HTTP method (POST only)
 * - Extracts contact name from form data
 * - Attempts to delete contact from directory
 * - Redirects back to home page with success/error message
 */
func handleDelete(w http.ResponseWriter, r *http.Request) {
	// Enforce POST method for data modification operations
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Extract contact name to delete from form data
	name := r.FormValue("name")

	// Attempt to delete contact from directory
	err := dir.DeleteContact(name)

	// Prepare redirect URL with appropriate success/error message
	redirectURL := "/"
	if err != nil {
		// Format error message for user display
		message := fmt.Sprintf("Error: %v", err)
		redirectURL = fmt.Sprintf("/?message=%s&type=error", url.QueryEscape(message))
	} else {
		// Format success message with deleted contact name
		message := fmt.Sprintf("Contact %s deleted successfully from local memory", name)
		redirectURL = fmt.Sprintf("/?message=%s&type=success", url.QueryEscape(message))
	}

	// Redirect back to home page to display result
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

/**
 * handleExport prepares contact data for download as JSON file
 *
 * This handler:
 * - Validates HTTP method (POST only)
 * - Extracts or defaults the filename for export
 * - Creates a temporary directory for export files
 * - Exports the contact directory to a JSON file
 * - Redirects with a download link or error message
 */
func handleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	filename := r.FormValue("filename")
	if filename == "" {
		filename = "contacts_export.json"
	}

	// Create temp directory if it doesn't exist
	tempDir := "temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		message := "Error creating temporary directory"
		redirectURL := fmt.Sprintf("/?message=%s&type=error", url.QueryEscape(message))
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Full path of temporary file
	tempFile := filepath.Join(tempDir, filename)

	err := dir.ExportToJSON(tempFile)

	// Prepare redirect URL with message
	redirectURL := "/"
	if err != nil {
		message := fmt.Sprintf("Export error: %v", err)
		redirectURL = fmt.Sprintf("/?message=%s&type=error", url.QueryEscape(message))
	} else {
		downloadURL := fmt.Sprintf("/download/%s", filename)
		message := fmt.Sprintf(`Export successful! <a href="%s" class="download-btn">Download %s</a>`, downloadURL, filename)
		redirectURL = fmt.Sprintf("/?message=%s&type=success", url.QueryEscape(message))
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// handleDownload serves exported files for download
// Automatically deletes temporary files after serving
func handleDownload(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL
	filename := r.URL.Path[len("/download/"):]

	// Full file path
	filepath := filepath.Join("temp", filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set download headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/json")

	// Copy file content to response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Download error", http.StatusInternalServerError)
		return
	}

	// Delete temporary file after download
	go func() {
		os.Remove(filepath)
	}()
}

/**
 * handleImport processes uploaded JSON files and imports contact data
 *
 * This handler:
 * - Validates HTTP method (POST only)
 * - Parses the multipart form data containing the file
 * - Creates a temporary file for the uploaded content
 * - Imports contact data from the JSON file into the directory
 * - Redirects with success/error message
 */
func handleImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		message := fmt.Sprintf("Form parsing error: %v", err)
		redirectURL := fmt.Sprintf("/?message=%s&type=error", url.QueryEscape(message))
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		message := fmt.Sprintf("File retrieval error: %v", err)
		redirectURL := fmt.Sprintf("/?message=%s&type=error", url.QueryEscape(message))
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	defer file.Close()

	// Create temporary file
	tempDir := "temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		message := "Error creating temporary directory"
		redirectURL := fmt.Sprintf("/?message=%s&type=error", url.QueryEscape(message))
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	tempFile := filepath.Join(tempDir, "import_"+header.Filename)
	dst, err := os.Create(tempFile)
	if err != nil {
		message := fmt.Sprintf("Temporary file creation error: %v", err)
		redirectURL := fmt.Sprintf("/?message=%s&type=error", url.QueryEscape(message))
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	defer dst.Close()
	defer os.Remove(tempFile) // Clean up temporary file

	// Copy uploaded file content
	_, err = io.Copy(dst, file)
	if err != nil {
		message := fmt.Sprintf("File copy error: %v", err)
		redirectURL := fmt.Sprintf("/?message=%s&type=error", url.QueryEscape(message))
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Close file before importing
	dst.Close()

	// Import data
	err = dir.ImportFromJSON(tempFile)

	// Prepare redirect URL with message
	redirectURL := "/"
	if err != nil {
		message := fmt.Sprintf("Import error from %s: %v", header.Filename, err)
		redirectURL = fmt.Sprintf("/?message=%s&type=error", url.QueryEscape(message))
	} else {
		message := fmt.Sprintf("Data imported successfully from %s (%d contacts loaded)", header.Filename, dir.ContactCount())
		redirectURL = fmt.Sprintf("/?message=%s&type=success", url.QueryEscape(message))
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

/**
 * handleClear removes all contacts from local memory
 *
 * @param {http.ResponseWriter} w - HTTP response writer for redirect responses
 * @param {*http.Request} r - HTTP request (POST method required)
 *
 * This handler provides a complete reset functionality by:
 * - Creating a new empty directory instance
 * - Replacing the global directory variable
 * - Redirecting with success confirmation message
 *
 * Note: This operation only affects the in-memory data, not any saved files
 */
func handleClear(w http.ResponseWriter, r *http.Request) {
	// Enforce POST method for data modification operations
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Replace global directory with new empty instance
	// This effectively clears all contacts from memory
	dir = annuaire.NewDirectory()

	// Prepare success message and redirect to home page
	message := "Local memory cleared successfully"
	redirectURL := fmt.Sprintf("/?message=%s&type=success", url.QueryEscape(message))
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
