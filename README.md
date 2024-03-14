# quant

## 安装

```bash
go install github.com/tensorchen/quant/cmd/quant@latest
```

## 配置

```yaml
tquant:
  port:
  token:

long_bridge:
  app_key:
  app_secret:
  access_token:

discord:
  id:
  token:
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
  "token": "{{token}}",
  "trade": {
    "ticker": "{{ticker}}",
    "exchange":"{{exchange}}",
    "strategy": {
      "order": {
        "action": "{{strategy.order.action}}",
        "contracts": "{{strategy.order.contracts}}"
      }
    }
  }
}
```