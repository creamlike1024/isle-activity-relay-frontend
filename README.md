# isle-activity-relay-frontend

relay.isle.moe

原理是从 [activity-relay](https://github.com/yukimochi/Activity-Relay) 的 redis 获取订阅信息，然后通过 markdown 生成页面

你可以修改 `template.md` 和 `template_head.html` 两个模板来生成自己的页面，配合 crontab 定时运行程序生成页面

配置文件：`config.yml`