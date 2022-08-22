package utils

import (
	"ARTeamViewService/global"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func OssDeleteFile(arrSidPrefix []string) {
	endpoint := global.GConfig.OssEndpoint
	accessKeyId := global.GConfig.StorageConfig.AccessKey
	accessKeySecret := global.GConfig.StorageConfig.SecretKey
	bucketName := global.GConfig.StorageConfig.Bucket
	//fileNamePrefix := global.GConfig.StorageConfig.FileNamePrefix

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		global.GLogger.Error("oss.New err: ", err)
		return
	}

	// yourBucketName填写存储空间名称。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		global.GLogger.Error("client.Bucket err: ", err)
		return
	}

	//ar/synergy/X2177364665_15717986
	//tmpPrefix := ""
	//
	//if len(fileNamePrefix) > 0 {
	//	for _, v := range fileNamePrefix {
	//		tmpPrefix += v + "/"
	//	}
	//} else {
	//	global.GLogger.Info("fileNamePrefix: ", fileNamePrefix)
	//	tmpPrefix += "/"
	//}

	if len(arrSidPrefix) > 0 {
		for _, tmpPrefix := range arrSidPrefix {
			//tmpPrefix := sid
			global.GLogger.Error("tmpPrefix: ", tmpPrefix)

			// 列举所有包含指定前缀的文件并删除。
			marker := oss.Marker("")
			// 如果您需要删除所有前缀为src的文件，则prefix设置为src。设置为src后，所有前缀为src的非目录文件、src目录以及目录下的所有文件均会被删除。
			prefix := oss.Prefix(tmpPrefix)
			// 如果您仅需要删除src目录及目录下的所有文件，则prefix设置为src/。
			// prefix := oss.Prefix("src/")
			count := 0
			for {
				lor, err := bucket.ListObjects(marker, prefix)
				if err != nil {
					global.GLogger.Error("bucket.ListObjects err: ", err)
					break
				}

				objects := []string{}
				for _, object := range lor.Objects {
					objects = append(objects, object.Key)
				}
				// 将oss.DeleteObjectsQuiet设置为true，表示不返回删除结果。
				delRes, err := bucket.DeleteObjects(objects, oss.DeleteObjectsQuiet(true))
				if err != nil {
					global.GLogger.Error("bucket.DeleteObjects err: ", err)
					break
				}

				if len(delRes.DeletedObjects) > 0 {
					global.GLogger.Error("these objects deleted failure, delRes.DeletedObjects: ", delRes.DeletedObjects)
					break
				}

				count += len(objects)

				prefix = oss.Prefix(lor.Prefix)
				marker = oss.Marker(lor.NextMarker)
				if !lor.IsTruncated {
					break
				}
			}
			global.GLogger.Info("OssDeleteFile tmpPrefix: ", tmpPrefix, ", delete object count: ", count)
		}
	} else {
		global.GLogger.Error("arrSidPrefix: ", arrSidPrefix)
	}
}
