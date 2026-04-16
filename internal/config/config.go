package config

type Config struct {
	DryRun       bool     `json:"dryRun"`
	MongoURI     string   `json:"mongoUri"`
	Database     string   `json:"database"`
	DBData       string   `json:"dbdata"`
	WFExportData string   `json:"wfexportdata"`
	WFCustomData string   `json:"wfcustomdata"`
	Index        int      `json:"index"`
	Mode         string   `json:"mode"`
	Global       string   `json:"global"`
	Debug        bool     `json:"debug"`
	Args         []string `json:"args"`
}

type WFData struct {
	Exports Exports
	Custom  Custom
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
}

type Custom struct {
	WarframesU41 []byte
	WarframesU42 []byte
	AllScans     []byte
}
