// # Description: Get os tcp connection state
// #
// # Author: dzz_dkmen@163.com
// #
// # Create date: 2024-03-15
// #
// # Modify history:
// # No  Name                 Date         Description
// # --- -------------------- -----------  ------------------------------------------------
// # 1   Create  iptracking   2024-03-15
// # --------------------------------------------------------------------------------------

// 需求:
// 已知源头IP，寻找上下游关联IP端口及进程; => v0.1.3-rc1 已实现获取上下游关联IP端口
// 已知目标IP，寻找上下游关联IP端口及进程; => v0.1.3-rc1 已实现获取上下游关联IP端口
// 已知道恶意IP，回放访问路径及端口及链接状态; => v1.0
// 已知道重保IP，实时监控重保IP上下游关联IP端口关系; => v1.0

// 实现思路:
// Linux Kernel 网络子系统 , /proc/net/tcp 文件提供当前系统中 TCP 协议栈的状态信息。
// 采用通过解析/proc/net/tcp文件中已建立的 TCP 连接列表、连接的状态、本地地址和端口、远程地址和端口等数据并实时落入ClickHouse。

 //   linux /proc/net/tcp 字段描述
 //   46: 010310AC:9C4C 030310AC:1770 01 
 //   |      |      |      |      |   |--> connection state（套接字状态）
 //   |      |      |      |      |------> remote TCP port number（远端端口，主机字节序）
 //   |      |      |      |-------------> remote IPv4 address（远端IP，网络字节序）
 //   |      |      |--------------------> local TCP port number（本地端口，主机字节序）
 //   |      |---------------------------> local IPv4 address（本地IP，网络字节序）
 //   |----------------------------------> number of entry

 //  00000150:00000000 01:00000019 00000000  
 //     |        |     |     |       |--> number of unrecovered RTO timeouts（超时重传次数）
 //     |        |     |     |----------> number of jiffies until timer expires（超时时间，单位是jiffies）
 //     |        |     |----------------> timer_active (定时器类型，see below)
 //     |        |----------------------> receive-queue（根据状态不同有不同表示,see below）
 //     |-------------------------------> transmit-queue(发送队列中数据长度)

 //    connection state（套接字状态）:
 //    TCP_ESTABLISHED:1
 //    TCP_SYN_SENT:2
 //    TCP_SYN_RECV:3
 //    TCP_FIN_WAIT1:4
 //    TCP_FIN_WAIT2:5
 //    TCP_TIME_WAIT:6
 //    TCP_CLOSE:7
 //    TCP_CLOSE_WAIT:8
 //    TCP_LAST_ACL:9
 //    TCP_LISTEN:10
 //    TCP_CLOSING:11

 // input /proc/net/tcp file ;
 // output srt "Last Modified Time: 2024-03-15T13:54:30+08:00 , Local Address: 172.16.10.39:10050 , Remote Address: 10.0.0.108:39020"


package main
import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
    "net"
    "encoding/binary"
    "encoding/hex"
    "time"
)

func main() {
    // 打开 /proc/net/tcp 文件
    // file, err := os.Open("/proc/net/tcp")
    file, err := os.Open("tcp")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }

    defer file.Close()

    filePath := "tcp" // 文件路径
    // filePath := "/proc/net/tcp" // 文件路径
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }    
    modTime := fileInfo.ModTime()
    // fmt.Println("Last Modified Time:", modTime.Format(time.RFC3339))

    // 逐行读取文件内容
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        // 忽略文件头部分
        if strings.HasPrefix(line, "  sl") {
            continue
        }
        // 解析每行内容
        fields := strings.Fields(line)

        localAddr, remoteAddr, err := parseAddresses(fields[1], fields[2])
        tpcstatcode := fields[3]

        if err != nil {
            fmt.Println("Error parsing addresses:", err)
            continue
        }

        fmt.Printf("LastModifiedTime %s LocalAddress %s RemoteAddress %s TCPCode %s\n", modTime.Format(time.RFC3339), localAddr, remoteAddr ,tpcstatcode)

    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
    }
}

// 解析本地和远程地址
func parseAddresses(local, remote string) (string, string, error) {
    // fmt.Printf("%s \n", local)

    localAddr := strings.Split(local, ":")[0]
    localPort := strings.Split(local, ":")[1]
    // fmt.Printf("%s \n", localAddr)

    remoteAddr := strings.Split(remote, ":")[0]
    remotePort := strings.Split(remote, ":")[1]
    // fmt.Printf("%s", remoteAddrhexPort)

    localAddrPort :=  parseAddressesPorttoStr(localAddr,localPort)
    // fmt.Printf("localAddrPort->%s === \n", localAddrPort)

    remoteAddrPort :=  parseAddressesPorttoStr(remoteAddr,remotePort)
    // fmt.Printf("remoteAddrPort->%s === \n", remoteAddrPort)

    return fmt.Sprintf("%s", localAddrPort,),fmt.Sprintf("%s", remoteAddrPort), nil

}

// 解析进程信息
func parseProcessInfo(inode string) string {
    return inode
}

// 解析ip和端口
func parseAddressesPorttoStr(Addr string,Port string) string{

    bytesIP, err := hex.DecodeString(Addr)
    fmt.Sprintf("%s", err)
    uint32IP := binary.LittleEndian.Uint32(bytesIP) //转换为主机字节序
    IP := make(net.IP, 4)
    binary.BigEndian.PutUint32(IP, uint32IP)
    strAddr := IP.String()
    strPort, err := strconv.ParseUint(Port, 16, 32)
    // fmt.Printf("strAddr->%s and strPort->%s ", strAddr,strPort)
    // fmt.Printf("%s:%s", strAddr,strPort)

    return fmt.Sprintf("%s:%d", strAddr,strPort)

}
