package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/jblim0125/redredis-expire/common/appdata"
)

// Version for version
type Version struct{}

// New create versio ninstance.
func (Version) New() *Version {
	return &Version{}
}

// GetVersion return app version
// @Summary Get Server Version
// @Description get server version info
// @Tags version
// @Accept  json
// @Produce  json
// @success 200 {object} controllers.HTTPResponse{data=appdata.VersionInfo} "app info(name, version, hash)"
// @Router /version [get]
func (controller *Version) GetVersion(c echo.Context) error {
	res := HTTPResponse{}.ReturnSuccess(
		&appdata.VersionInfo{
			Name:      appdata.Name,
			Version:   appdata.Version,
			BuildHash: appdata.BuildHash})
	return c.JSON(http.StatusOK, res)
}
