package agency

type (
	// Agency stores information on an agency in Blossom
	Agency struct {
		// Name is the unique name of the agency
		Name string `json:"name"`
		// ATO is the Authority To Operate document
		ATO string `json:"ato"` // TODO string is placeholder for actual object
		// MSPID is the Membership Service Provider ID
		MSPID string `json:"mspid"`
		// Users contains the users of the organization that will access Blossom
		Users Users `json:"users"`
		// Status of an agency within the Blossom system
		Status Status `json:"status"`
	}

	// Status represents the status of an agency within the blossom system
	Status string

	// Users that will access blossom on behalf of an agency
	Users struct {
		// SystemOwner is responsible for administrative tasks for the agency system
		SystemOwner string `json:"system_owner"`
		// AcquisitionSpecialist authorizes transaction requests for the agency
		AcquisitionSpecialist string `json:"acquisition_specialist"`
		// SystemAdministrator interacts with the smart contracts to checkin and checkout software licenses for the agency
		SystemAdministrator string `json:"system_administrator"`
	}
)

// TODO finalize possible statuses for accounts
const (
	Pending  Status = "pending"
	Active   Status = "active"
	Inactive Status = "inactive"
)
