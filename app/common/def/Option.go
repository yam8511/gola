package def

// Option 選項
type Option struct {
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
