package application

import (
	"errors"
	"os"

	"github.com/shogo82148/androidbinary/apk"
	"github.com/xmxu/aab-parser"
)

var (
	ErrAnalyze = errors.New("could not analyze app")
	ErrInfoAab = errors.New("could not get info for aab")
	ErrInfoApk = errors.New("could not get info for apk")
)

type ApplicationType int

func (t ApplicationType) String() string {
	switch t {
	case ApplicationTypeAab:
		return "AAB"
	case ApplicationTypeApk:
		return "APK"
	case ApplicationTypeUnknown:
		fallthrough
	default:
		return "Unknown"
	}
}

const (
	ApplicationTypeUnknown ApplicationType = iota
	ApplicationTypeApk
	ApplicationTypeAab
)

type ApplicationInfo struct {
	ApplicationType ApplicationType
	PackageName     string
	VersionName     string
	VersionCode     int64
}

type Analyzer struct{}

func New() Analyzer {
	return Analyzer{}
}

func (a Analyzer) Analyze(fileName string) (ApplicationInfo, error) {
	_, err := os.Stat(fileName)
	if err != nil {
		return ApplicationInfo{}, errors.Join(ErrAnalyze, err)
	}
	info, err := a.infoAab(fileName)
	if err == nil {
		return info, nil
	}
	info, err = a.infoApk(fileName)
	if err == nil {
		return info, nil
	}
	return ApplicationInfo{}, errors.Join(ErrAnalyze, err)
}

func (a Analyzer) infoAab(fileName string) (ApplicationInfo, error) {
	aab, err := aab.OpenFile(fileName)
	if err != nil {
		return ApplicationInfo{}, errors.Join(ErrInfoAab, err)
	}
	applicationInfo := ApplicationInfo{
		ApplicationType: ApplicationTypeAab,
		PackageName:     aab.PackageName(),
		VersionName:     aab.Manifest().VersionName,
		VersionCode:     aab.Manifest().VersionCode,
	}
	return applicationInfo, nil
}

func (a Analyzer) infoApk(fileName string) (ApplicationInfo, error) {
	apk, err := apk.OpenFile(fileName)
	if err != nil {
		return ApplicationInfo{}, errors.Join(ErrInfoApk, err)
	}
	versionName, err := apk.Manifest().VersionName.String()
	if err != nil {
		return ApplicationInfo{}, errors.Join(ErrInfoApk, err)
	}
	versionCode, err := apk.Manifest().VersionCode.Int32()
	if err != nil {
		return ApplicationInfo{}, errors.Join(ErrInfoApk, err)
	}
	applicationInfo := ApplicationInfo{
		ApplicationType: ApplicationTypeApk,
		PackageName:     apk.PackageName(),
		VersionName:     versionName,
		VersionCode:     int64(versionCode),
	}
	return applicationInfo, nil
}
