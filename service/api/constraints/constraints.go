package constraints

//TODO: Make sure all constraints are enforced in the backend

const MaxParticipants = 1000

const MaxGroupNameLength = 16
const MinGroupNameLength = 1

const MaxUsernameLength = 16
const MinUsernameLength = 3

const MaxMessageLength = 65536
const MinMessageLength = 1

var AllowedMimeTypes = []string{
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/webp",
}
