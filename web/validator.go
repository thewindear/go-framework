package web

import (
    "github.com/go-playground/locales/zh"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    zhTranslations "github.com/go-playground/validator/v10/translations/zh"
    "reflect"
    "strings"
)

type InvalidField struct {
    Field string `json:"field"`
    Tag   string `json:"tag"`
    Value string `json:"value"`
    Error string `json:"error"`
}

var validate = validator.New()
var cn = zh.New()
var uni = ut.New(cn, cn)
var trans, _ = uni.GetTranslator("zh")

func init() {
    _ = zhTranslations.RegisterDefaultTranslations(validate, trans)
}

func ValidateStruct(data interface{}) []*InvalidField {
    var errors []*InvalidField
    err := validate.Struct(data)
    elem := reflect.TypeOf(data).Elem()
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            var element InvalidField
            structField, has := elem.FieldByNameFunc(func(s string) bool {
                return s == err.Field()
            })
            if !has {
                continue
            }
            element.Field = structField.Tag.Get("json")
            element.Tag = err.Tag()
            element.Value = err.Param()
            element.Error = strings.ReplaceAll(err.Translate(trans), err.Field(), element.Field)
            errors = append(errors, &element)
        }
    }
    return errors
}
