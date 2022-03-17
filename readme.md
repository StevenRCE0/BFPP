# 备份盘盘

## 简介

这是一个使用 Go 编写的用来定期备份的小玩意儿。需要提前安装 [bypy](https://github.com/houtianze/bypy)。

## 使用

1. 配置 bypy。
2. 按说明填写两个 json。
3. ```go build main.go```
4. ```main <path/to/configuration.json> <path/to/backupList.json>```

## macOS 定期备份

1. 按说明修改 plist，搁进 `~/Library/LaunchAgents`。
2. ```id -u``` 看看 UID。
3. 启用配置。

   ```sudo launchctl bootstrap gui/<UID> ~/Library/LaunchAgents/com.coils.pan.plist```
4. 运行一次测试。

   ```sudo launchctl kickstart gui/<UID>/com.coils.pan```
