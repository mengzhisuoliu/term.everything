package wayland

//go:generate bash -c "go run ./generate ./protocols . $(go list) WlSurface XdgPositioner XdgSurface WlPointer WlSubsurface XdgToplevel"
