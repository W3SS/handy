package handy

type FactoriesMap map[string]func() interface{}
type ProvidersMap map[string]func(*Context) func() interface{}


type Contextualized interface {
	Context() *Context
}

type Context struct {
	Context                   *Context
	factoriesMap              FactoriesMap
	providersMap              ProvidersMap
	gc                        []func()
}

func (context *Context) GC() {
	for _, v := range context.gc {
		v()
	}
}

func (context *Context) CleanupFunc(a func()) {
	context.gc = append(context.gc, a)
}


func (t *Context) SetFactory(key string, loader func() interface{}) *Context {
	if t.factoriesMap == nil {
		t.factoriesMap = FactoriesMap{}
	}
	t.factoriesMap[key] = loader
	return t
}

func (t *Context) SetProvider(key string, provider func(*Context) func() interface{}) *Context {
	if t.providersMap == nil {
		t.providersMap = ProvidersMap{}
	}
	t.providersMap[key] = provider
	return t
}


func (context *Context) SetValue(key string, value interface{}) *Context {
	context.SetFactory(key, func() interface{} {
		return value
	})
	return context
}


func (t *Context) MapFactories(factoriesMap FactoriesMap) *Context {
	for k, v := range factoriesMap {
		t.SetFactory(k, v)
	}
	return t
}


func (t *Context) MapProviders(providersMap ProvidersMap) *Context {
	for k, v := range providersMap {
		t.SetProvider(k, v)
	}
	return t
}


func (context *Context) GetFactory(key string) func() interface{} {

	if factory, ok := context.factoriesMap[key]; ok {
		return factory
	}

	if context.Context != nil {
		return context.Context.GetFactory(key)
	}

	return nil
}


func (context *Context) GetProvider(key string) func(*Context) func() interface{} {
	if provider, ok := context.providersMap[key]; ok {
		return provider
	}
	if context.Context != nil {
		return context.Context.GetProvider(key)
	}
	return nil
}


func (t *Context) Get(key string) interface{} {

	if interFC, ok := t.factoriesMap[key]; ok {
		return interFC()
	}

	provider := t.GetProvider(key)

	if provider != nil {
		loader := provider(t)
		t.SetFactory(key, loader)
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
