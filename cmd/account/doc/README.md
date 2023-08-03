# 功能点

## 注册账号
新用户注册新账号。
用户输入(accountName,password)，查询MySQL：账号是否存在
* 存在: 返回账号已存在
* 不存在：插入新记录到MySQL

## 注销账号
已有用户注销账号。
用户输入(accountName)，查询MySQL：账号是否存在
* 存在：删除该记录, MQ通知删除Redis关于该账号的缓存
* 不存在: 返回账号不存在

## 用户登录
用户输入(accountName,password)，查询Redis：该accountName对应的Token是否存在
* 存在：直接返回Token
* 不存在：查询数据库，查询账号是否存在
  * 存在：生成Token, 设置过期时间，存到Redis中
  * 不存在: 返回账号不存在

## 用户退出登录
用户输入(accountName),查询Redis：该accountName对应的记录是否存在
* 存在: 删除，并返回退出成功
* 不存在: 返回失败，用户早已退出