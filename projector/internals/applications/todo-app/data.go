package application

import "github.com/charmbracelet/bubbles/list"

func (b *App) initLists() {
	// Init To Do
	todos := newColumn(StatusTODO)
	todos.list.Title = "To do"
	todos.list.SetItems([]list.Item{
		Todo{status: StatusTODO, title: "buy milk", description: "strawberry milk"},
		Todo{status: StatusTODO, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Todo{status: StatusTODO, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	// Init done
	dones := newColumn(StatusDone)
	dones.list.Title = "Done"
	dones.list.SetItems([]list.Item{
		Todo{status: StatusDone, title: "stay cool", description: "as a cucumber"},
	})
	b.cols = map[status]*column{
		StatusTODO: &todos,
		StatusDone: &dones,
	}
}
