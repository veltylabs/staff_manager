package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	staffmanager "github.com/veltylabs/staff_manager"
)

func TestStaffByUserFound(t *testing.T) {
	m := setup(t, nil)
	// Insert staff "s1" with user_id "u1", name "Ana", specialty "Pediatría"
	_, err := m.UpsertStaff(staffmanager.StaffMember{
		Id:        "s1",
		UserId:    "u1",
		Name:      "Ana",
		Specialty: "Pediatría",
		Rut:       "1-9",
		IsActive:  true,
		TenantId:  testTenant,
	})
	assert.NoError(t, err)

	staffID, name, specialty, found, err := m.StaffByUser(testTenant, "u1")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "s1", staffID)
	assert.Equal(t, "Ana", name)
	assert.Equal(t, "Pediatría", specialty)
}

func TestStaffByUserMissing(t *testing.T) {
	m := setup(t, nil)

	staffID, name, specialty, found, err := m.StaffByUser(testTenant, "u9")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Equal(t, "", staffID)
	assert.Equal(t, "", name)
	assert.Equal(t, "", specialty)
}

func TestStaffByUserOtherTenant(t *testing.T) {
	m := setup(t, nil)
	// staff of testTenant with user "u1"
	_, err := m.UpsertStaff(staffmanager.StaffMember{
		Id:        "s2",
		UserId:    "u1",
		Name:      "Bob",
		Specialty: "General",
		Rut:       "2-7",
		IsActive:  true,
		TenantId:  testTenant,
	})
	assert.NoError(t, err)

	// query with tenant "t1" which is not testTenant
	staffID, name, specialty, found, err := m.StaffByUser("t1", "u1")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Equal(t, "", staffID)
	assert.Equal(t, "", name)
	assert.Equal(t, "", specialty)
}
