# go语言实现的网站流量分析|统计脚本  
使用nginx生成1 X 1的空白图片,用户访问页面,携带统计数据请求1 X 1图片, 统计信息存入nginx访问日志  
使用go读取nginx日志, 分析统计日志数据  
  
环境需要
* nginx_gif
* golang
* redis  
  
[流程Process](https://github.com/zheng59521/web_log_record/tree/dev/process.png)

## 进度
    ** 格式化pu uv数据为json数据 ** 
    


