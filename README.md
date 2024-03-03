# quant

## 安装

```bash
go install github.com/tensorchen/quant/cmd/quant@latest
```

## 配置

```bash
# 设置服务启动监听的端口号
export TQUANT_PORT=""
# 设置量化交易密钥
export TQUANT_TOKEN=""

# 将交易信息推送到discord的相关配置
export DISCORD_ID=""
export DISCORD_TOKEN=""

# 配置长桥证券的相关配置
export LONGPORT_APP_KEY=""
export LONGPORT_APP_SECRET=""
export LONGPORT_ACCESS_TOKEN=""
```

## 使用

| 字段                             | 解释   | 备注 |
|--------------------------------|------|----|
| token                          | 密钥   |    |
| trade.ticker                   | 股票代码 |    |
| trade.exchange                 | 交易所  |    |
| trade.strategy.order.action    | 交易操作 |    |
| trade.strategy.order.contracts | 交易数量 |    |

```json
{
  "token": "",
  "trade": {
    "ticker": "",
    "exchange":"",
    "strategy": {
      "order": {
        "action": "",
        "contracts": ""
      }
    }
  }
}
```