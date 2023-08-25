package config

type Config struct {
	LogLevel              int      `mapstructure:"log_level" `
	ServerBind            string   `mapstructure:"server_bind"`
	GitlabURL             string   `mapstructure:"gitlab_url"`
	GitlabToken           string   `mapstructure:"gitlab_token"`
	GitlabUser            string   `mapstructure:"gitlab_user"`
	EsURL                 string   `mapstructure:"es_url"`
	SearchProvider        string   `mapstructure:"search_provider"`
	ReposPath             string   `mapstructure:"repos_path"`
	SyncProjectLimit      int      `mapstructure:"sync_Project_limit"`
	OAuthClientID         string   `mapstructure:"oauth_client_id"`
	OAuthClientSecret     string   `mapstructure:"oauth_client_secret"`
	TokenKey              string   `mapstructure:"token_key"`
	EnableAuth            bool     `mapstructure:"enable_auth"`
	HTTPHost              string   `mapstructure:"http_host"`
	DirBlacklist          []string `mapstructure:"dir_blacklist"`
	FileBlacklist         []string `mapstructure:"file_blacklist"`
	EnableCodeFullPreview bool     `mapstructure:"enable_code_full_preview"`
	MaxSearchTotal        int64    `mapstructure:"max_search_total"`
	MaxFileSize           int      `mapstructure:"max_file_size"`
	MaxLineLength         int      `mapstructure:"max_line_length"`
}

var Conf Config
