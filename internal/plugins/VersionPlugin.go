package plugins

type VersionPlugin struct{}

func (v VersionPlugin) Name() string {
    return "VersionPlugin"
}

func (v VersionPlugin) Collect() (map[string]interface{}, error) {
    return map[string]interface{}{
        "version": "1.0.0",
        "status":  "stable",
    }, nil
}