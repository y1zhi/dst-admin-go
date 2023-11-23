package constant

import (
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/systemUtils"
	"fmt"
	"path/filepath"
)

var HOME_PATH string

// var dstConfig = dstConfigUtils.GetDstConfig()

func init() {
	home, err := systemUtils.Home()
	if err != nil {
		panic("Home path error: " + err.Error())
	}
	HOME_PATH = home
	fmt.Println("home path: " + HOME_PATH)
}

var (
	PASSWORD_PATH = "./password.txt"

	//饥荒安装位置
	DST_INSTALL_DIR = dstConfigUtils.GetDstConfig().Force_install_dir
	//steam cmd 安装的位置
	STEAMCMD = dstConfigUtils.GetDstConfig().Steamcmd
	//房间位置
	CLUSTER = dstConfigUtils.GetDstConfig().Cluster

	/**
	 * 全局的地面进程、存档的名称
	 */
	DST_MASTER = "Master"

	/**
	 * 全局的洞穴进程、存档的名称
	 */
	DST_CAVES = "Caves"

	/**
	 * 地面的screen任务的名称 DST_MASTER
	 */
	SCREEN_WORK_MASTER_NAME = "DST_MASTER"

	/**
	 * 洞穴的screen任务的名称 DST_CAVES
	 */
	SCREEN_WORK_CAVES_NAME = "DST_CAVES"

	/**
	 * 启动地面进程命令 设置名称为 DST_MASTER
	 */
	START_MASTER_CMD = "cd " + DST_INSTALL_DIR + "/bin ; screen -d -m -S \"" + SCREEN_WORK_MASTER_NAME + "\"  ./dontstarve_dedicated_server_nullrenderer -console -cluster " + CLUSTER + " -shard " + DST_MASTER + "  ;"
	// cd ~/dst/bin/ ; screen -d -m -S \"DST_MASTER\"  ./dontstarve_dedicated_server_nullrenderer -console -level MyDediServer -shard Master  ;

	/**
	 * 启动洞穴进程命令 设置名称为 DST_CAVES
	 */
	START_CAVES_CMD = "cd " + DST_INSTALL_DIR + "/bin ; screen -d -m -S \"" + SCREEN_WORK_CAVES_NAME + "\"  ./dontstarve_dedicated_server_nullrenderer -console -cluster " + CLUSTER + " -shard " + DST_CAVES + " ;"

	/**
	 * 检查目前所有的screen作业，并删除已经无法使用的screen作业
	 */
	CLEAR_SCREEN_CMD = "screen -wipe "

	/**
	 * 查询地面进程号命令
	 */
	FIND_MASTER_CMD = " ps -ef | grep -v grep |grep '" + DST_MASTER + "'|sed -n '1P'|awk '{print $2}' "

	/**
	 * 查询洞穴进程号命令
	 */
	FIND_CAVES_CMD = " ps -ef | grep -v grep |grep '" + DST_CAVES + "'|sed -n '1P'|awk '{print $2}' "

	/**
	 * 杀死地面进程命令
	 */
	STOP_MASTER_CMD = "ps -ef | grep -v grep |grep '" + DST_MASTER + "' |sed -n '1P'|awk '{print $2}' |xargs kill -9"

	/**
	 * 杀死洞穴进程命令
	 */
	STOP_CAVES_CMD = "ps -ef | grep -v grep |grep '" + DST_CAVES + "' |sed -n '1P'|awk '{print $2}' |xargs kill -9"

	/**
	 * 更新游戏目录
	 */
	UPDATE_GAME_CMD = "cd " + STEAMCMD + " ; ./steamcmd.sh +login anonymous +force_install_dir " + DST_INSTALL_DIR + " +app_update 343050 validate +quit"

	/**
	 * 删除地面游戏记录
	 */
	DEL_RECORD_MASTER_CMD = "rm -r ~/.klei/DoNotStarveTogether/" + CLUSTER + "/" + DST_MASTER + "/save"

	/**
	 * 删除地面游戏记录
	 */
	DEL_RECORD_CAVES_CMD = "rm -r ~/.klei/DoNotStarveTogether/" + CLUSTER + "/" + DST_CAVES + "/save"

	/**
	 * 获取地面的玩家 替换99999999关键字
	 */
	MASTER_PLAYLIST_CMD = "screen -S \"" + SCREEN_WORK_MASTER_NAME + "\" -p 0 -X stuff \"for i, v in ipairs(TheNet:GetClientTable()) do  print(string.format(\\\"playerlist %s [%d] %s %s %s\\\", 99999999, i-1, v.userid, v.name, v.prefab )) end$(printf \\\\r)\"\n"

	/**
	 * 饥荒的启动程序
	 */
	DST_START_PROGRAM = "dontstarve_dedicated_server_nullrenderer"

	/**
	 * 单斜杠
	 */
	SINGLE_SLASH = "/"

	/**
	 * 备份的存档文件的扩展名
	 */
	BACKUP_FILE_EXTENSION = ".tar"

	/**
	 * 备份的存档文件的扩展名
	 */
	BACKUP_FILE_EXTENSION_NON_POINT = "tar"
	/**
	 * 备份的存档文件的扩展名zip
	 */
	BACKUP_FILE_EXTENSION_NON_POINT_ZIP = "zip"

	/**
	 * 不允许下载文件路径中存在改字符
	 */
	BACKUP_ERROR_PATH = "../"

	/**
	 * 游戏文档
	 */
	DST_DOC_PATH = ".klei/DoNotStarveTogether"

	/**
	 * 饥荒游戏用户存档位置
	 */
	DST_USER_GAME_CONFG_PATH = "/.klei/DoNotStarveTogether/" + CLUSTER

	/**
	 * 饥荒游戏存档路径
	 */
	DST_USER_SAVE_FILE = "save"

	/**
	 * 游戏配置的名称
	 */
	DST_USER_SERVER_INI_NAME = "server.ini"

	/**
	 * 游戏房间设置的文件名称
	 */
	DST_USER_CLUSTER_INI_NAME = "cluster.ini"

	/**
	 * token设置文件
	 */
	DST_USER_CLUSTER_TOKEN = "cluster_token.txt"

	/**
	 * 地上mod保存地址
	 */
	DST_USER_GAME_MASTER_MOD_PATH = ".klei/DoNotStarveTogether/" + CLUSTER + "/" + DST_MASTER + "/modoverrides.lua"

	/**
	 * 洞穴mod保存位置
	 */
	DST_USER_GAME_CAVES_MOD_PATH = ".klei/DoNotStarveTogether/" + CLUSTER + "/" + DST_CAVES + "/modoverrides.lua"

	/**
	 * 地面地图配置地址
	 */
	DST_USER_GAME_MASTER_MAP_PATH = ".klei/DoNotStarveTogether/" + CLUSTER + "/" + DST_MASTER + "/leveldataoverride.lua"

	/**
	 * 洞穴地图配置地址
	 */
	DST_USER_GAME_CAVES_MAP_PATH = ".klei/DoNotStarveTogether/" + CLUSTER + "/" + DST_CAVES + "/leveldataoverride.lua"

	/**
	 * 游戏配置文件
	 */
	DST_USER_GAME_CONFIG_PATH = ".klei/DoNotStarveTogether/" + CLUSTER + "/cluster.ini"

	/**
	 * 地面游戏运行日志位置
	 */
	DST_MASTER_SERVER_LOG_PATH = ".klei/DoNotStarveTogether/" + CLUSTER + "/" + DST_MASTER + "/server_log.txt"

	/**
	 * 地面用户聊天信息
	 */
	DST_MASTER_SERVER_CHAT_LOG_PATH = ".klei/DoNotStarveTogether/" + CLUSTER + "/" + DST_MASTER + "/server_chat_log.txt"

	/**
	 * 洞穴游戏运行日志位置
	 */
	DST_CAVES_SERVER_LOG_PATH = ".klei/DoNotStarveTogether/" + CLUSTER + "/" + DST_CAVES + "/server_log.txt"

	/**
	 * 管理员存储位置
	 */
	DST_ADMIN_LIST_PATH = ".klei/DoNotStarveTogether/" + CLUSTER + "/adminlist.txt"

	/**
	 * 黑名单存储位置
	 */
	DST_PLAYER_BLOCK_LIST_PATH = ".klei/DoNotStarveTogether/" + CLUSTER + "/blocklist.txt"

	/**
	 * 游戏mod设置
	 */
	DST_MOD_SETTING_PATH = DST_INSTALL_DIR + "/mods/dedicated_server_mods_setup.lua"

	/**
	 * master的session目录
	 */
	DST_USER_GAME_MASTER_SESSION = DST_MASTER + "/save/session"
)

func GET_DST_USER_GAME_CONFG_PATH() string {
	dstConfig := dstConfigUtils.GetDstConfig()
	cluster := dstConfig.Cluster
	donot_starve_server_directory := dstConfig.DoNotStarveServerDirectory
	persistent_storage_root := dstConfig.Persistent_storage_root
	if donot_starve_server_directory == "" {
		filepath.Join(HOME_PATH, "Documents", "klei", "DoNotStarveTogether", cluster)
		return filepath.Join(HOME_PATH, "Documents", "klei", "DoNotStarveTogether", cluster)
	}

	return filepath.Join(persistent_storage_root, donot_starve_server_directory, cluster)
}
