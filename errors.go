package main

type (
	NotFoundError struct {
		Message string `json:"message"`
	}
	InternalServerError struct {
		Message string `json:"message"`
	}
	BadRequestError struct {
		Message string            `json:"message"`
		Fields  map[string]string `json:"fields"`
	}
	ConflictError struct {
		Message string `json:"message"`
	}
)

func (error *NotFoundError) Error() string {
	bytes, _ := JSONEncoder(error)
	return string(bytes)
}

func (error *InternalServerError) Error() string {
	bytes, _ := JSONEncoder(error)
	return string(bytes)
}

func (error *BadRequestError) Error() string {
	bytes, _ := JSONEncoder(error)
	return string(bytes)
}

func (error *ConflictError) Error() string {
	bytes, _ := JSONEncoder(error)
	return string(bytes)
}
