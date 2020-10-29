package setting

type ServerSettingS struct {
	RunModel string
}

type DatabaseSettingS struct {
	DBType       string
	Username     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

type AppSettingS struct {
	MaxPageSize     int
	DefaultPageSize int
}
