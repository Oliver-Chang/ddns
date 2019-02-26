# DDNS
实现一个自己的 ddns 服务，只支持 cloudflare。只更新 ipv6 地址。

## Usage
添加配置文件到 $HOME `.ddns.toml` 或 `.ddns.yaml`

```toml
AuthEmail="xxx@xxx.xxx"
AuthKey="xxx"
Subdomain="xxx.xxx.xxx"
```

然后启动该程序。
