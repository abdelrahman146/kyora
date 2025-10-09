package webcontext

type SideBarItem struct {
	Icon  string
	Title string
	Link  string
}

func GetSideBarItems() []SideBarItem {
	return []SideBarItem{
		{Title: "Home", Link: "/", Icon: "ti ti-home"},
		{Title: "Profile", Link: "/profile", Icon: "ti ti-user"},
		{Title: "Settings", Link: "/settings", Icon: "ti ti-settings"},
	}
}
