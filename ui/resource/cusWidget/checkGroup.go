package cusWidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 创建一个 checkGroup

type CheckedGroupClickCallback func(label string, checked bool)
type CheckGroup struct {
	CheckList     []*widget.Check
	CheckedLabel  string
	ClickCallback CheckedGroupClickCallback
}

func (group *CheckGroup) Create(labels []string, h bool) *fyne.Container {
	group.CheckList = make([]*widget.Check, 0)
	for _, label := range labels {
		group.CheckList = append(
			group.CheckList,
			widget.NewCheck(label, func(b bool) {
				if b {
					for _, check := range group.CheckList {
						check.SetChecked(false)
					}
				}
				group.CheckedLabel = label

			}),
		)
	}
	if h {
		res := container.NewHBox()
		for _, check := range group.CheckList {
			res.Add(check)
		}
		return res
	} else {
		res := container.NewVBox()
		for _, check := range group.CheckList {
			res.Add(check)
		}
		return res
	}

}

type LabelAndInit struct {
	Label     string
	InitCheck bool
}

func CreateCheckGroup(labelsAndInit []LabelAndInit, h bool, singleCheck bool, clickCallback CheckedGroupClickCallback) *fyne.Container {

	group := CheckGroup{
		CheckList:     nil,
		CheckedLabel:  "",
		ClickCallback: clickCallback,
	}

	group.CheckList = make([]*widget.Check, 0)
	for i := 0; i < len(labelsAndInit); i++ {

		checkName := ""
		checkName += labelsAndInit[i].Label

		newCheck := widget.NewCheck(checkName, func(b bool) {
			if b && singleCheck {
				for _, otherCheck := range group.CheckList {
					if otherCheck.Text != checkName {
						otherCheck.SetChecked(false)
					}
				}
			}
			group.ClickCallback(checkName, b)
		})

		newCheck.SetChecked(labelsAndInit[i].InitCheck)

		group.CheckList = append(
			group.CheckList,
			newCheck,
		)

	}

	if h {
		res := container.NewHBox()
		for _, check := range group.CheckList {
			res.Add(check)
		}
		return res
	} else {
		res := container.NewVBox()
		for _, check := range group.CheckList {
			res.Add(check)
		}
		return res
	}

}

func CreateDropDown(labelsAndInit []LabelAndInit, clickCallback CheckedGroupClickCallback) *fyne.Container {
	options := make([]string, len(labelsAndInit))
	var selectedOption string

	for i, item := range labelsAndInit {
		options[i] = item.Label
		if item.InitCheck {
			selectedOption = item.Label
		}
	}

	selectWidget := widget.NewSelect(options, func(selected string) {
		for _, item := range labelsAndInit {
			if item.Label == selected {
				clickCallback(item.Label, true)
			} else {
				clickCallback(item.Label, false)
			}
		}
	})

	selectWidget.SetSelected(selectedOption)

	return container.NewVBox(selectWidget)
}
