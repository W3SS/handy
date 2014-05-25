package handy

type ContextLoaderMap map[string]func() interface{}
type ContextProviderMap map[string]func(*Context) func() interface{}

type Context struct {
	Context *Context
	Map             ContextLoaderMap
	ProviderMap     ContextProviderMap
	gc              []func()
}

func (context *Context) GC() {
	for _, v := range context.gc {
		v()
	}
}

func (context *Context) CleanupFunc(a func()) {
	context.gc = append(context.gc, a)
}


func (t *Context) Set(key string, loader func() interface{}) *Context {
	if t.Map == nil {
		t.Map = ContextLoaderMap{}
	}
	t.Map[key] = loader
	return t
}

func (t *Context) SetProvider(key string, provider func(*Context) func() interface{}) *Context {
	if t.ProviderMap == nil {
		t.ProviderMap = ContextProviderMap{}
	}
	t.ProviderMap[key] = provider
	return t
}


func (context *Context) SetValue(key string, value interface{}) *Context {
	context.Set(key, func() interface{} {
		return value
	})
	return context
}


func (t *Context) SetMap(maps ...ContextLoaderMap) *Context {
	for _, map_value := range maps {
		for k, v := range map_value {
			t.Set(k, v)
		}
	}
	return t
}


func (t *Context) SetProviderMap(maps ...ContextProviderMap) *Context {
	for _, map_value := range maps {
		for k, v := range map_value {
			t.SetProvider(k, v)
		}
	}
	return t
}

func (context *Context) GetProvider(key string) func(*Context) func() interface{} {
	if provider, ok := context.ProviderMap[key]; ok {
		return provider
	}
	if context.Context != nil {
		return context.Context.GetProvider(key)
	}
	return nil
}


func (t *Context) Get(key string) interface{} {

	if interfc, ok := t.Map[key]; ok {
		return interfc()
	}

	provider := t.GetProvider(key)

	if provider != nil {
		loader := provider(t)
		t.Set(key, loader)
		return loader()
	}

	return nil
}

func (t *Context) NewContext() *Context {
	return &Context{Context:t}
}

func NewContext() *Context {
	return &Context{}
}
