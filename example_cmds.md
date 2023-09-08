#### run chain and only log error msgs
```
ignite chain serve -v | grep -v "INF"
```

#### submit eth/usd val
```
luqchaind tx luqchain submit-val ad 1900 --from alice
```
