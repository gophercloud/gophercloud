package gophercloud

type Context struct {
	providerMap map[string]*Provider
}

func TestContext() *Context {
	return &Context{
		providerMap: make(map[string]*Provider),
	}
}
