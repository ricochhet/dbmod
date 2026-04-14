package config

type Config struct {
	DryRun   bool     `json:"dryRun"`
	MongoURI string   `json:"mongoUri"`
	Database string   `json:"database"`
	DBData   string   `json:"dbdata"`
	WFData   string   `json:"wfdata"`
	Index    int      `json:"index"`
	Mode     string   `json:"mode"`
	Global   string   `json:"global"`
	Debug    bool     `json:"debug"`
	Args     []string `json:"args"`
}

type Exports struct {
	Achievements []byte
	Codex        []byte
	Customs      []byte
	Enemies      []byte
	Resources    []byte
	Virtuals     []byte
	Flavor       []byte
	Regions      []byte
	Weapons      []byte
	Warframes    []byte
	Sentinels    []byte
	AllScans     []byte
}
