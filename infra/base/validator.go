package base

import (
	"errors"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	vtzh "gopkg.in/go-playground/validator.v9/translations/zh"
	"one/infra"
)

var validate *validator.Validate
var translator ut.Translator

func Validate() *validator.Validate {
	Check(validate)
	return validate
}

func Translate() ut.Translator {
	Check(translator)
	return translator
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	validate = validator.New()
	//创建消息国际化通用翻译器
	cn := zh.New()
	uni := ut.New(cn, cn)
	var found bool
	translator, found = uni.GetTranslator("zh")
	if found {
		err := vtzh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			log.Error(err)
		}
	} else {
		log.Error("Not found translator: zh")
	}
}

func (v *ValidatorStarter) Setup(ctx infra.StarterContext) {
	// 自定义验证器
	_ = validate.RegisterValidation("onSale", onSale)
	// 自定义验证输出
	_ = validate.RegisterTranslation("onSale", translator, func(ut ut.Translator) error {
		return ut.Add("onSale", "{0} 是一个无效的sku!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("onSale", fe.Field())
		return t
	})
}

func ValidateStruct(s interface{}) (err error) {
	err = Validate().Struct(s)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			log.Error("验证错误", err)
			return err
		}
		// 在终端打印，接口也需要显示
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			var slice string
			for _, e := range errs {
				log.Error(e.Translate(Translate()))
				slice += e.Translate(Translate()) + " | "
			}
			return errors.New(slice)
		}
	}
	return nil
}
