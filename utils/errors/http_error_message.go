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

	// rcic error message
	CreateHTTPRequestFailed   = "create HTTP request failed"
	RequestRemoteServerFailed = "request remote server failed"
	ReadResponseBodyFailed    = "read response body failed"

	// post create message
	NotValidSubject    = "not valid subject"
	NotValidContent    = "not valid content"
	NotValidCategories = "not valid categories"
	NotValidImages     = "not valid images"

	// search page message
	NotValidLimit = "not valid Limit"

	// delete post
	OwnerNotMatch   = "owner not match"
	InvalidImageUrl = "invalid image URL"

	//update profile

	UseNameInvalid = "use name in invalid"

	// message module
	SubjectInvalid = "subject in invalid"
	ContentInvalid = "content in invalid"
)
