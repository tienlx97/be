package log

// https://www.sobyte.net/post/2022-03/uber-zap-advanced-usage/
func InitLog() {
	var tops = []TeeOption{
		{
			Filename: "access.log",
			Ropt: RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 3,
				Compress:   true,
			},
			Lef: func(lvl Level) bool {
				return lvl <= InfoLevel
			},
		},
		{
			Filename: "error.log",
			Ropt: RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 3,
				Compress:   true,
			},
			Lef: func(lvl Level) bool {
				return lvl > InfoLevel
			},
		},
	}

	logger := NewTeeWithRotate(tops)
	ResetDefault(logger)
}
