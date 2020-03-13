package config

type (
	FileConfig struct {
		Log LogConfig
	}

	LogConfig struct {
		//EnableConsole bool
		//ConsoleLevel string
		EnableFile  bool
		FileLevel   string
		Filename    string
		Skip        int
		FileMaxSize int
		FileMaxAge  int
		Compress    bool
	}
)
