package response

import "gola/app/common/errorcode"

// WithData 回傳附帶資料的選項
func WithData(data interface{}) ArgsF {
	return func(opt *Option) {
		opt.Data = data
	}
}

// WithText 回傳客製化文字的選項
func WithText(text string) ArgsF {
	return func(opt *Option) {
		opt.ErrorText = text
	}
}

// WithStatusCode 回傳客製化狀態的選項
func WithStatusCode(status int) ArgsF {
	return func(opt *Option) {
		opt.HTTPCode = status
	}
}

// Option 選項
type Option struct {
	HTTPCode  int
	ErrorCode errorcode.Code
	ErrorText string
	Data      interface{}
}

// ArgsF 參數
type ArgsF func(opt *Option)

// WithOption 組合參數
func WithOption(opt *Option, args ...ArgsF) {
	if opt == nil {
		return
	}

	for _, fn := range args {
		if fn != nil {
			fn(opt)
		}
	}
}

// WithNewOption 組合參數，並回傳新的Option
func WithNewOption(args ...ArgsF) Option {
	opt := Option{}

	for _, fn := range args {
		if fn != nil {
			fn(&opt)
		}
	}

	return opt
}
