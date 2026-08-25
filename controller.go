package hazart

// Controller is an interface for registering a group of routes (e.g. UserController, AuthController)
type Controller interface {
	RegisterRoutes(g *Group)
}

// ControllerApp is an interface for registering routes directly to the App
type ControllerApp interface {
	RegisterRoutesApp(a *App)
}

// MountController mounts a controller struct under a group route prefix
func (g *Group) MountController(c Controller) {
	c.RegisterRoutes(g)
}

// MountController mounts a controller struct under an app group prefix
func (a *App) MountController(prefix string, c Controller) {
	group := a.Group(prefix)
	c.RegisterRoutes(group)
}
