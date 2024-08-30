package errors

// custom error message
const (
	JSONParseError = "JSON Parse Error"
	UnmarshalError = "unmarshal Error Error"
	// Create account
	NotValidEmail        = "not valid email"
	PasswordError        = "password error"
	InternalError        = "internal error"
	InsertDBError        = "insert DB error"
	AccountAlreadyExists = "account already exists"

	// Login error message
	AccountNotExist        = "account not exist"
	PasswordDecryptedError = "password decrypt error"
	PasswordIncorrect      = "password incorrect"
	UserProfileNotFound    = "user profile not found"

	//update profile
	UseNameInvalid = "use name in invalid"

	// message module
	SubjectInvalid = "subject in invalid"
	ContentInvalid = "content in invalid"

	// common
	DBProcessError = "db process error"

	TokenNotFound   = "token not found"
	TokenConsumed   = "token is already consumed"
	TokenHasExpired = "Token has expired"
)
