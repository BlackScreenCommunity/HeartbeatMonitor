package plugins

type Plugin interface {
    Name() string
    Collect() (map[string]interface{}, error) 
}

