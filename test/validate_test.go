package test

import (
    "github.com/thewindear/go-web-framework/web"
    "testing"
)

type User struct {
    Name string `validate:"required,min=3,max=32" json:"name"`
    // use `*bool` here otherwise the validation will fail for `false` values
    // Ref: https://github.com/go-playground/validator/issues/319#issuecomment-339222389
    IsActive *bool  `validate:"required" json:"isActive"`
    Email    string `validate:"required,email,min=6,max=32" json:"email"`
}

func TestValidator(t *testing.T) {
    b := false
    user := &User{
        Name:     "abc",
        IsActive: &b,
        Email:    "",
    }
    t.Log(web.ValidateStructJson(user)[0])
}
