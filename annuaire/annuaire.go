package annuaire

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Contact represents a single contact entry in the directory
// This structure defines the core data model for storing individual contact information
// Each contact contains a last name, first name, and phone number
type Contact struct {
	Name  string `json:"name"`  // Last name of the contact (required, used as primary identifier)
	First string `json:"first"` // First name of the contact (required)
	Phone string `json:"phone"` // Phone number of the contact (required, part of composite key)
}

// Directory manages a collection of contacts using an in-memory map
// The directory uses a composite key (name_phone) to allow multiple contacts
// with the same name but different phone numbers
// This design choice enables storing family members or business contacts with shared names
type Directory struct {
	contacts map[string]Contact // Internal storage using composite keys for uniqueness
}

/**
 * NewDirectory creates and returns a new empty directory instance
 *
 * @return {*Directory} A pointer to a newly initialized directory with an empty contact map
 *
 * Usage:
 *   dir := NewDirectory()
 *
 * The returned directory is ready to use with all CRUD operations available
 */
func NewDirectory() *Directory {
	return &Directory{
		contacts: make(map[string]Contact), // Initialize empty map for contact storage
	}
}

/**
 * AddContact adds a new contact to the directory with comprehensive validation
 *
 * @param {string} name - Last name of the contact (cannot be empty)
 * @param {string} first - First name of the contact (cannot be empty)
 * @param {string} phone - Phone number of the contact (cannot be empty)
 * @return {error} Returns an error if validation fails or contact already exists
 *
 * Validation rules:
 * - All fields must be non-empty strings
 * - Combination of name and phone must be unique (allows same name with different phones)
 *
 * Usage:
 *   err := dir.AddContact("Smith", "John", "555-1234")
 *   if err != nil {
 *       // Handle error (empty fields or duplicate name+phone combination)
 *   }
 */
func (d *Directory) AddContact(name, first, phone string) error {
	// Input validation - ensure all required fields are provided
	if name == "" || first == "" || phone == "" {
		return errors.New("all fields are required")
	}

	// Create composite key to allow multiple contacts with same name but different phones
	// This design enables storing contacts like "Smith, John (home)" and "Smith, John (work)"
	key := fmt.Sprintf("%s_%s", name, phone)

	// Check for duplicate entries using the composite key
	if _, exists := d.contacts[key]; exists {
		return errors.New("a contact with this name and phone already exists")
	}

	// Store the contact with the composite key for fast lookup
	d.contacts[key] = Contact{
		Name:  name,
		First: first,
		Phone: phone,
	}

	return nil
}

/**
 * SearchContact searches for and returns the first contact matching the search term
 *
 * @param {string} searchTerm - Term to search for (matches name, first name, or phone)
 * @return {Contact} The found contact (empty if not found)
 * @return {bool} True if contact was found, false otherwise
 *
 * Search behavior:
 * - Performs exact string matching (case-sensitive)
 * - Searches across name, first name, and phone fields
 * - Returns the first match found (order not guaranteed due to map iteration)
 *
 * Usage:
 *   contact, found := dir.SearchContact("Smith")
 *   if found {
 *       fmt.Printf("Found: %s %s", contact.First, contact.Name)
 *   }
 */
func (d *Directory) SearchContact(searchTerm string) (Contact, bool) {
	// DEBUG: Log search initiation for troubleshooting search operations
	log.Printf("SearchContact: Looking for '%s'", searchTerm)
	// DEBUG: Display total contacts to verify directory state during search
	log.Printf("Total contacts in directory: %d", len(d.contacts))

	// Iterate through all contacts to find exact matches
	for key, contact := range d.contacts {
		// DEBUG: Log each contact being checked to trace search execution path
		log.Printf("Checking contact: key='%s', name='%s', first='%s', phone='%s'",
			key, contact.Name, contact.First, contact.Phone)

		// Check if search term matches any of the contact's fields exactly
		if contact.Name == searchTerm || contact.First == searchTerm || contact.Phone == searchTerm {
			// DEBUG: Log successful match for debugging search results
			log.Printf("Found match: %+v", contact)
			return contact, true
		}
	}

	// DEBUG: Log when no match is found to help diagnose search issues
	log.Printf("No match found for '%s'", searchTerm)
	return Contact{}, false
}

/**
 * FilterContacts searches for and returns all contacts matching the search term
 *
 * @param {string} searchTerm - Term to search for (matches name, first name, or phone)
 * @return {[]Contact} Slice of all contacts that match the search criteria
 *
 * This method differs from SearchContact by returning ALL matches instead of just the first one
 * Useful for scenarios where multiple contacts might match (e.g., same last name)
 *
 * Usage:
 *   matches := dir.FilterContacts("Smith")
 *   fmt.Printf("Found %d contacts named Smith", len(matches))
 */
func (d *Directory) FilterContacts(searchTerm string) []Contact {
	// DEBUG: Log filter operation start for debugging multi-match scenarios
	log.Printf("FilterContacts: Looking for '%s'", searchTerm)
	// DEBUG: Show directory size to verify data state before filtering
	log.Printf("Total contacts in directory: %d", len(d.contacts))

	var matches []Contact

	// Scan all contacts for matches
	for key, contact := range d.contacts {
		// DEBUG: Trace each contact evaluation during filtering process
		log.Printf("Checking contact: key='%s', name='%s', first='%s', phone='%s'",
			key, contact.Name, contact.First, contact.Phone)

		// Apply same matching logic as SearchContact but collect all results
		if contact.Name == searchTerm || contact.First == searchTerm || contact.Phone == searchTerm {
			// DEBUG: Log each match found during filtering
			log.Printf("Found match: %+v", contact)
			matches = append(matches, contact)
		}
	}

	// DEBUG: Report final filter results for verification
	log.Printf("Found %d matches for '%s'", len(matches), searchTerm)
	return matches
}

/**
 * ListContacts returns a slice containing all contacts in the directory
 *
 * @return {[]Contact} Slice of all contacts (empty slice if no contacts exist)
 *
 * Note: Order of contacts is not guaranteed due to underlying map structure
 * For sorted output, the caller should sort the returned slice
 *
 * Usage:
 *   allContacts := dir.ListContacts()
 *   fmt.Printf("Total contacts: %d", len(allContacts))
 */
func (d *Directory) ListContacts() []Contact {
	// Pre-allocate slice with known capacity for better performance
	contacts := make([]Contact, 0, len(d.contacts))

	// Convert map values to slice for easier iteration by callers
	for _, contact := range d.contacts {
		contacts = append(contacts, contact)
	}
	return contacts
}

/**
 * DeleteContact removes the first contact with the specified name from the directory
 *
 * @param {string} name - Last name of the contact to delete
 * @return {error} Returns an error if no contact with the given name is found
 *
 * Deletion behavior:
 * - Searches by last name only (not first name or phone)
 * - Removes the first matching contact found
 * - If multiple contacts have the same last name, only one is deleted
 *
 * Usage:
 *   err := dir.DeleteContact("Smith")
 *   if err != nil {
 *       // Handle case where no contact named Smith exists
 *   }
 */
func (d *Directory) DeleteContact(name string) error {
	found := false

	// Search through all contacts to find the first match by last name
	for key, contact := range d.contacts {
		if contact.Name == name {
			// Remove the contact from the map using its composite key
			delete(d.contacts, key)
			found = true
			break // Exit after first match to maintain single-delete behavior
		}
	}

	// Return error if no matching contact was found
	if !found {
		return errors.New("contact not found")
	}
	return nil
}

/**
 * UpdateContact modifies an existing contact's first name and/or phone number
 *
 * @param {string} name - Last name of the contact to update (used for lookup)
 * @param {string} newFirst - New first name (empty string means no change)
 * @param {string} newPhone - New phone number (empty string means no change)
 * @return {error} Returns an error if no contact with the given name is found
 *
 * Update behavior:
 * - Searches by last name to find the contact
 * - Only updates fields that have non-empty values provided
 * - Preserves existing values for empty parameters
 * - Updates the first matching contact found
 *
 * Usage:
 *   // Update only phone number
 *   err := dir.UpdateContact("Smith", "", "555-9999")
 *
 *   // Update both first name and phone
 *   err := dir.UpdateContact("Smith", "Jane", "555-8888")
 */
func (d *Directory) UpdateContact(name, newFirst, newPhone string) error {
	// Search for the contact to update by last name
	for key, contact := range d.contacts {
		if contact.Name == name {
			// Update first name only if a new value is provided
			if newFirst != "" {
				contact.First = newFirst
			}
			// Update phone number only if a new value is provided
			if newPhone != "" {
				contact.Phone = newPhone
			}
			// Save the updated contact back to the map
			d.contacts[key] = contact
			return nil
		}
	}
	// Return error if no contact with the specified name exists
	return errors.New("contact not found")
}

/**
 * ContactCount returns the total number of contacts in the directory
 *
 * @return {int} The number of contacts currently stored
 *
 * This is a simple utility method that provides the size of the internal contacts map
 * Useful for statistics, validation, and UI display purposes
 *
 * Usage:
 *   count := dir.ContactCount()
 *   fmt.Printf("You have %d contacts", count)
 */
func (d *Directory) ContactCount() int {
	return len(d.contacts)
}

/**
 * ExportToJSON exports all contacts to a JSON file at the specified path
 *
 * @param {string} filename - Full path where the JSON file should be created
 * @return {error} Returns an error if file operations or JSON marshaling fails
 *
 * File operations:
 * - Creates directory structure if it doesn't exist
 * - Overwrites existing files without warning
 * - Uses proper JSON formatting with indentation for readability
 * - Converts internal map structure to array for standard JSON format
 *
 * Usage:
 *   err := dir.ExportToJSON("backup/contacts.json")
 *   if err != nil {
 *       // Handle file system or JSON encoding errors
 *   }
 */
func (d *Directory) ExportToJSON(filename string) error {
	// Create directory structure if it doesn't exist (recursive creation)
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Convert internal map to slice for proper JSON array structure
	// This ensures the JSON file contains a standard array format
	contacts := make([]Contact, 0, len(d.contacts))
	for _, contact := range d.contacts {
		contacts = append(contacts, contact)
	}

	// Marshal to JSON with indentation for human readability
	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return err
	}

	// Write JSON data to file with appropriate permissions
	return os.WriteFile(filename, data, 0644)
}

/**
 * ImportFromJSON imports contacts from a JSON file and replaces current data
 *
 * @param {string} filename - Path to the JSON file to import
 * @return {error} Returns an error if file doesn't exist or JSON parsing fails
 *
 * Import behavior:
 * - Completely replaces existing contacts (not additive)
 * - Expects JSON array format with Contact objects
 * - Reconstructs internal composite keys from imported data
 * - Validates JSON structure but not individual contact data
 *
 * Usage:
 *   err := dir.ImportFromJSON("contacts.json")
 *   if err != nil {
 *       // Handle file not found or malformed JSON errors
 *   }
 */
func (d *Directory) ImportFromJSON(filename string) error {
	// Check if file exists before attempting to read
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return errors.New("file not found")
	}

	// Read entire file content into memory
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Parse JSON array into slice of Contact structs
	var contacts []Contact
	if err := json.Unmarshal(data, &contacts); err != nil {
		return err
	}

	// Clear existing contacts and rebuild internal map structure
	d.contacts = make(map[string]Contact)
	for _, contact := range contacts {
		// Reconstruct composite key for internal storage
		key := fmt.Sprintf("%s_%s", contact.Name, contact.Phone)
		d.contacts[key] = contact
	}

	return nil
}

/**
 * DebugPrintContacts prints all contacts for debugging purposes
 *
 * This utility method outputs the complete internal state of the directory
 * including composite keys and contact data for troubleshooting
 *
 * Output format shows:
 * - Total number of contacts
 * - Each contact's internal key and all field values
 * - Clear visual separation for easy reading
 *
 * Usage:
 *   dir.DebugPrintContacts() // Call when debugging contact storage issues
 */
func (d *Directory) DebugPrintContacts() {
	fmt.Printf("=== DEBUG: Directory Contents ===\n")
	fmt.Printf("Total contacts: %d\n", len(d.contacts))

	// Display each contact with its internal storage key for debugging
	for key, contact := range d.contacts {
		fmt.Printf("Key: %s -> Name: %s, First: %s, Phone: %s\n",
			key, contact.Name, contact.First, contact.Phone)
	}
	fmt.Printf("================================\n")
}

// =============================================================================
// LEGACY COMPATIBILITY LAYER
// =============================================================================
// The following section provides backward compatibility for existing code
// that uses French method names. These methods are deprecated and should
// not be used in new code.

// Legacy type alias for backward compatibility with existing French code
type Annuaire = Directory

/**
 * NewAnnuaire creates a new directory instance (legacy function name)
 *
 * @deprecated Use NewDirectory instead for new code
 * @return {*Directory} A pointer to a newly initialized directory
 *
 * This function exists solely for backward compatibility with existing
 * French-named code and will be removed in future versions.
 */
func NewAnnuaire() *Directory {
	return NewDirectory()
}

// =============================================================================
// DEPRECATED FRENCH METHOD NAMES
// =============================================================================
// These methods maintain the French naming convention for existing code
// All new development should use the English method names above

/**
 * AjouterContact adds a contact using the legacy French method name
 *
 * @deprecated Use AddContact instead
 */
func (d *Directory) AjouterContact(nom, prenom, telephone string) error {
	return d.AddContact(nom, prenom, telephone)
}

/**
 * RechercherContact searches for a contact using the legacy French method name
 *
 * @deprecated Use SearchContact instead
 */
func (d *Directory) RechercherContact(nom string) (Contact, bool) {
	return d.SearchContact(nom)
}

/**
 * ListerContacts lists all contacts using the legacy French method name
 *
 * @deprecated Use ListContacts instead
 */
func (d *Directory) ListerContacts() []Contact {
	return d.ListContacts()
}

/**
 * SupprimerContact deletes a contact using the legacy French method name
 *
 * @deprecated Use DeleteContact instead
 */
func (d *Directory) SupprimerContact(nom string) error {
	return d.DeleteContact(nom)
}

/**
 * ModifierContact updates a contact using the legacy French method name
 *
 * @deprecated Use UpdateContact instead
 */
func (d *Directory) ModifierContact(nom, nouveauPrenom, nouveauTelephone string) error {
	return d.UpdateContact(nom, nouveauPrenom, nouveauTelephone)
}

/**
 * NombreContacts returns the contact count using the legacy French method name
 *
 * @deprecated Use ContactCount instead
 */
func (d *Directory) NombreContacts() int {
	return d.ContactCount()
}

/**
 * SaveToJSON exports to JSON using the legacy method name
 *
 * @deprecated Use ExportToJSON instead
 */
func (d *Directory) SaveToJSON(nomFichier string) error {
	return d.ExportToJSON(nomFichier)
}

/**
 * LoadFromJSON imports from JSON using the legacy method name
 *
 * @deprecated Use ImportFromJSON instead
 *
 * Note: For backward compatibility, this method doesn't fail if file doesn't exist
 * This differs from the new ImportFromJSON method which properly reports missing files
 */
func (d *Directory) LoadFromJSON(nomFichier string) error {
	// Legacy behavior: silently ignore missing files for backward compatibility
	if _, err := os.Stat(nomFichier); os.IsNotExist(err) {
		return nil
	}
	return d.ImportFromJSON(nomFichier)
}
