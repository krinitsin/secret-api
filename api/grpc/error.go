package grpc

import (
	rpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcError struct{
	error
	grpcCode rpccodes.Code
}

type Error interface{
	error
	ErrorWithCode() error
	UpdateCode(code rpccodes.Code)
}

func (err grpcError) ErrorWithCode() error {
	return status.Errorf(err.grpcCode, err.Error())
}

func NewError(code rpccodes.Code,err error) *grpcError {
	return &grpcError{
		error:    err,
		grpcCode: code,
	}
}

func (err grpcError)UpdateCode(code rpccodes.Code) {
	err.grpcCode = code
}



