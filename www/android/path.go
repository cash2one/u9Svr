package androidPackage

import (
	"strconv"
	"time"
	"u9/models"
)

const (
	productRootPath = "package/game"
	channelRootPath = "package/channel"
	packageRootPath = "package/uncompile"
	buildIdRootPath = "package/build_id"
	u9sdkFile = "package/sdk/1.1/u9.tar"
)

func GetApkName(product *models.Product, productVersion *models.ProductVersion) string {
	return strconv.Itoa(product.CpId) + "_" +
		strconv.Itoa(product.Id) + "_" + productVersion.VersionName
}

func GetProductPath(product *models.Product, apkName string) string {
	return productRootPath + "/" + strconv.Itoa(product.CpId) + "/" +
		strconv.Itoa(product.Id) + "/" + apkName
}

func GetChannelPath(channel *models.Channel) string {
	return channelRootPath + "/" + channel.Type
}

func GetPackagePath(packageTaskId int, apkName string) (ret string) {
	if apkName == "" {
		ret = packageRootPath + "/" + strconv.Itoa(packageTaskId)
	} else {
		ret = packageRootPath + "/" + strconv.Itoa(packageTaskId) + "/" + apkName
	}
	return
}

func GetPublishPath(product *models.Product, productVersion *models.ProductVersion) (ret string) {
	return "publish/" + strconv.Itoa(product.CpId) + "/" + strconv.Itoa(product.Id) +
		"/" + productVersion.VersionName
}

func GetPublishApk(product *models.Product, productVersion *models.ProductVersion, channel *models.Channel) (ret string) {
	publishPath := GetPublishPath(product, productVersion)
	packageTime := time.Now().Format("20060102150405")

	return publishPath + "/" + product.Code + "_" + channel.Type + "_" +
		productVersion.VersionName + "_" + packageTime + ".apk"
}

func GetPackageIconPath(channel *models.Channel, packageParam *models.PackageParam, productVersion *models.ProductVersion) (channelRet string, gameRet string) {
	iconType := packageParam.IconType
	switch iconType {
	case 1:
		return "", productVersion.IconUrl
	case 2:
		return channel.IconLeftTop, productVersion.IconUrl
	case 3:
		return channel.IconLeftBottom, productVersion.IconUrl
	case 4:
		return channel.IconRightTop, productVersion.IconUrl
	case 5:
		return channel.IconRightBottom, productVersion.IconUrl
	case 6:
		return "", packageParam.PackageIcon
	}
	return
}

func GetBuildIdPath(packageTaskId int, project string) (ret string) {
	if project == "" {
		ret = buildIdRootPath + "/" + strconv.Itoa(packageTaskId)
	} else {
		ret = buildIdRootPath + "/" + strconv.Itoa(packageTaskId) + project
	}
	return
}

func GetU9SDKFile() (ret string){
	return u9sdkFile
}
