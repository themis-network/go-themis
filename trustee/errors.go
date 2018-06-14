package trustee

const (
	//NO_WINNER = -1
	//Decrypt_ERROR = -2
)

type noWinnerError struct{ message string }

//func (e *noWinnerError) ErrorCode() int { return -NO_WINNER }

func (e *noWinnerError) Error() string { return e.message }


type decryptError struct{ message string }

//func (e *decryptError) ErrorCode() int { return -Decrypt_ERROR }

func (e *decryptError) Error() string { return e.message }
