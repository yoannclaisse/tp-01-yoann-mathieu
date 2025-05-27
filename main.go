package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"tp1/annuaire"
	"tp1/server"
)

// Default data file path for persistent contact storage
// This file serves as the primary storage location for CLI operations
const defaultDataFile = "data/contacts.json"

/**
 * main is the entry point of the application
 *
 * This function serves as the command-line interface dispatcher, handling:
 * - Command-line argument parsing and validation
 * - Routing to appropriate action handlers (CLI operations vs web server)
 * - Global error handling and application exit codes
 * - Data directory initialization and file management
 *
 * The application supports two primary modes:
 * 1. CLI mode: Direct command execution with immediate results
 * 2. Web server mode: HTTP server providing browser-based interface
 */
func main() {
	// Define command-line flags with comprehensive help descriptions
	var action = flag.String("action", "", "Action to perform (add, list, search, delete, update, export, import)")
	var name = flag.String("name", "", "Contact last name")
	var first = flag.String("first", "", "Contact first name")
	var phone = flag.String("phone", "", "Phone number")
	var file = flag.String("file", "", "JSON file for import/export (required for export/import)")
	var webserver = flag.Bool("server", false, "Start web server")

	// Parse all command-line arguments
	flag.Parse()

	// Check for web server mode and start HTTP server if requested
	if *webserver {
		server.StartServer() // This call blocks until server shutdown
		return
	}

	// Initialize data storage directory structure
	// Create the data directory if it doesn't exist to ensure file operations succeed
	if err := os.MkdirAll(filepath.Dir(defaultDataFile), 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize directory instance for CLI operations
	dir := annuaire.NewDirectory()

	// Load existing contacts from persistent storage
	// This provides continuity between CLI sessions
	if err := dir.ImportFromJSON(defaultDataFile); err != nil {
		// Only show warning for actual errors, not missing files
		if !os.IsNotExist(err) {
			fmt.Printf("Warning: Error loading contacts: %v\n", err)
		}
		// Continue execution with empty directory if file doesn't exist
	}

	// Route to appropriate action handler based on command-line arguments
	switch *action {
	case "add":
		handleAddAction(dir, *name, *first, *phone)
	case "list":
		handleListAction(dir)
	case "search":
		handleSearchAction(dir, *name)
	case "delete":
		handleDeleteAction(dir, *name)
	case "update":
		handleUpdateAction(dir, *name, *first, *phone)
	case "export":
		handleExportAction(dir, *file)
	case "import":
		handleImportAction(dir, *file)
	case "":
		// No action specified - show usage information
		printUsage()
	default:
		// Unknown action specified
		fmt.Printf("Action '%s' not implemented\n", *action)
		os.Exit(1)
	}
}

/**
 * handleAddAction processes the add contact command
 *
 * @param {*annuaire.Directory} dir - Directory instance to add contact to
 * @param {string} name - Last name of the contact
 * @param {string} first - First name of the contact
 * @param {string} phone - Phone number of the contact
 *
 * This function performs comprehensive validation and provides user feedback:
 * - Validates that all required fields are provided
 * - Attempts to add contact with error handling
 * - Automatically saves changes to persistent storage
 * - Provides success confirmation or error messages
 */
func handleAddAction(dir *annuaire.Directory, name, first, phone string) {
	// Validate that all required fields are provided
	if name == "" || first == "" || phone == "" {
		fmt.Println("Error: name, first name and phone required")
		os.Exit(1)
	}

	// Attempt to add contact to directory
	err := dir.AddContact(name, first, phone)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Save changes to persistent storage to maintain data between sessions
	if err := dir.ExportToJSON(defaultDataFile); err != nil {
		fmt.Printf("Warning: Error saving: %v\n", err)
	}

	// Confirm successful addition to user
	fmt.Printf("Contact %s %s added successfully\n", first, name)
}

/**
 * handleListAction processes the list contacts command
 *
 * @param {*annuaire.Directory} dir - Directory instance to list contacts from
 *
 * This function provides formatted output of all contacts:
 * - Handles empty directory case with user-friendly message
 * - Shows contact count statistics
 * - Formats contact information consistently
 */
func handleListAction(dir *annuaire.Directory) {
	contacts := dir.ListContacts()

	// Handle empty directory case
	if len(contacts) == 0 {
		fmt.Println("No contacts found")
	} else {
		// Display contact count and formatted list
		fmt.Printf("Contact list (%d total):\n", len(contacts))
		for _, contact := range contacts {
			fmt.Printf("- %s %s: %s\n", contact.First, contact.Name, contact.Phone)
		}
	}
}

/**
 * handleSearchAction processes the search contact command
 *
 * @param {*annuaire.Directory} dir - Directory instance to search
 * @param {string} searchTerm - Term to search for
 *
 * This function provides single-result search functionality:
 * - Validates that search term is provided
 * - Searches across name, first name, and phone fields
 * - Provides clear feedback for found/not found cases
 */
func handleSearchAction(dir *annuaire.Directory, searchTerm string) {
	// Validate that search term is provided
	if searchTerm == "" {
		fmt.Println("Error: search term required")
		os.Exit(1)
	}

	// Perform search operation
	contact, exists := dir.SearchContact(searchTerm)
	if exists {
		// Display found contact information
		fmt.Printf("Contact found: %s %s - %s\n", contact.First, contact.Name, contact.Phone)
	} else {
		// Inform user that no match was found
		fmt.Printf("No contact found matching: %s\n", searchTerm)
	}
}

/**
 * handleDeleteAction processes the delete contact command
 *
 * @param {*annuaire.Directory} dir - Directory instance to delete from
 * @param {string} name - Last name of contact to delete
 *
 * This function provides safe deletion with persistence:
 * - Validates that contact name is provided
 * - Attempts deletion with error handling
 * - Automatically saves changes to persistent storage
 * - Provides success confirmation or error messages
 */
func handleDeleteAction(dir *annuaire.Directory, name string) {
	// Validate that contact name is provided
	if name == "" {
		fmt.Println("Error: name required")
		os.Exit(1)
	}

	// Attempt to delete contact
	err := dir.DeleteContact(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Save changes to persistent storage
	if err := dir.ExportToJSON(defaultDataFile); err != nil {
		fmt.Printf("Warning: Error saving: %v\n", err)
	}

	// Confirm successful deletion
	fmt.Printf("Contact %s deleted successfully\n", name)
}

/**
 * handleUpdateAction processes the update contact command
 *
 * @param {*annuaire.Directory} dir - Directory instance to update
 * @param {string} name - Last name of contact to update (required)
 * @param {string} first - New first name (optional)
 * @param {string} phone - New phone number (optional)
 *
 * This function provides flexible update functionality:
 * - Validates that contact name is provided (required for lookup)
 * - Allows partial updates (empty fields are not changed)
 * - Automatically saves changes to persistent storage
 * - Provides success confirmation or error messages
 */
func handleUpdateAction(dir *annuaire.Directory, name, first, phone string) {
	// Validate that contact name is provided for lookup
	if name == "" {
		fmt.Println("Error: name required")
		os.Exit(1)
	}

	// Attempt to update contact (empty fields will be ignored)
	err := dir.UpdateContact(name, first, phone)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Save changes to persistent storage
	if err := dir.ExportToJSON(defaultDataFile); err != nil {
		fmt.Printf("Warning: Error saving: %v\n", err)
	}

	// Confirm successful update
	fmt.Printf("Contact %s updated successfully\n", name)
}

/**
 * handleExportAction processes the export contacts command
 *
 * @param {*annuaire.Directory} dir - Directory instance to export from
 * @param {string} file - Target file path for export
 *
 * This function provides data backup and sharing functionality:
 * - Validates that file path is provided
 * - Exports all contacts to specified JSON file
 * - Provides success confirmation or error messages
 */
func handleExportAction(dir *annuaire.Directory, file string) {
	// Validate that file path is provided
	if file == "" {
		fmt.Println("Error: file path required for export (-file)")
		os.Exit(1)
	}

	// Attempt to export contacts to specified file
	err := dir.ExportToJSON(file)
	if err != nil {
		fmt.Printf("Export error: %v\n", err)
		os.Exit(1)
	}

	// Confirm successful export
	fmt.Printf("Contacts exported to %s\n", file)
}

/**
 * handleImportAction processes the import contacts command
 *
 * @param {*annuaire.Directory} dir - Directory instance to import into
 * @param {string} file - Source file path for import
 *
 * This function provides data restoration and sharing functionality:
 * - Validates that file path is provided
 * - Imports contacts from specified JSON file
 * - Automatically saves imported data to default storage
 * - Provides success confirmation or error messages
 */
func handleImportAction(dir *annuaire.Directory, file string) {
	// Validate that file path is provided
	if file == "" {
		fmt.Println("Error: file path required for import (-file)")
		os.Exit(1)
	}

	// Attempt to import contacts from specified file
	err := dir.ImportFromJSON(file)
	if err != nil {
		fmt.Printf("Import error: %v\n", err)
		os.Exit(1)
	}

	// Save imported data to default storage location for future CLI sessions
	if err := dir.ExportToJSON(defaultDataFile); err != nil {
		fmt.Printf("Warning: Error saving: %v\n", err)
	}

	// Confirm successful import
	fmt.Printf("Contacts imported from %s\n", file)
}

/**
 * printUsage displays available commands and usage information
 *
 * This function provides comprehensive help information including:
 * - List of all available actions with descriptions
 * - Required and optional parameters for each action
 * - Information about persistent storage location
 * - Command-line flag documentation
 */
func printUsage() {
	fmt.Println("📞 Go Directory - Contact Management System")
	fmt.Println("===========================================")
	fmt.Println()
	fmt.Println("Available actions:")
	fmt.Println("  add      - Add a contact (name, first, phone required)")
	fmt.Println("  list     - List all contacts")
	fmt.Println("  search   - Search for a contact by name, first name, or phone (name required)")
	fmt.Println("  delete   - Delete a contact (name required)")
	fmt.Println("  update   - Update a contact (name required)")
	fmt.Println("  export   - Export to JSON file (file required)")
	fmt.Println("  import   - Import from JSON file (file required)")
	fmt.Println("  server   - Start web interface")
	fmt.Println()
	fmt.Printf("📁 Contacts are automatically saved to: %s\n", defaultDataFile)
	fmt.Println()
	fmt.Println("Command-line flags:")
	flag.PrintDefaults()
}
