package service

import (
	"dst-admin-go/constant/consts"
	dst_cli_window "dst-admin-go/dst-cli-window"
	"dst-admin-go/model"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/levelConfigUtils"
	"dst-admin-go/utils/systemUtils"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"

	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/vo"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type GameService struct {
	lock      sync.Mutex
	c         HomeService
	logRecord LogRecordService
}

func (g *GameService) GetLastDstVersion() int64 {

	url := "http://ver.tugos.cn/getLocalVersion"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	s := string(body)
	veriosn, err := strconv.Atoi(s)
	if err != nil {
		veriosn = 0
	}
	return int64(veriosn)
}

func (g *GameService) GetLocalDstVersion(clusterName string) int64 {
	cluster := clusterUtils.GetCluster(clusterName)
	versionTextPath := filepath.Join(cluster.ForceInstallDir, "version.txt")
	log.Println("versionTextPath", versionTextPath)
	version, err := fileUtils.ReadFile(versionTextPath)
	if err != nil {
		log.Println(err)
		return 0
	}
	version = strings.Replace(version, "\r", "", -1)
	version = strings.Replace(version, "\n", "", -1)
	l, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	return l
}

func ClearScreen() bool {
	result, err := shellUtils.Shell(consts.ClearScreenCmd)
	if err != nil {
		return false
	}
	res := strings.Split(result, "\n")[0]
	return res != ""
}

func (g *GameService) UpdateGame(clusterName string) error {

	g.lock.Lock()
	defer g.lock.Unlock()
	// TODO 关闭相应的世界
	// g.StopGame(clusterName)

	updateGameCMd := dstUtils.GetDstUpdateCmd(clusterName)
	log.Println("正在更新游戏", "cluster: ", clusterName, "command: ", updateGameCMd)
	result, err := shellUtils.ExecuteCommandInWin(updateGameCMd)
	log.Println(result)

	levelConfig, _ := levelConfigUtils.GetLevelConfig(clusterName)
	for i := range levelConfig.LevelList {
		level := homeServe.GetLevel(clusterName, levelConfig.LevelList[i].File)
		modoverrides := level.Modoverrides
		dstUtils.DedicatedServerModsSetup2(clusterName, modoverrides)
	}

	return err
}

func (g *GameService) GetLevelStatus(clusterName, level string) bool {

	//start := time.Now().Nanosecond()
	//// 替换为你要查找的窗口标题
	//targetWindowTitle := clusterName + "_" + level
	//
	//// 构建 PowerShell 命令
	//psCommand := fmt.Sprintf("Get-Process | Where-Object { $_.MainWindowTitle -eq '%s' } | Format-Table -Property Id, ProcessName, MainWindowTitle", targetWindowTitle)
	//
	//// 执行 PowerShell 命令
	//cmd := exec.Command("powershell", "-Command", psCommand)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//err := cmd.Run()
	//if err != nil {
	//	fmt.Println("Error executing PowerShell command:", err)
	//	return false
	//}
	//
	//// 打印输出结果
	//// fmt.Println(out.String())
	//// log.Println("查询世界状态", cmd, out.String(), out.String() != "")
	//end := time.Now().Nanosecond()
	//log.Println("消耗时间:", (start-end)/1000)
	//return out.String() != ""
	status, err := dst_cli_window.DstCliClient.DstStatus(clusterName, level)
	if err != nil {
		return false
	}
	return status.Status
}

func (g *GameService) shutdownLevel(clusterName, level string) {
	if !g.GetLevelStatus(clusterName, level) {
		return
	}

	shell := "c_shutdown(true)"
	log.Println("正在shutdown世界", "cluster: ", clusterName, "level: ", level, "command: ", shell)
	_, err := dst_cli_window.DstCliClient.Command(clusterName, level, shell)
	if err != nil {
		log.Println("shut down " + clusterName + " " + level + " error: " + err.Error())
		log.Println("shutdown 失败，将强制杀掉世界")
	}
}

// TODO 强制kill 掉进程
func (g *GameService) killLevel(clusterName, level string) {

}

func (g *GameService) StartLevel(clusterName, level string, bin, beta int) {
	g.StopLevel(clusterName, level)
	g.LaunchLevel(clusterName, level, bin, beta)
	ClearScreen()
}

func (g *GameService) LaunchLevel(clusterName, level string, bin, beta int) {

	g.logRecord.RecordLog(clusterName, level, model.RUN)

	cluster := clusterUtils.GetCluster(clusterName)
	dstInstallDir := cluster.ForceInstallDir
	ugcDirectory := cluster.Ugc_directory
	persistent_storage_root := cluster.Persistent_storage_root
	conf_dir := cluster.Conf_dir

	command := ""
	title := cluster.ClusterName + "_" + level
	if bin == 64 {
		dstInstallDir = filepath.Join(dstInstallDir, "bin64")
		command = "cd /d " + dstInstallDir + " && Start \"" + title + "\" dontstarve_dedicated_server_nullrenderer_x64 -console -cluster " + clusterName + " -shard " + level
	} else {
		dstInstallDir = filepath.Join(dstInstallDir, "bin")
		command = "cd /d " + dstInstallDir + " && Start \"" + title + "\" dontstarve_dedicated_server_nullrenderer -console -cluster " + clusterName + " -shard " + level
	}

	if ugcDirectory != "" {
		command += " -ugc_directory " + ugcDirectory
	}
	if persistent_storage_root != "" {
		command += " -persistent_storage_root " + persistent_storage_root
	}
	if conf_dir != "" {
		command += " -conf_dir " + conf_dir
	}

	// 创建一个命令对象
	cmd := exec.Command("cmd.exe", "/C", command)
	log.Println(command)
	// 设置新窗口属性
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.CmdLine = "cmd.exe /K " + command

	// 启动命令
	err := cmd.Start()
	if err != nil {
		log.Panicln(err)
	}

}

func (g *GameService) StopLevel(clusterName, level string) {

	g.logRecord.RecordLog(clusterName, level, model.STOP)

	g.shutdownLevel(clusterName, level)

	time.Sleep(3 * time.Second)

	if g.GetLevelStatus(clusterName, level) {
		var i uint8 = 1
		for {
			if g.GetLevelStatus(clusterName, level) {
				break
			}
			g.shutdownLevel(clusterName, level)
			time.Sleep(1 * time.Second)
			i++
			if i > 3 {
				break
			}
		}
	}
	g.killLevel(clusterName, level)
}

// StopGame TODO 适配windows
func (g *GameService) StopGame(clusterName string) {

	config, err := levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		log.Panicln(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(config.LevelList))
	for i := range config.LevelList {
		go func(i int) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			levelName := config.LevelList[i].File
			g.StopLevel(clusterName, levelName)
		}(i)
	}
	wg.Wait()
}

func (g *GameService) StartGame(clusterName string) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.StopGame(clusterName)
	cluster := clusterUtils.GetCluster(clusterName)
	bin := cluster.Bin
	beta := cluster.Beta

	config, err := levelConfigUtils.GetLevelConfig(clusterName)
	if err != nil {
		log.Panicln(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(config.LevelList))
	for i := range config.LevelList {
		go func(i int) {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			levelName := config.LevelList[i].File
			g.LaunchLevel(clusterName, levelName, bin, beta)
		}(i)
	}
	ClearScreen()
	wg.Wait()
}

func (g *GameService) PsAuxSpecified(clusterName, level string) *vo.DstPsVo {
	dstPsVo := vo.NewDstPsVo()
	cmd := "ps -aux | grep -v grep | grep -v tail | grep " + clusterName + "  | grep " + level + " | sed -n '2P' |awk '{print $3, $4, $5, $6}'"

	info, err := shellUtils.Shell(cmd)
	if err != nil {
		log.Println(cmd + " error: " + err.Error())
		return dstPsVo
	}
	if info == "" {
		return dstPsVo
	}

	arr := strings.Split(info, " ")
	dstPsVo.CpuUage = strings.Replace(arr[0], "\n", "", -1)
	dstPsVo.MemUage = strings.Replace(arr[1], "\n", "", -1)
	dstPsVo.VSZ = strings.Replace(arr[2], "\n", "", -1)
	dstPsVo.RSS = strings.Replace(arr[3], "\n", "", -1)

	return dstPsVo
}

type SystemInfo struct {
	HostInfo      *systemUtils.HostInfo `json:"host"`
	CpuInfo       *systemUtils.CpuInfo  `json:"cpu"`
	MemInfo       *systemUtils.MemInfo  `json:"mem"`
	DiskInfo      *systemUtils.DiskInfo `json:"disk"`
	PanelMemUsage uint64                `json:"panelMemUsage"`
	PanelCpuUsage float64               `json:"panelCpuUsage"`
}

func (g *GameService) GetSystemInfo(clusterName string) *SystemInfo {
	var wg sync.WaitGroup
	wg.Add(5)

	dashboardVO := SystemInfo{}
	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		dashboardVO.HostInfo = systemUtils.GetHostInfo()
	}()

	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		dashboardVO.CpuInfo = systemUtils.GetCpuInfo()
	}()

	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		dashboardVO.MemInfo = systemUtils.GetMemInfo()
	}()

	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		dashboardVO.DiskInfo = systemUtils.GetDiskInfo()
	}()

	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		dashboardVO.PanelMemUsage = m.Alloc / 1024 // 将字节转换为MB

		// 获取当前程序使用的CPU信息
		//startCPU, _ := cpu.Percent(time.Second, false)
		//time.Sleep(1 * time.Second) // 假设程序运行1秒
		//endCPU, _ := cpu.Percent(time.Second, false)
		//cpuUsage := endCPU[0] - startCPU[0]
		//dashboardVO.PanelCpuUsage = cpuUsage

	}()

	wg.Wait()
	return &dashboardVO
}
