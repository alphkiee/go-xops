package initialize

import (
	"go-xops/pkg/common"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
)

// 初始化校验器
func Validate() {
	// 实例化需要转换的语言, 中文
	chinese := zh.New()
	uni := ut.New(chinese, chinese)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()

	// 注册转换的语言为默认语言
	_ = translations.RegisterDefaultTranslations(validate, trans)

	common.Validate = validate
	common.Translator = trans
}
