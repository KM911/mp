::用户变量所在位置：`HKEY_CURRENT_USER\Environment`
set USERregpath=HKEY_CURRENT_USER\Environment

::系统变量所在位置：`HKEY_LOCAL_MACHINE\SYSTEM\ControlSet001\Control\Session Manager\Environment`
set MACHINEregpath=HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment

::用户环境变量
::新建
// 
reg add "%USERregpath%" /v "GOPATH" /t REG_SZ /d "%UserProfile%\gopath" /f

::PATH 追加
::查询原有 PATH 的值
for /F "tokens=3*" %i in ('reg query "%USERregpath%" /v "path" 2^>nul') do echo %i
for /F "tokens=3*" %i in ('reg query "%MACHINEregpath%" /v "path" 2^>nul') do echo %i
::在 .bat 或者 .cmd 批处理文件中，%i 应该写成 %%i
for /F "tokens=3*" %%i in ('reg query "%USERregpath%" /v "path" 2^>nul') do echo %%i
for /F "tokens=3*" %%i in ('reg query "%MACHINEregpath%" /v "path" 2^>nul') do echo %%i

for /F "tokens=3*" %i in ('reg query "%USERregpath%" /v "path" 2^>nul') do ( set USERpath=%i)
echo USERpath=%USERpath%
reg add "%USERregpath%" /v "Path" /t REG_EXPAND_SZ /d ""%USERpath%"%GOPATH%\bin;" /f
:: 经过测试，巨硬公司的不同 Windows 10 版本的 PATH 变量竟然写法不一样，有的以分号结尾，有的没有分号。所以命令还是要加上分号分割，结尾分号取消
reg add "%USERregpath%" /v "Path" /t REG_EXPAND_SZ /d "%USERpath%;%GOPATH%\bin" /f

::系统环境变量
::新建
reg add "%MACHINEregpath%" /v "GOROOT" /t REG_SZ /d "C:\go" /f

::PATH 追加
for /F "tokens=3*" %i in ('reg query "%MACHINEregpath%" /v "path" 2^>nul') do ( set MACHINEpath=%i)
echo MACHINEpath=%MACHINEpath%
reg add "%MACHINEregpath%" /v "Path" /t REG_EXPAND_SZ /d "%USERpath%;%GOROOT%\bin" /f