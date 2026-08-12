# WsureDev Xray redirect compatibility build

本分支基于 v2rayA `v2.4.11`（commit
`1550d843d02b1ad8dd7b832197ba7d97f59d27fe`），用于固化 v2rayA 与新版
Xray-core 的源码级合并构建，并修复 Linux `redirect` 透明代理路径。

对应的 Xray-core 是 `XTLS/Xray-core` `main` 分支的精确 commit
`7d214f8b094f75322fa3990f8aadad1c912f24f5`（紧随 `v26.7.28` 快照）。该版本
已经固定在 `core/xray` Git submodule 中；Xray 源码本身没有私有修改。

## 本分支的修改

- redirect 的 dokodemo-door IPv4 入站改为监听 `0.0.0.0`，使
  `PREROUTING` 转发进入的 TCP 流量能被接收，并保持
  `followRedirect: true` 与 `sockopt.tproxy: redirect`。
- IPv6 启用时创建独立的 `[::]`/`v6only` 入站，避免与 IPv4 wildcard
  socket 冲突；IPv6 未启用时不创建该入站。
- 明确将 redirect 数据路径限制为 TCP；TCP/UDP DNS 使用 v2rayA 的
  `52353` DNS listener。
- 修正 legacy iptables DNS 规则：TCP 和 UDP 分别声明协议，兼容
  `iptables v1.8.x (nf_tables)` 对 `REDIRECT --to-port` 的校验。
- 调整 DNS hook 的插入顺序，使 TCP/53 优先进入 DNS listener，而不是
  通用的 `52345` redirect 入站。
- 避免 redirect 模式重复安装另一套 DNS REDIRECT hook。
- 增加 redirect 配置与 iptables 规则的回归测试。

网页日志不是源码行为修改。下面的 Compose 通过
`V2RAYA_LOG_FILE=/etc/v2raya/v2raya.log` 启用文件日志，因此 Web UI 可以读取
日志，并随 `/etc/v2raya` 数据目录持久化。

## 获取源码

```sh
git clone --recurse-submodules \
  --branch main \
  git@github.com:WsureDev/v2rayA.git
cd v2rayA
```

已有 clone 需要执行：

```sh
git submodule sync --recursive
git submodule update --init --recursive
```

可用以下命令核对 Xray 版本：

```sh
git -C core/xray rev-parse HEAD
# 7d214f8b094f75322fa3990f8aadad1c912f24f5
```

## Docker 构建与使用（推荐）

需要 Docker BuildKit/Buildx 和 Compose v2：

```sh
docker compose -f compose.xray-redirect.yaml build
docker compose -f compose.xray-redirect.yaml up -d
docker compose -f compose.xray-redirect.yaml ps
```

Compose 从当前 Fork 源码构建 `v2raya` 和合并 Xray 的 `v2raya_core`，并以
与 v2rayA 2.4.11 匹配的官方镜像层作为 runtime。默认接口为：

| 用途 | 地址/端口 |
| --- | --- |
| Web 管理 | `http://<宿主机 IP>:2018` |
| SOCKS5 | `<宿主机 IP>:20180` |
| HTTP | `<宿主机 IP>:20181` |
| 规则分流 HTTP | `<宿主机 IP>:20182` |

首次进入 Web UI 后，由用户创建账号、导入订阅、连接节点并选择 redirect
透明代理模式。数据默认保存在仓库下的 `./data`；可通过
`V2RAYA_DATA_DIR=/path/to/data` 指定其他目录。其他 Linux Docker 容器可通过
宿主端口使用代理，例如：

```yaml
extra_hosts:
  - "host.docker.internal:host-gateway"
environment:
  HTTPS_PROXY: http://host.docker.internal:20181
```

也可以让业务容器加入 `v2raya-xray_default` 网络，直接使用
`v2raya-xray:20170`（SOCKS5）、`:20171`（HTTP）或 `:20172`（规则 HTTP）。

停止服务但保留数据：

```sh
docker compose -f compose.xray-redirect.yaml down
```

## 原生源码构建

需要 Go 1.26、Node.js 24、Yarn、Git，以及已经初始化的 Xray submodule：

```sh
./build.sh
```

生成物位于仓库根目录：

- `v2raya`：Web/系统服务端；
- `v2raya_core`：与上述固定 Xray-core 源码合并构建的 core。

修改 redirect 后的针对性测试可在 Go 1.26 环境中执行。v2rayA 当前的参数
解析器不能接收 Go 1.26 普通 `go test` 注入的 `-test.testlogfile`，因此先编译
再运行测试二进制：

```sh
(cd service && \
  go test -c -o /tmp/v2raya-iptables.test ./kernel/iptables && \
  /tmp/v2raya-iptables.test && \
  go test -c -o /tmp/v2raya-kernel.test ./kernel/v2ray && \
  /tmp/v2raya-kernel.test)
```

---

# v2rayA [![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/v2rayA/v2raya)](https://hub.docker.com/r/mzz2017/v2raya) [![Travis (.org)](https://img.shields.io/travis/v2rayA/v2rayA?label=travis-ci%20build)](https://travis-ci.org/v2rayA/v2rayA)

[**English**](https://github.com/v2rayA/v2rayA/blob/main/README.md)&nbsp;&nbsp;&nbsp;[**简体中文**](https://github.com/v2rayA/v2rayA/blob/main/README_zh.md)

v2rayA is a V2Ray client supporting global transparent proxy on Linux and system proxy on Windows and macOS, it is compatible with SS, SSR, Trojan(trojan-go), Tuic and [Juicity](https://github.com/juicity) protocols. [[SSR protocol list]](https://github.com/v2rayA/shadowsocksR/blob/main/README.md#ss-encrypting-algorithm)

We are committed to providing the simplest operation and meet most needs.

Thanks to the advantages of Web GUI, you can not only use it on your local computer, but also easily deploy it on a router or NAS.

Project：https://github.com/v2rayA/v2rayA


## Usage

v2rayA mainly provides the following methods of installation:

1. Install from apt-source or AUR
2. Docker
3. Our self-built [scoop bucket](https://github.com/v2rayA/v2raya-scoop) (for Windows users)
4. Our self-built [homebrew tap](https://github.com/v2rayA/homebrew-v2raya)
5. Our self-built [OpenWrt repo](https://github.com/v2rayA/v2raya-openwrt) and OpenWrt's official repo(from OpenWrt version 22.03)
6. Microsoft winget: https://winstall.app/apps/v2rayA.v2rayA
7. Ubuntu Snap: https://snapcraft.io/v2raya
8. Binary file and installation package from GitHub releases

See [**v2rayA - Docs**](https://v2raya.org/en/docs/prologue/introduction/)


## Screenshot

<img src="https://i.loli.net/2020/04/19/gt3NqOMiafYbp7L.png" border="0">

## Statement

1. The program does not store any user data in the cloud, and all user data is stored in local.
2. **Do not use this project for illegal purposes.**

## Credits

[hq450/fancyss](https://github.com/hq450/fancyss)

[ToutyRater/v2ray-guide](https://github.com/ToutyRater/v2ray-guide/blob/master/routing/sitedata.md)

[nadoo/glider](https://github.com/nadoo/glider)

[Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)

[zfl9/ss-tproxy](https://github.com/zfl9/ss-tproxy/blob/master/ss-tproxy)

## Stargazers over time

[![Stargazers over time](https://starchart.cc/v2rayA/v2rayA.svg)](https://starchart.cc/v2rayA/v2rayA)

## License

[![License: AGPL v3-only](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
