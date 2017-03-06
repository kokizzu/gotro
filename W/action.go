package W

type Action func(controller *Context)
type Routes map[string]Action
