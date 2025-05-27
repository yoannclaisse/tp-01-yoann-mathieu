package annuaire

import (
	"testing"
)

// TestAddContact tests the AddContact functionality with various scenarios
func TestAddContact(t *testing.T) {
	dir := NewDirectory()

	// Test successful addition
	err := dir.AddContact("Dupont", "Jean", "0123456789")
	if err != nil {
		t.Errorf("Error during addition: %v", err)
	}

	// Test adding contact with same name but different phone
	err = dir.AddContact("Dupont", "Pierre", "0987654321")
	if err != nil {
		t.Errorf("Unexpected error when adding contact with same name but different phone: %v", err)
	}

	// Test adding with empty fields
	err = dir.AddContact("", "Test", "0000000000")
	if err == nil {
		t.Error("Expected error for empty name")
	}
}

// TestSearchContact tests the SearchContact functionality
func TestSearchContact(t *testing.T) {
	dir := NewDirectory()
	dir.AddContact("Martin", "Alice", "0123456789")

	// Test searching existing contact
	contact, exists := dir.SearchContact("Martin")
	if !exists {
		t.Error("Contact not found when it should exist")
	}

	if contact.First != "Alice" || contact.Phone != "0123456789" {
		t.Errorf("Incorrect data: %+v", contact)
	}

	// Test searching non-existent contact
	_, exists = dir.SearchContact("NonExistent")
	if exists {
		t.Error("Contact found when it shouldn't exist")
	}
}

// TestDeleteContact tests the DeleteContact functionality
func TestDeleteContact(t *testing.T) {
	dir := NewDirectory()
	dir.AddContact("Test", "Contact", "0000000000")

	err := dir.DeleteContact("Test")
	if err != nil {
		t.Errorf("Error during deletion: %v", err)
	}

	if dir.ContactCount() != 0 {
		t.Error("Contact was not deleted")
	}

	err = dir.DeleteContact("NonExistent")
	if err == nil {
		t.Error("Expected error for non-existent contact")
	}
}

// TestUpdateContact tests the UpdateContact functionality
func TestUpdateContact(t *testing.T) {
	dir := NewDirectory()
	dir.AddContact("Update", "Test", "0000000000")

	err := dir.UpdateContact("Update", "NewFirst", "1111111111")
	if err != nil {
		t.Errorf("Error during update: %v", err)
	}

	contact, _ := dir.SearchContact("Update")
	if contact.First != "NewFirst" || contact.Phone != "1111111111" {
		t.Errorf("Update failed: %+v", contact)
	}
}

// TestSearchContactWithMultipleSameNames tests searching when multiple contacts have the same last name
func TestSearchContactWithMultipleSameNames(t *testing.T) {
	dir := NewDirectory()

	// Add two contacts with same last name
	dir.AddContact("Bernard", "Jean", "0654321876")
	dir.AddContact("Bernard", "Pierre", "11111")

	// Search should find one of them (the first match)
	contact, exists := dir.SearchContact("Bernard")
	if !exists {
		t.Error("Contact not found when it should exist")
	}

	// Should find a Bernard (either Jean or Pierre)
	if contact.Name != "Bernard" {
		t.Errorf("Expected last name 'Bernard', got '%s'", contact.Name)
	}

	// Verify we can search by first name too
	contact, exists = dir.SearchContact("Jean")
	if !exists {
		t.Error("Contact 'Jean' not found")
	}
	if contact.First != "Jean" || contact.Name != "Bernard" {
		t.Errorf("Expected 'Jean Bernard', got '%s %s'", contact.First, contact.Name)
	}

	contact, exists = dir.SearchContact("Pierre")
	if !exists {
		t.Error("Contact 'Pierre' not found")
	}
	if contact.First != "Pierre" || contact.Name != "Bernard" {
		t.Errorf("Expected 'Pierre Bernard', got '%s %s'", contact.First, contact.Name)
	}
}

// TestImportAndAddFunctionality tests that imported and manually added contacts work together
func TestImportAndAddFunctionality(t *testing.T) {
	dir := NewDirectory()

	// Simulate importing a contact
	dir.AddContact("Bernard", "Jean", "0654321876")

	// Simulate manually adding another contact
	dir.AddContact("Bernard", "Pierre", "11111")

	// Should have 2 contacts
	if dir.ContactCount() != 2 {
		t.Errorf("Expected 2 contacts, got %d", dir.ContactCount())
	}

	// Should be able to find both by last name (will return first match)
	_, exists := dir.SearchContact("Bernard")
	if !exists {
		t.Error("No Bernard found")
	}

	// Should be able to find by first names
	jean, exists := dir.SearchContact("Jean")
	if !exists {
		t.Error("Jean not found")
	}
	if jean.Phone != "0654321876" {
		t.Errorf("Expected Jean's phone to be 0654321876, got %s", jean.Phone)
	}

	pierre, exists := dir.SearchContact("Pierre")
	if !exists {
		t.Error("Pierre not found")
	}
	if pierre.Phone != "11111" {
		t.Errorf("Expected Pierre's phone to be 11111, got %s", pierre.Phone)
	}
}
