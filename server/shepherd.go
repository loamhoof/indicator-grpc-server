package main

import (
	"log"
	"path/filepath"

	"github.com/conformal/gotk3/gtk"
	"github.com/doxxan/appindicator"
	"github.com/doxxan/appindicator/gtk-extensions/gotk3"
	"golang.org/x/net/context"

	pb "github.com/loamhoof/indicator"
)

type indicatorShepherd struct {
	iconDir    string
	logger     *log.Logger
	indicators map[string]*gotk3.AppIndicatorGotk3
}

func newIndicatorShepherd(iconDir string, logger *log.Logger) *indicatorShepherd {
	return &indicatorShepherd{
		iconDir:    iconDir,
		logger:     logger,
		indicators: make(map[string]*gotk3.AppIndicatorGotk3),
	}
}

func (is *indicatorShepherd) Update(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	if err := is.ensureIndicator(req.Id); err != nil {
		is.logger.Println(err)

		return &pb.Response{Err: err.Error()}, err
	}

	if req.Label != "" {
		is.setLabel(req.Id, req.Label, req.LabelGuide)
	}
	if req.Icon != "" {
		is.setIcon(req.Id, req.Icon)
	}
	is.setStatus(req.Id, req.Active)

	return &pb.Response{}, nil
}

func (is *indicatorShepherd) ensureIndicator(id string) error {
	if _, ok := is.indicators[id]; ok {
		return nil
	}

	is.logger.Println(id, "newIndicator")

	indicator := gotk3.NewAppIndicator(id, "", appindicator.CategoryOther)

	menu, err := gtk.MenuNew()
	if err != nil {
		return err
	}

	menuItem, err := gtk.MenuItemNewWithLabel("")
	if err != nil {
		return err
	}

	menu.Append(menuItem)

	menuItem.Show()
	indicator.SetMenu(menu)

	is.indicators[id] = indicator

	return nil
}

func (is *indicatorShepherd) setLabel(id string, label, labelGuide string) {
	indicator := is.indicators[id]

	if currLabel, _ := indicator.GetLabel(); label != currLabel {
		is.logger.Println(id, "setLabel", label, labelGuide)

		indicator.SetLabel(label, labelGuide)
	}
}

func (is *indicatorShepherd) setIcon(id string, icon string) {
	indicator := is.indicators[id]

	iconFullPath := filepath.Join(is.iconDir, icon)

	if currIcon, _ := indicator.GetIcon(); iconFullPath != currIcon {
		is.logger.Println(id, "setIcon", icon)

		indicator.SetIcon(iconFullPath, "")
	}
}

func (is *indicatorShepherd) setStatus(id string, active bool) {
	indicator := is.indicators[id]

	var status appindicator.Status

	if active {
		status = appindicator.StatusActive
	} else {
		status = appindicator.StatusPassive
	}

	if indicator.GetStatus() != status {
		is.logger.Println(id, "setStatus", active)

		indicator.SetStatus(status)
	}
}
