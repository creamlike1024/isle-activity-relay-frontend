# isle-activity-relay-frontend
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcreamlike1024%2Fisle-activity-relay-frontend.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcreamlike1024%2Fisle-activity-relay-frontend?ref=badge_shield)


[relay.isle.moe](https://relay.isle.moe)

原理是从 [activity-relay](https://github.com/yukimochi/Activity-Relay) 的 redis 获取订阅信息，然后生成信息页面

你可以修改 `template_info.html`, `template_head.html` 和 `template_log.html` 三个模板来生成自己的页面，配合 crontab 定时运行程序生成页面

配置文件：`config.yml`

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fcreamlike1024%2Fisle-activity-relay-frontend.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fcreamlike1024%2Fisle-activity-relay-frontend?ref=badge_large)