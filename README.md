mp -- manage your environment path

## 功能

- [ ] 添加环境变量
- [ ] 导出/加载用于备份和恢复
- [ ] 清除无效值

## 安装 

```bash

git clone  

cd mp 
# 将mp添加到PATH
pwd | ./mp
# 重启终端 使修改生效 
cd folder_your_want_to_add_to_path
pwd | ./mp
# 或者
mp folder_your_want_to_add_to_path

# 添加新的环境变量 
mp hello world 
```

## 原理

linux和win一样,添加完环境变量后,需要重启软件重新加载环境变量,最新的修改才可以生效.

### Linux

向 `~/.bashrc`文件中添加 `export ` 语句,实现添加环境变量的效果.

未实现失效环境变量移除的效果.

### windows

利用win提供的`reg`命令,修改注册名表.

比如添加一个全新的环境变量
```
reg add HKEY_CURRENT_USER\Environment /v "k" /t REG_SZ /d "v" /f
```

win添加环境变量后需要重启才可以生效,所以立即使用该环境变量是无效的.


